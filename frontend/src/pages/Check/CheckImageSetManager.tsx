// src/pages/Check/CheckImageSetManager.tsx
import { useState, useRef, useCallback } from 'react';
import { useNavigate, useSearchParams, useLocation } from 'react-router-dom';
import {
  ArrowLeft,
  Plus,
  Trash2,
  ArrowUpDown,
  RefreshCw,
  Camera,
} from 'lucide-react';
import { motion } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import styles from './CheckImageSetManager.module.css';

interface ImageSlot {
  id: 'task' | 'answer';
  label: string;
  file: File | null;
  preview: string | null;
}

const INITIAL_IMAGES_TWO_PHOTO: ImageSlot[] = [
  { id: 'task', label: 'Задание', file: null, preview: null },
  { id: 'answer', label: 'Ответ', file: null, preview: null },
];

const INITIAL_IMAGES_SINGLE_PHOTO: ImageSlot[] = [
  { id: 'task', label: 'Задание и ответ', file: null, preview: null },
];

export default function CheckImageSetManager() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const analytics = useAnalytics();
  const fileInputRefs = useRef<{ task: HTMLInputElement | null; answer: HTMLInputElement | null }>({
    task: null,
    answer: null,
  });

  const scenario = searchParams.get('scenario') || 'two_photo';
  const isSinglePhoto = scenario === 'single_photo';

  // Получаем данные для режима исправления ошибок
  let existingTaskImage = location.state?.existingTaskImage as string | null;

  // Если нет в state, пробуем получить из sessionStorage
  if (!existingTaskImage) {
    const lastTaskImageStr = sessionStorage.getItem('check_last_task_image');
    if (lastTaskImageStr) {
      try {
        const lastTaskImage = JSON.parse(lastTaskImageStr);
        // Используем только если сценарий совпадает
        if (lastTaskImage.scenario === scenario) {
          existingTaskImage = lastTaskImage.base64;
        }
      } catch (e) {
        console.error('Failed to parse last task image:', e);
      }
    }
  }

  // Получаем данные для режима "из помощи" - переход с страницы подсказок
  const fromHelpMode = location.state?.mode === 'from_help';
  let helpTaskImage = location.state?.taskImage as string | null;

  const fixErrorsMode = location.state?.mode === 'fix_errors' || !!existingTaskImage;

  // Инициализируем слоты в зависимости от сценария
  const [images, setImages] = useState<ImageSlot[]>(() => {
    if (isSinglePhoto) {
      // Для 1 фото в режиме "из помощи" предзаполняем задание
      if (fromHelpMode && helpTaskImage) {
        return [
          {
            id: 'task',
            label: 'Задание и ответ',
            file: null,
            preview: helpTaskImage // Существующее фото задания
          },
        ];
      }
      return INITIAL_IMAGES_SINGLE_PHOTO;
    }

    // Для 2 фото в режиме "из помощи" предзаполняем задание (ПРОВЕРЯЕМ ПЕРЕД fixErrorsMode!)
    if (fromHelpMode && helpTaskImage) {
      return [
        {
          id: 'task',
          label: 'Задание',
          file: null,
          preview: helpTaskImage // Существующее фото задания
        },
        {
          id: 'answer',
          label: 'Ответ',
          file: null,
          preview: null
        },
      ];
    }

    // Для 2 фото в режиме исправления ошибок предзаполняем задание
    if (fixErrorsMode && existingTaskImage) {
      return [
        {
          id: 'task',
          label: 'Задание',
          file: null,
          preview: existingTaskImage // Существующее фото задания
        },
        {
          id: 'answer',
          label: 'Ответ (исправленный)',
          file: null,
          preview: null
        },
      ];
    }

    return INITIAL_IMAGES_TWO_PHOTO;
  });
  const [isReordered, setIsReordered] = useState(false);

  const allUploaded = images.every((img) => img.file !== null || ((fixErrorsMode || fromHelpMode) && img.id === 'task' && img.preview));

  const displayedImages = isReordered ? [...images].reverse() : images;

  const handleFileSelect = useCallback(
    (slotId: 'task' | 'answer') => {
      fileInputRefs.current[slotId]?.click();
    },
    []
  );

  const handleCameraCapture = useCallback(
    (slotId: 'task' | 'answer') => {
      // Для мобильных — можно открыть getUserMedia
      // Пока используем тот же input с capture
      fileInputRefs.current[slotId]?.click();
    },
    []
  );

  const handleFileChange = useCallback(
    async (slotId: 'task' | 'answer', e: React.ChangeEvent<HTMLInputElement>) => {
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
      analytics.trackEvent('check_image_uploaded', {
        child_profile_id: childProfileId,
        scenario,
        image_role: slotId,
      });

      // Конвертируем файл в base64 для надёжного хранения
      const reader = new FileReader();
      reader.onload = () => {
        const base64 = reader.result as string;
        
        setImages((prev) =>
          prev.map((img) =>
            img.id === slotId ? { ...img, file, preview: base64 } : img
          )
        );
      };
      reader.readAsDataURL(file);

      // Сбрасываем input
      e.target.value = '';
    },
    [analytics, scenario]
  );

  const handleReplace = useCallback(
    (slotId: 'task' | 'answer') => {
      handleFileSelect(slotId);
    },
    [handleFileSelect]
  );

  const handleDelete = useCallback(
    (slotId: 'task' | 'answer') => {
      const childProfileId = vkStorage.getItem(storageKeys.PROFILE_ID);
      analytics.trackEvent('check_image_deleted', {
        child_profile_id: childProfileId,
        scenario,
        image_role: slotId,
      });

      setImages((prev) =>
        prev.map((img) =>
          img.id === slotId ? { ...img, file: null, preview: null } : img
        )
      );
    },
    [analytics, scenario]
  );

  const handleReorder = useCallback(() => {
    setIsReordered((prev) => !prev);
  }, []);

  const handleContinue = useCallback(async () => {
    if (!allUploaded) return;

    const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);

    analytics.trackEvent('check_images_upload_completed', {
      child_profile_id: childProfileId,
      scenario,
    });

    const taskImage = images.find((img) => img.id === 'task');

    if (isSinglePhoto) {
      // Для 1 фото сохраняем в один ключ и переходим на quality-single
      if (!taskImage?.file && !taskImage?.preview) return;

      sessionStorage.setItem('check_single_photo_data', JSON.stringify({
        fileName: taskImage.file?.name || 'existing.jpg',
        fileType: taskImage.file?.type || 'image/jpeg',
        fileSize: taskImage.file?.size || 0,
        base64: taskImage.preview || '',
      }));

      // Сохраняем также для режима исправления ошибок
      sessionStorage.setItem('check_last_task_image', JSON.stringify({
        base64: taskImage.preview || '',
        scenario: 'single_photo',
      }));

      navigate('/check/quality-single');
    } else {
      // Для 2 фото сохраняем оба и переходим на quality-two
      const answerImage = images.find((img) => img.id === 'answer');

      // В режиме исправления ошибок используем существующее фото задания
      if (fixErrorsMode && existingTaskImage) {
        if (!answerImage?.file) return;

        // Сохраняем существующее задание и новый ответ
        sessionStorage.setItem('check_task_photo', JSON.stringify({
          fileName: 'existing_task.jpg',
          fileType: 'image/jpeg',
          fileSize: 0,
          base64: existingTaskImage,
        }));

        sessionStorage.setItem('check_answer_photo', JSON.stringify({
          fileName: answerImage.file.name,
          fileType: answerImage.file.type,
          fileSize: answerImage.file.size,
          base64: answerImage.preview,
        }));
      } else if (fromHelpMode && helpTaskImage) {
        // В режиме "из помощи" используем фото задания из помощи
        if (!answerImage?.file) return;

        // Сохраняем задание из помощи и новый ответ
        sessionStorage.setItem('check_task_photo', JSON.stringify({
          fileName: 'help_task.jpg',
          fileType: 'image/jpeg',
          fileSize: 0,
          base64: helpTaskImage,
        }));

        sessionStorage.setItem('check_answer_photo', JSON.stringify({
          fileName: answerImage.file.name,
          fileType: answerImage.file.type,
          fileSize: answerImage.file.size,
          base64: answerImage.preview,
        }));

        // Сохраняем изображение задания для возможного исправления ошибок
        sessionStorage.setItem('check_last_task_image', JSON.stringify({
          base64: helpTaskImage,
          scenario: 'two_photo',
        }));
      } else {
        // Обычный режим - оба фото новые
        if (!taskImage?.file || !answerImage?.file) return;

        const storeImage = (img: ImageSlot, key: string) => {
          if (img.preview) {
            sessionStorage.setItem(key, JSON.stringify({
              fileName: img.file!.name,
              fileType: img.file!.type,
              fileSize: img.file!.size,
              base64: img.preview,
            }));
          }
        };

        storeImage(taskImage, 'check_task_photo');
        storeImage(answerImage, 'check_answer_photo');

        // Сохраняем изображение задания для возможного исправления ошибок
        sessionStorage.setItem('check_last_task_image', JSON.stringify({
          base64: taskImage.preview,
          scenario: 'two_photo',
        }));
      }

      navigate('/check/quality-two');
    }
  }, [allUploaded, images, navigate, analytics, scenario, isSinglePhoto, fixErrorsMode, existingTaskImage, fromHelpMode, helpTaskImage]);

  // Cleanup при размонтировании (теперь не нужен для blob URLs, но оставляем для совместимости)
  useState(() => {
    return () => {
      // Больше не нужно отзывать blob URLs, так как используем base64
    };
  });

  return (
    <div className={styles.container}>
      {/* Скрытые inputs для каждого слота */}
      <input
        ref={(el) => { fileInputRefs.current.task = el; }}
        type="file"
        accept="image/*"
        onChange={(e) => handleFileChange('task', e)}
        className={styles.hiddenInput}
      />
      {!isSinglePhoto && (
        <input
          ref={(el) => { fileInputRefs.current.answer = el; }}
          type="file"
          accept="image/*"
          onChange={(e) => handleFileChange('answer', e)}
          className={styles.hiddenInput}
        />
      )}

      {/* Кнопка назад */}
      <button
        onClick={() => navigate(-1)}
        className={styles.backButton}
      >
        <ArrowLeft size={20} />
        <span className={styles.backText}>Назад</span>
      </button>

      {/* Заголовок */}
      <h2 className={styles.title}>
        {isSinglePhoto ? 'Загрузи фото' : 'Изображения'}
      </h2>
      <p className={styles.subtitle}>
        {isSinglePhoto 
          ? 'Загрузи фото с заданием и решением' 
          : 'Загрузи все необходимые фото'}
      </p>

      {/* Слоты изображений */}
      <div className={styles.slots}>
        {displayedImages.map((img, index) => (
          <motion.div
            key={img.id}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
            className={styles.slotCard}
          >
            {/* Лейбл + кнопки действий */}
            <div className={styles.slotHeader}>
              <span className={styles.slotLabel}>{img.label}</span>
              {(img.file || (fixErrorsMode && img.id === 'task' && img.preview)) && (
                <div className={styles.slotActions}>
                  <button
                    onClick={() => handleReplace(img.id)}
                    className={`${styles.slotActionButton} ${styles.slotActionButtonReplace}`}
                  >
                    <RefreshCw size={14} />
                  </button>
                  {!fixErrorsMode && !fromHelpMode && (
                    <button
                      onClick={() => handleDelete(img.id)}
                      className={`${styles.slotActionButton} ${styles.slotActionButtonDelete}`}
                    >
                      <Trash2 size={14} />
                    </button>
                  )}
                </div>
              )}
            </div>

            {/* Область превью или загрузки */}
            {img.preview ? (
              <div className={styles.previewContainer}>
                <img
                  src={img.preview}
                  alt={img.label}
                  className={styles.previewImage}
                />
                {fixErrorsMode && img.id === 'task' && (
                  <div className={styles.existingImageOverlay}>
                    <span>✓ Задание уже загружено</span>
                  </div>
                )}
                {fromHelpMode && img.id === 'task' && (
                  <div className={styles.existingImageOverlay}>
                    <span>✓ Задание уже загружено</span>
                  </div>
                )}
              </div>
            ) : (
              <div className={styles.uploadArea}>
                <div className={styles.uploadButtons}>
                  <button
                    onClick={() => handleFileSelect(img.id)}
                    className={styles.uploadButton}
                  >
                    <Plus size={20} />
                    <span>Выбрать</span>
                  </button>
                  <button
                    onClick={() => handleCameraCapture(img.id)}
                    className={`${styles.uploadButton} ${styles.uploadButtonCamera}`}
                  >
                    <Camera size={16} />
                    <span>Камера</span>
                  </button>
                </div>
              </div>
            )}
          </motion.div>
        ))}
      </div>

      {/* Кнопка смены порядка (только для 2 фото) */}
      {!isSinglePhoto && (
        <button
          onClick={handleReorder}
          className={styles.reorderButton}
        >
          <ArrowUpDown size={16} className={styles.reorderIcon} />
          Поменять порядок
        </button>
      )}

      {/* Кнопка продолжить */}
      <button
        onClick={handleContinue}
        disabled={!allUploaded}
        className={`${styles.continueButton} ${
          allUploaded ? styles.continueButtonEnabled : styles.continueButtonDisabled
        }`}
      >
        Продолжить
      </button>
    </div>
  );
}
