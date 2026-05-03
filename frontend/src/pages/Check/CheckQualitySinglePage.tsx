// src/pages/check/CheckQualitySinglePage.tsx
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Check, RefreshCw, Crop } from 'lucide-react';
import { motion } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { checkAPI } from '@/api/check';
import imageCompression from 'browser-image-compression';
import type { CheckScenario } from '@/types/check';
import { ImageCropModal } from '@/components/ui/ImageCropModal';
import styles from './CheckQualityTwoPage.module.css';

interface StoredFileData {
  fileName: string;
  fileType: string;
  fileSize: number;
  base64: string;
}

export default function CheckQualitySinglePage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();

  const [preview, setPreview] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [compressing, setCompressing] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [showCropModal, setShowCropModal] = useState(false);

  // Load file from sessionStorage on mount
  useEffect(() => {
    const storedData = sessionStorage.getItem('check_single_photo_data');
    if (!storedData) {
      navigate('/check/scenario');
      return;
    }

    try {
      const data: StoredFileData = JSON.parse(storedData);
      setPreview(data.base64);

      // Convert base64 data URI to File without using fetch()
      const base64ToFile = (data: StoredFileData): File => {
        const base64Index = data.base64.indexOf(',');
        const base64String = base64Index !== -1 ? data.base64.substring(base64Index + 1) : data.base64;
        const binaryString = atob(base64String);
        const bytes = new Uint8Array(binaryString.length);
        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i);
        }
        return new File([bytes], data.fileName, { type: data.fileType });
      };

      try {
        const file = base64ToFile(data);
        setFile(file);
      } catch (err) {
        console.error('[CheckQualitySinglePage] Failed to convert base64 to File:', err);
        navigate('/check/upload-images?scenario=single_photo');
      }
    } catch (error) {
      console.error('[CheckQualitySinglePage] Failed to parse stored data:', error);
      navigate('/check/upload-images?scenario=single_photo');
    }
  }, [navigate]);

  if (!file || !preview) {
    return null; // Will redirect in useEffect
  }

  const handleConfirm = async () => {
    try {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      if (!childProfileId) return;

      analytics.trackEvent('check_quality_confirmed', {
        child_profile_id: childProfileId,
        scenario: 'single_photo',
      });

      // 1. Создаём attempt с scenario single_photo
      setCompressing(true);
      const attempt = await checkAPI.createAttempt(childProfileId, 'single_photo' as CheckScenario);

      // 2. Сжимаем изображение
      const compressOptions = {
        maxSizeMB: 7,
        maxWidthOrHeight: 1920,
        useWebWorker: true,
        initialQuality: 0.8,
      };

      const compressedFile = await imageCompression(file, compressOptions);

      // 3. Загружаем изображение как task
      setUploading(true);
      await checkAPI.uploadImage(attempt.id, 'task', compressedFile);

      analytics.trackEvent('check_image_uploaded', {
        child_profile_id: childProfileId,
        attempt_id: attempt.id,
        image_role: 'task',
      });

      // Для single_photo загружаем то же изображение как answer
      // (бэкенд ожидает оба изображения для check attempts)
      await checkAPI.uploadImage(attempt.id, 'answer', compressedFile);

      analytics.trackEvent('check_image_uploaded', {
        child_profile_id: childProfileId,
        attempt_id: attempt.id,
        image_role: 'answer',
      });

      // Clear sessionStorage
      sessionStorage.removeItem('check_single_photo_data');

      // 4. Переходим к обработке
      navigate('/check/processing', { 
        state: { 
          attemptId: attempt.id,
          scenario: 'single_photo',
        } 
      });
    } catch (error) {
      console.error('[CheckQualitySinglePage] Failed:', error);
      alert('Не удалось загрузить изображение. Попробуй ещё раз.');
      navigate('/check/upload-images?scenario=single_photo');
    }
  };

  const handleRetake = () => {
    sessionStorage.removeItem('check_single_photo_data');
    navigate('/check/upload-images?scenario=single_photo');
  };

  const handleCropSave = async (croppedFile: File) => {
    // Создаём превью для обрезанного изображения
    const reader = new FileReader();
    reader.onloadend = () => {
      const base64 = reader.result as string;
      setPreview(base64);
      setFile(croppedFile);

      // Обновляем в sessionStorage
      const storedData = {
        fileName: croppedFile.name,
        fileType: croppedFile.type,
        fileSize: croppedFile.size,
        base64,
      };
      sessionStorage.setItem('check_single_photo_data', JSON.stringify(storedData));
    };
    reader.readAsDataURL(croppedFile);

    setShowCropModal(false);
  };

  const isLoading = compressing || uploading;

  const getLoadingText = () => {
    if (compressing) return 'Сжимаем изображение...';
    if (uploading) return 'Загружаем изображение...';
    return '';
  };

  return (
    <div className={styles.container}>
      {/* Кнопка назад */}
      <button
        onClick={() => {
          sessionStorage.removeItem('check_single_photo_data');
          navigate('/check/upload-images?scenario=single_photo');
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
          Проверь, что задание хорошо видно на фото
        </p>
      </div>

      {/* Превью изображения */}
      <div className={styles.previews}>
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          className={styles.previewCard}
        >
          <div className={styles.previewHeader}>
            <span className={styles.previewLabel}>Задание и ответ</span>
            <button
              onClick={() => setShowCropModal(true)}
              className={styles.cropButton}
              disabled={isLoading}
            >
              <Crop size={16} />
            </button>
          </div>
          <div className={styles.previewImageContainer}>
            <img
              src={preview}
              alt="Задание"
              className={styles.previewImage}
            />
          </div>
        </motion.div>
      </div>

      {/* Действия */}
      <div className={styles.actions}>
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

      {/* Модальное окно обрезки */}
      {showCropModal && (
        <ImageCropModal
          image={preview!}
          onSave={handleCropSave}
          onClose={() => setShowCropModal(false)}
          title="Обрезать изображение"
        />
      )}
    </div>
  );
}
