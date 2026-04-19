// src/pages/Help/HelpImageUploader.tsx
import { useState, useRef, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Plus, RefreshCw, Trash2, Camera } from 'lucide-react';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';
import styles from './HelpImageUploader.module.css';

export default function HelpImageUploader() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [image, setImage] = useState<{
    file: File | null;
    preview: string | null;
  }>({ file: null, preview: null });

  const [isUploading, setIsUploading] = useState(false);

  const handleFileSelect = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  const handleCameraCapture = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  const handleFileChange = useCallback(
    async (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (!file) return;

      // Валидация
      if (!file.type.startsWith('image/')) {
        return;
      }
      if (file.size > 10 * 1024 * 1024) {
        return;
      }

      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      analytics.trackEvent('help_choose_file_clicked', {
        child_profile_id: childProfileId,
      });

      // Конвертируем файл в base64 для надёжного хранения
      const reader = new FileReader();
      reader.onload = () => {
        const base64 = reader.result as string;
        setImage({ file, preview: base64 });
      };
      reader.readAsDataURL(file);

      // Сбрасываем input
      e.target.value = '';
    },
    [analytics]
  );

  const handleReplace = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  const handleDelete = useCallback(() => {
    setImage({ file: null, preview: null });
  }, []);

  const handleContinue = useCallback(async () => {
    if (!image.file && !image.preview) return;
    if (isUploading) return;

    setIsUploading(true);

    try {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      if (!childProfileId) {
        setIsUploading(false);
        return;
      }

      analytics.trackEvent('help_upload_retry_clicked', {
        child_profile_id: childProfileId,
      });

      // 1. Создаём попытку
      const attempt = await helpAPI.createAttempt(childProfileId);

      // 2. Загружаем изображение
      await helpAPI.uploadImage(attempt.id, image.file!);

      // Сохраняем в sessionStorage для ProcessingPage
      sessionStorage.setItem('help_photo_data', JSON.stringify({
        fileName: image.file?.name || 'photo.jpg',
        fileType: image.file?.type || 'image/jpeg',
        fileSize: image.file?.size || 0,
        base64: image.preview || '',
      }));

      // 3. Переходим на страницу обработки с attemptId
      navigate(ROUTES.HELP_PROCESSING, {
        state: {
          attemptId: attempt.id,
          imageUrl: image.preview,
        }
      });
    } catch (error) {
      console.error('[HelpImageUploader] Failed to upload:', error);
      setIsUploading(false);
      alert('Не удалось загрузить изображение. Попробуй ещё раз.');
    }
  }, [image, navigate, analytics, isUploading]);

  const isUploaded = image.file !== null || image.preview !== null;

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
        onClick={() => navigate(ROUTES.HOME)}
        className={styles.backButton}
      >
        <ArrowLeft size={20} />
        <span className={styles.backText}>Назад</span>
      </button>

      {/* Заголовок */}
      <h2 className={styles.title}>Загрузи фото</h2>
      <p className={styles.subtitle}>Загрузи фото задания из учебника</p>

      {/* Слот изображения */}
      <div className={styles.slotCard}>
        {/* Лейбл + кнопки действий */}
        <div className={styles.slotHeader}>
          <span className={styles.slotLabel}>Задание</span>
          {isUploaded && (
            <div className={styles.slotActions}>
              <button
                onClick={handleReplace}
                className={`${styles.slotActionButton} ${styles.slotActionButtonReplace}`}
                disabled={isUploading}
              >
                <RefreshCw size={14} />
              </button>
              <button
                onClick={handleDelete}
                className={`${styles.slotActionButton} ${styles.slotActionButtonDelete}`}
                disabled={isUploading}
              >
                <Trash2 size={14} />
              </button>
            </div>
          )}
        </div>

        {/* Область превью или загрузки */}
        {image.preview ? (
          <div className={styles.previewContainer}>
            <img
              src={image.preview}
              alt="Задание"
              className={styles.previewImage}
            />
          </div>
        ) : (
          <div className={styles.uploadArea}>
            <div className={styles.uploadButtons}>
              <button
                onClick={handleFileSelect}
                className={styles.uploadButton}
                disabled={isUploading}
              >
                <Plus size={20} />
                <span>Выбрать</span>
              </button>
              <button
                onClick={handleCameraCapture}
                className={`${styles.uploadButton} ${styles.uploadButtonCamera}`}
                disabled={isUploading}
              >
                <Camera size={16} />
                <span>Камера</span>
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Кнопка продолжить */}
      <button
        onClick={handleContinue}
        disabled={!isUploaded || isUploading}
        className={`${styles.continueButton} ${
          isUploaded && !isUploading ? styles.continueButtonEnabled : styles.continueButtonDisabled
        }`}
      >
        {isUploading ? 'Загрузка...' : 'Продолжить'}
      </button>
    </div>
  );
}
