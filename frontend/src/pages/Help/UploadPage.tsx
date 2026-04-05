// src/pages/Help/UploadPage.tsx
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { ArrowLeft, CheckCircle, XCircle } from 'lucide-react';
import { motion } from 'framer-motion';
import { Spinner } from '@/components/ui/Spinner';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';
import imageCompression from 'browser-image-compression';
import styles from './UploadPage.module.css';

type UploadStatus = 'compressing' | 'uploading' | 'success' | 'error';

export default function UploadPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();

  const [childProfileId, setChildProfileId] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [status, setStatus] = useState<UploadStatus>('compressing');
  const [errorMessage, setErrorMessage] = useState<string>('');

  const file = location.state?.file as File;
  const source = location.state?.source as string;

  useEffect(() => {
    const loadProfile = async () => {
      const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      setChildProfileId(profileId);
    };

    loadProfile();
  }, []);

  useEffect(() => {
    if (!file) {
      navigate(ROUTES.HELP);
      return;
    }

    if (childProfileId) {
      uploadImage();
    }
  }, [file, childProfileId]);

  const uploadImage = async () => {
    if (!childProfileId) return;

    try {
      // 1. Create attempt
      const attempt = await helpAPI.createAttempt(childProfileId);

      analytics.trackEvent('help_image_upload_started', {
        child_profile_id: childProfileId,
        upload_source: source,
        attempt_id: attempt.id,
      });

      // 2. Compress image
      setStatus('compressing');
      const maxSizeMB = 7; // 7 МБ, чтобы после base64 (~9.3 МБ) не превысить лимит бэкенда 10 МБ
      const quality = 0.8;

      const compressedFile = await imageCompression(file, {
        maxSizeMB,
        maxWidthOrHeight: 1920,
        useWebWorker: true,
        initialQuality: quality,
      });

      // 3. Upload image
      setStatus('uploading');
      const result = await helpAPI.uploadImage(
        attempt.id,
        compressedFile,
        (progress) => setUploadProgress(progress)
      );

      setStatus('success');

      analytics.trackEvent('help_image_upload_completed', {
        child_profile_id: childProfileId,
        attempt_id: attempt.id,
      });

      // Navigate to processing
      setTimeout(() => {
        navigate(ROUTES.HELP_PROCESSING, {
          state: {
            attemptId: attempt.id,
            imageUrl: result.imageUrl,
          },
        });
      }, 500);
    } catch (error) {
      console.error('[UploadPage] Upload failed:', error);
      setStatus('error');
      setErrorMessage(error instanceof Error ? error.message : 'Неизвестная ошибка');

      if (childProfileId) {
        analytics.trackEvent('help_image_upload_failed', {
          child_profile_id: childProfileId,
          error_message: error instanceof Error ? error.message : 'Unknown error',
        });
      }
    }
  };

  const handleRetry = () => {
    if (childProfileId) {
      analytics.trackEvent('help_upload_retry_clicked', {
        child_profile_id: childProfileId,
      });
    }
    navigate(ROUTES.HELP);
  };

  const handleBack = () => {
    if (childProfileId) {
      analytics.trackEvent('help_upload_back_clicked', {
        child_profile_id: childProfileId,
      });
    }
    navigate(ROUTES.HELP);
  };

  const getStatusMessage = () => {
    switch (status) {
      case 'compressing':
        return 'Сжимаем изображение...';
      case 'uploading':
        return 'Загружаем изображение...';
      case 'success':
        return 'Изображение загружено!';
      case 'error':
        return 'Ошибка загрузки';
      default:
        return '';
    }
  };

  return (
    <div className={styles.container}>
      {/* Header */}
      <button onClick={handleBack} className={styles.backButton}>
        <ArrowLeft size={20} />
        <span>Назад</span>
      </button>

      {/* Content */}
      <div className={styles.content}>
        {status === 'compressing' && (
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className={styles.statusContainer}
          >
            <Spinner size="lg" />
            <h2 className={styles.title}>{getStatusMessage()}</h2>
            <p className={styles.subtitle}>Это займет несколько секунд</p>
          </motion.div>
        )}

        {status === 'uploading' && (
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className={styles.statusContainer}
          >
            <h2 className={styles.title}>{getStatusMessage()}</h2>
            <div className={styles.progressBarContainer}>
              <div className={styles.progressBar}>
                <div
                  className={styles.progressFill}
                  style={{ width: `${uploadProgress}%` }}
                />
              </div>
              <span className={styles.progressLabel}>{uploadProgress}%</span>
            </div>
          </motion.div>
        )}

        {status === 'success' && (
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className={styles.statusContainer}
          >
            <div className={styles.iconCircle}>
              <CheckCircle size={48} className={styles.successIcon} />
            </div>
            <h2 className={styles.title}>{getStatusMessage()}</h2>
          </motion.div>
        )}

        {status === 'error' && (
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className={styles.statusContainer}
          >
            <div className={styles.iconCircle}>
              <XCircle size={48} className={styles.errorIcon} />
            </div>
            <h2 className={styles.title}>{getStatusMessage()}</h2>
            <p className={styles.errorMessage}>{errorMessage || 'Попробуйте еще раз'}</p>
            <button onClick={handleRetry} className={styles.retryButton}>
              Повторить
            </button>
          </motion.div>
        )}
      </div>
    </div>
  );
}
