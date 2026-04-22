// src/pages/Help/SourcePicker.tsx
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Camera, Image as ImageIcon, AlertCircle } from 'lucide-react';
import { motion } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { usePlatform } from '@/hooks/usePlatform';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { ROUTES } from '@/config/routes';
import type { UploadSource } from '@/types/help';
import styles from './SourcePicker.module.css';

export default function SourcePicker() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const { isDesktop } = usePlatform();
  const [childProfileId, setChildProfileId] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isMobile = !isDesktop;

  useEffect(() => {
    const loadProfile = async () => {
      const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      setChildProfileId(profileId);

      if (profileId) {
        analytics.trackEvent('help_source_picker_opened', {
          child_profile_id: profileId,
        });
      }
    };

    loadProfile();
  }, [analytics]);

  const handleImageSelected = (file: File, source: UploadSource) => {
    // Валидация типа файла
    if (!file.type.startsWith('image/')) {
      setError('Пожалуйста, выберите изображение (JPG, PNG)');
      return;
    }

    // Валидация размера (макс 10MB для загрузки, будет сжато до 7MB)
    const maxSize = 10 * 1024 * 1024;
    if (file.size > maxSize) {
      setError('Размер файла не должен превышать 10 МБ');
      return;
    }

    setError(null);

    if (childProfileId) {
      analytics.trackEvent('help_image_selected', {
        child_profile_id: childProfileId,
        source,
        file_size: file.size,
        mime_type: file.type,
      });
    }

    navigate('/help/upload-progress', {
      state: { file, source },
    });
  };

  const handleFileSelect = () => {
    setError(null);

    if (childProfileId) {
      analytics.trackEvent('help_choose_file_clicked', {
        child_profile_id: childProfileId,
      });
    }

    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/jpeg,image/png,image/jpg';
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        handleImageSelected(file, 'file');
      }
    };
    input.click();
  };

  const handleCameraClick = () => {
    setError(null);

    if (childProfileId) {
      analytics.trackEvent('help_camera_clicked', {
        child_profile_id: childProfileId,
      });
    }

    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.setAttribute('capture', 'environment');
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        handleImageSelected(file, 'camera');
      }
    };
    input.click();
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const file = e.dataTransfer.files[0];
    if (file) {
      if (!file.type.startsWith('image/')) {
        setError('Пожалуйста, выберите изображение (JPG, PNG)');
        return;
      }
      handleImageSelected(file, 'dragdrop');
    }
  };

  return (
    <div className={styles.container}>
      {/* Header */}
      <button onClick={() => navigate(ROUTES.HOME)} className={styles.backButton}>
        <ArrowLeft size={20} />
        <span>Назад</span>
      </button>

      {/* Title Section */}
      <div className={styles.titleSection}>
        <div className={styles.iconCircle}>
          <ImageIcon size={30} className={styles.iconPrimary} />
        </div>
        <h1 className={styles.title}>Помоги разобраться</h1>
        <p className={styles.subtitle}>Загрузи фото задания из учебника</p>
      </div>

      {/* Upload Options */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className={styles.optionsContainer}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        {/* File Select Button */}
        <button
          onClick={handleFileSelect}
          className={`${styles.uploadButton} ${isDragging ? styles.dragging : ''}`}
        >
          <div className={styles.buttonIcon}>
            <ImageIcon size={24} className={styles.iconPrimary} />
          </div>
          <div className={styles.buttonText}>
            <p className={styles.buttonTitle}>
              {isDragging ? 'Отпустите файл' : 'Выбрать изображение'}
            </p>
            <p className={styles.buttonSubtitle}>
              {isDragging ? 'Загрузка начнётся автоматически' : 'JPG, PNG (до 10 МБ)'}
            </p>
          </div>
        </button>

        {/* Camera Button - только на мобильных */}
        {isMobile && (
          <button onClick={handleCameraClick} className={styles.uploadButton}>
            <div className={`${styles.buttonIcon} ${styles.cameraIcon}`}>
              <Camera size={24} className={styles.iconCamera} />
            </div>
            <div className={styles.buttonText}>
              <p className={styles.buttonTitle}>Сфотографировать</p>
              <p className={styles.buttonSubtitle}>Открыть камеру</p>
            </div>
          </button>
        )}

        {/* Drag & Drop hint - только на десктопе */}
        {!isMobile && (
          <p className={styles.dragHint}>
            или перетащите файл в любую область экрана
          </p>
        )}

        {/* Error Message */}
        {error && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className={styles.errorBox}
          >
            <AlertCircle size={20} className={styles.errorIcon} />
            <div className={styles.errorContent}>
              <p className={styles.errorText}>{error}</p>
              <button onClick={() => setError(null)} className={styles.errorDismiss}>
                Понятно
              </button>
            </div>
          </motion.div>
        )}
      </motion.div>
    </div>
  );
}
