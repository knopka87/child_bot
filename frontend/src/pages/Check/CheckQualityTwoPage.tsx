// src/pages/Check/CheckQualityTwoPage.tsx
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Check, RefreshCw, Crop } from 'lucide-react';
import { motion } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { checkAPI } from '@/api/check';
import imageCompression from 'browser-image-compression';
import type { CheckScenario } from '@/types/check';
import styles from './CheckQualityTwoPage.module.css';

interface StoredFileData {
  fileName: string;
  fileType: string;
  fileSize: number;
  base64: string;
}

export default function CheckQualityTwoPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();

  const [taskPreview, setTaskPreview] = useState<string | null>(null);
  const [answerPreview, setAnswerPreview] = useState<string | null>(null);
  const [taskFile, setTaskFile] = useState<File | null>(null);
  const [answerFile, setAnswerFile] = useState<File | null>(null);
  const [compressing, setCompressing] = useState<string | null>(null);
  const [uploading, setUploading] = useState<string | null>(null);

  // Load files from sessionStorage on mount
  useEffect(() => {
    const taskDataStr = sessionStorage.getItem('check_task_photo');
    const answerDataStr = sessionStorage.getItem('check_answer_photo');

    if (!taskDataStr || !answerDataStr) {
      navigate('/check/scenario');
      return;
    }

    try {
      const taskData: StoredFileData = JSON.parse(taskDataStr);
      const answerData: StoredFileData = JSON.parse(answerDataStr);

      setTaskPreview(taskData.base64);
      setAnswerPreview(answerData.base64);

      // Convert base64 data URI to File without using fetch()
      const base64ToFile = (data: StoredFileData): File => {
        // Extract the base64 part (remove data:image/...;base64, prefix)
        const base64Index = data.base64.indexOf(',');
        const base64String = base64Index !== -1 ? data.base64.substring(base64Index + 1) : data.base64;
        
        // Decode base64 to binary string
        const binaryString = atob(base64String);
        
        // Convert binary string to Uint8Array
        const bytes = new Uint8Array(binaryString.length);
        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i);
        }
        
        // Create File from blob
        return new File([bytes], data.fileName, { type: data.fileType });
      };

      try {
        const tFile = base64ToFile(taskData);
        const aFile = base64ToFile(answerData);
        
        setTaskFile(tFile);
        setAnswerFile(aFile);
      } catch (err) {
        console.error('[CheckQualityTwoPage] Failed to convert base64 to Files:', err);
        navigate('/check/upload-images?scenario=two_photo');
      }
    } catch (error) {
      console.error('[CheckQualityTwoPage] Failed to parse stored data:', error);
      navigate('/check/upload-images?scenario=two_photo');
    }
  }, [navigate]);

  if (!taskFile || !answerFile || !taskPreview || !answerPreview) {
    return null; // Will redirect in useEffect
  }

  const handleConfirm = async () => {
    try {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      if (!childProfileId) return;

      analytics.trackEvent('check_quality_confirmed', {
        child_profile_id: childProfileId,
        scenario: 'two_photo',
      });

      // 1. Создаём attempt с scenario two_photo
      setCompressing('both');
      const attempt = await checkAPI.createAttempt(childProfileId, 'two_photo' as CheckScenario);

      // 2. Сжимаем изображения
      const compressOptions = {
        maxSizeMB: 7,
        maxWidthOrHeight: 1920,
        useWebWorker: true,
        initialQuality: 0.8,
      };

      const [compressedTask, compressedAnswer] = await Promise.all([
        imageCompression(taskFile, compressOptions),
        imageCompression(answerFile, compressOptions),
      ]);

      // 3. Загружаем task image
      setUploading('task');
      await checkAPI.uploadImage(attempt.id, 'task', compressedTask);

      // 4. Загружаем answer image
      setUploading('answer');
      await checkAPI.uploadImage(attempt.id, 'answer', compressedAnswer);

      analytics.trackEvent('check_images_uploaded', {
        child_profile_id: childProfileId,
        attempt_id: attempt.id,
      });

      // Clear sessionStorage
      sessionStorage.removeItem('check_task_photo');
      sessionStorage.removeItem('check_answer_photo');

      // 5. Переходим к обработке
      navigate('/check/processing', { 
        state: { 
          attemptId: attempt.id,
          scenario: 'two_photo',
        } 
      });
    } catch (error) {
      console.error('[CheckQualityTwoPage] Failed:', error);
      // В случае ошибки — показываем alert и возвращаем назад
      alert('Не удалось загрузить изображения. Попробуй ещё раз.');
      navigate('/check/upload-images?scenario=two_photo');
    }
  };

  const handleRetake = () => {
    sessionStorage.removeItem('check_task_photo');
    sessionStorage.removeItem('check_answer_photo');
    navigate('/check/upload-images?scenario=two_photo');
  };

  const isLoading = compressing !== null || uploading !== null;

  const getLoadingText = () => {
    if (compressing) return 'Сжимаем изображения...';
    if (uploading === 'task') return 'Загружаем задание...';
    if (uploading === 'answer') return 'Загружаем ответ...';
    return '';
  };

  return (
    <div className={styles.container}>
      {/* Кнопка назад */}
      <button
        onClick={() => {
          sessionStorage.removeItem('check_task_photo');
          sessionStorage.removeItem('check_answer_photo');
          navigate(-1);
        }}
        className={styles.backButton}
        disabled={isLoading}
      >
        <ArrowLeft size={20} />
        <span className={styles.backText}>Назад</span>
      </button>

      {/* Заголовок */}
      <div className={styles.header}>
        <h2 className={styles.title}>Всё ли видно?</h2>
        <p className={styles.subtitle}>
          Проверь, что оба фото хорошо видны
        </p>
      </div>

      {/* Превью изображений */}
      <div className={styles.previews}>
        {/* Задание */}
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          className={styles.previewCard}
        >
          <span className={styles.previewLabel}>Задание</span>
          <div className={styles.previewImageContainer}>
            <img
              src={taskPreview}
              alt="Задание"
              className={styles.previewImage}
            />
          </div>
        </motion.div>

        {/* Ответ */}
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.1 }}
          className={styles.previewCard}
        >
          <span className={styles.previewLabel}>Ответ</span>
          <div className={styles.previewImageContainer}>
            <img
              src={answerPreview}
              alt="Ответ"
              className={styles.previewImage}
            />
          </div>
        </motion.div>
      </div>

      {/* Действия */}
      <div className={styles.actions}>
        {/* Обрезать — пока заглушка */}
        <button
          onClick={() => {}}
          className={styles.actionButtonSecondary}
          disabled={isLoading}
        >
          <Crop size={18} className={styles.actionButtonIcon} />
          Обрезать
        </button>

        {/* Всё видно, продолжить */}
        <button
          onClick={handleConfirm}
          disabled={isLoading}
          className={`${styles.actionButtonPrimary} ${isLoading ? styles.actionButtonPrimaryDisabled : ''}`}
        >
          {isLoading ? (
            <span className={styles.loadingText}>{getLoadingText()}</span>
          ) : (
            <>
              <Check size={18} />
              Всё видно, продолжить
            </>
          )}
        </button>

        {/* Переснять */}
        <button
          onClick={handleRetake}
          className={styles.actionButtonOutline}
          disabled={isLoading}
        >
          <RefreshCw size={18} />
          Переснять
        </button>
      </div>
    </div>
  );
}
