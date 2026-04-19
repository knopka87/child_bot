// src/pages/Check/CheckUploadPage.tsx
import { useState, useRef, useCallback } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { ArrowLeft, Camera, Image, AlertCircle } from 'lucide-react';
import { motion } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import type { CheckScenario } from '@/types/check';
import styles from './CheckUploadPage.module.css';

export default function CheckUploadPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const analytics = useAnalytics();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const scenario = (searchParams.get('scenario') as CheckScenario) || 'single_photo';
  const isTwoPhoto = scenario === 'two_photo';

  const [error, setError] = useState<string | null>(null);

  const handleFileSelect = useCallback(async () => {
    setError(null);
    if (isTwoPhoto) {
      navigate('/check/upload-images?scenario=two_photo');
    } else {
      fileInputRef.current?.click();
    }
  }, [isTwoPhoto, navigate]);

  const handleCamera = useCallback(async () => {
    setError(null);
    if (isTwoPhoto) {
      navigate('/check/upload-images?scenario=two_photo');
    } else {
      // Для одного фото — открываем камеру напрямую
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      analytics.trackEvent('check_camera_clicked', {
        child_profile_id: childProfileId,
        scenario,
      });
      // Для single_photo пока ведём на тот же роут выбора файла
      // В будущем можно открыть getUserMedia
      fileInputRef.current?.click();
    }
  }, [isTwoPhoto, navigate, analytics, scenario]);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Валидация типа
    if (!file.type.startsWith('image/')) {
      setError('Пожалуйста, выбери изображение (JPG, PNG)');
      return;
    }

    // Валидация размера (макс 10 МБ)
    if (file.size > 10 * 1024 * 1024) {
      setError('Файл слишком большой. Максимум 10 МБ');
      return;
    }

    setError(null);

    const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
    analytics.trackEvent('check_file_selected', {
      child_profile_id: childProfileId,
      scenario,
    });

    // Convert file to base64 and store in sessionStorage for the quality page
    const reader = new FileReader();
    reader.onload = () => {
      const base64 = reader.result as string;
      sessionStorage.setItem('check_single_photo_data', JSON.stringify({
        fileName: file.name,
        fileType: file.type,
        fileSize: file.size,
        base64: base64,
      }));

      // Для single_photo — переходим на страницу проверки качества
      navigate('/check/quality-single');
    };
    reader.readAsDataURL(file);
  };

  return (
    <div className={styles.container}>
      {/* Скрытый input для выбора файла */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileChange}
        className={styles.hiddenInput}
      />

      {/* Кнопка назад */}
      <button
        onClick={() => navigate('/check/scenario')}
        className={styles.backButton}
      >
        <ArrowLeft size={20} />
        <span className={styles.backText}>Назад</span>
      </button>

      {/* Заголовок */}
      <div className={styles.header}>
        <div className={styles.iconCircle}>
          <Image size={30} className={styles.iconCircleIcon} />
        </div>
        <h2 className={styles.title}>Проверка ДЗ</h2>
        <p className={styles.subtitle}>
          {isTwoPhoto
            ? 'Загрузи первое фото'
            : 'Загрузи фото задания из учебника'}
        </p>
      </div>

      {/* Кнопки выбора */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className={styles.actions}
      >
        {/* Выбрать изображение */}
        <button
          onClick={handleFileSelect}
          className={styles.actionButton}
        >
          <div className={`${styles.actionIcon} ${styles.actionIconPrimary}`}>
            <Image size={24} />
          </div>
          <div className={styles.actionText}>
            <p className={styles.actionTitle}>Выбрать изображение</p>
            <p className={styles.actionDesc}>JPG, PNG</p>
          </div>
        </button>

        {/* Сфотографировать */}
        <button
          onClick={handleCamera}
          className={styles.actionButton}
        >
          <div className={`${styles.actionIcon} ${styles.actionIconCamera}`}>
            <Camera size={24} />
          </div>
          <div className={styles.actionText}>
            <p className={styles.actionTitle}>Сфотографировать</p>
            <p className={styles.actionDesc}>Открыть камеру</p>
          </div>
        </button>

        {/* Ошибка */}
        {error && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className={styles.errorBox}
          >
            <AlertCircle size={20} className={styles.errorIcon} />
            <div>
              <p className={styles.errorText}>{error}</p>
              <button
                onClick={() => setError(null)}
                className={styles.errorRetry}
              >
                Попробовать снова
              </button>
            </div>
          </motion.div>
        )}
      </motion.div>
    </div>
  );
}
