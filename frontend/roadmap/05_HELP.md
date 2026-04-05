# Phase 4: Поток "Помоги разобраться" (Help Flow)

**Длительность:** 5-7 дней
**Приоритет:** Критический
**Зависимости:** 02_CORE.md, 04_HOME.md

---

## Цель

Реализовать полный поток помощи с домашним заданием: выбор источника изображения, загрузка, проверка качества, обработка, получение подсказок и отправка ответа.

---

## Архитектура потока

### Этапы Help Flow

```
1. Source Selection   → Выбор источника (camera, file, clipboard, drag&drop)
2. Image Upload       → Загрузка изображения с прогрессом
3. Quality Check      → Проверка качества фото
4. Crop (optional)    → Обрезка изображения
5. Processing         → Обработка задания (loading, long-wait, save-and-wait)
6. Result & Hints     → Показ результата с подсказками
7. Answer Submit      → Отправка ответа пользователя
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/help.ts`

```typescript
export type UploadSource = 'camera' | 'file' | 'clipboard' | 'dragdrop';

export interface HelpAttempt {
  id: string;
  childProfileId: string;
  mode: 'help';
  status: HelpStatus;
  imageUrl?: string;
  thumbnailUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export type HelpStatus =
  | 'created'
  | 'uploading'
  | 'uploaded'
  | 'quality_check'
  | 'processing'
  | 'long_wait'
  | 'completed'
  | 'failed';

export interface Hint {
  id: string;
  level: 1 | 2 | 3; // Уровень подсказки
  title: string;
  content: string;
  order: number;
}

export interface HelpResult {
  attemptId: string;
  hints: Hint[];
  coinsEarned: number;
  damageDealt: number; // Урон злодею
}

export interface CropArea {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface UploadProgress {
  loaded: number;
  total: number;
  percentage: number;
}
```

---

## Часть 2: API Integration

### 2.1. Help API

**Файл:** `src/api/help.ts`

```typescript
import { apiClient } from './client';
import type {
  HelpAttempt,
  HelpResult,
  Hint,
  CropArea,
} from '@/types/help';

export const helpAPI = {
  /**
   * Создать новую попытку помощи
   */
  async createAttempt(childProfileId: string): Promise<HelpAttempt> {
    return apiClient.post<HelpAttempt>('/attempts', {
      childProfileId,
      mode: 'help',
    });
  },

  /**
   * Загрузить изображение задания
   */
  async uploadImage(
    attemptId: string,
    file: Blob,
    onProgress?: (progress: number) => void
  ): Promise<{ imageUrl: string; thumbnailUrl: string }> {
    return apiClient.upload<{ imageUrl: string; thumbnailUrl: string }>(
      `/attempts/${attemptId}/images`,
      file,
      onProgress
    );
  },

  /**
   * Подтвердить качество изображения
   */
  async confirmQuality(attemptId: string): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/confirm-quality`);
  },

  /**
   * Обрезать изображение
   */
  async cropImage(attemptId: string, cropArea: CropArea): Promise<{ imageUrl: string }> {
    return apiClient.post<{ imageUrl: string }>(
      `/attempts/${attemptId}/crop`,
      { cropArea }
    );
  },

  /**
   * Начать обработку задания
   */
  async processAttempt(attemptId: string): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/process`);
  },

  /**
   * Получить результат обработки (polling)
   */
  async getResult(attemptId: string): Promise<HelpResult> {
    return apiClient.get<HelpResult>(`/attempts/${attemptId}/result`);
  },

  /**
   * Получить следующую подсказку
   */
  async getNextHint(attemptId: string, currentLevel: number): Promise<Hint> {
    return apiClient.post<Hint>(`/attempts/${attemptId}/next-hint`, {
      currentLevel,
    });
  },

  /**
   * Отправить ответ пользователя
   */
  async submitAnswer(
    attemptId: string,
    answer: string
  ): Promise<{ success: boolean; coinsEarned: number }> {
    return apiClient.post<{ success: boolean; coinsEarned: number }>(
      `/attempts/${attemptId}/submit-answer`,
      { answer }
    );
  },

  /**
   * Сохранить и подождать (для long-wait)
   */
  async saveAndWait(attemptId: string): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/save-and-wait`);
  },

  /**
   * Отменить попытку
   */
  async cancelAttempt(attemptId: string): Promise<void> {
    return apiClient.delete<void>(`/attempts/${attemptId}`);
  },
};
```

---

## Часть 3: Компоненты

### 3.1. HelpFlow Container

**Файл:** `src/pages/Help/HelpFlow.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/config/routes';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import type { HelpAttempt, HelpStatus } from '@/types/help';

export default function HelpFlow() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const [attemptId, setAttemptId] = useState<string | null>(null);

  useEffect(() => {
    // Analytics: help_flow_started
    analytics.trackEvent('help_flow_started', {
      child_profile_id: profile?.id,
      mode: 'help',
      entry_point: 'home_button',
    });
  }, []);

  // Redirect to upload screen
  useEffect(() => {
    navigate(ROUTES.HELP_UPLOAD);
  }, []);

  return null;
}
```

---

### 3.2. SourcePicker Component

**Файл:** `src/pages/Help/SourcePicker.tsx`

```typescript
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Button } from '@/components/ui/Button';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { usePlatform } from '@/hooks/usePlatform';
import { ROUTES } from '@/config/routes';
import styles from './SourcePicker.module.css';

export default function SourcePicker() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const { platform } = usePlatform();

  useEffect(() => {
    // Analytics: help_source_picker_opened
    analytics.trackEvent('help_source_picker_opened', {
      child_profile_id: profile?.id,
    });
  }, []);

  const handleFileSelect = () => {
    // Analytics: help_choose_file_clicked
    analytics.trackEvent('help_choose_file_clicked', {
      child_profile_id: profile?.id,
    });

    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/jpeg,image/png';
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        handleImageSelected(file, 'file');
      }
    };
    input.click();
  };

  const handleCameraClick = async () => {
    // Analytics: help_camera_clicked
    analytics.trackEvent('help_camera_clicked', {
      child_profile_id: profile?.id,
    });

    try {
      // Request camera access via platform bridge
      const hasAccess = await platform.requestPhotoAccess();
      if (hasAccess) {
        // Open camera
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = 'image/*';
        input.capture = 'environment';
        input.onchange = (e) => {
          const file = (e.target as HTMLInputElement).files?.[0];
          if (file) {
            handleImageSelected(file, 'camera');
          }
        };
        input.click();
      }
    } catch (error) {
      console.error('[SourcePicker] Camera access failed:', error);
    }
  };

  const handleClipboardPaste = async () => {
    // Analytics: help_clipboard_clicked
    analytics.trackEvent('help_clipboard_clicked', {
      child_profile_id: profile?.id,
    });

    try {
      const clipboardItems = await navigator.clipboard.read();
      for (const item of clipboardItems) {
        const imageType = item.types.find((type) => type.startsWith('image/'));
        if (imageType) {
          const blob = await item.getType(imageType);
          handleImageSelected(blob, 'clipboard');
          break;
        }
      }
    } catch (error) {
      console.error('[SourcePicker] Clipboard read failed:', error);
    }
  };

  const handleDragDrop = (file: File) => {
    // Analytics: help_dragdrop_used
    analytics.trackEvent('help_dragdrop_used', {
      child_profile_id: profile?.id,
    });

    handleImageSelected(file, 'dragdrop');
  };

  const handleImageSelected = (file: Blob, source: string) => {
    // Analytics: help_image_selected
    analytics.trackEvent('help_image_selected', {
      child_profile_id: profile?.id,
      upload_source: source,
      file_size_bucket: getFileSizeBucket(file.size),
      mime_type: file.type,
    });

    // Navigate to upload with file
    navigate(ROUTES.HELP_UPLOAD, {
      state: { file, source },
    });
  };

  return (
    <Container className={styles.container}>
      <div className={styles.header}>
        <button
          className={styles.backButton}
          onClick={() => navigate(ROUTES.HOME)}
        >
          Назад
        </button>
        <h1 className={styles.title}>Помоги разобраться</h1>
        <p className={styles.subtitle}>
          Загрузи фото задания из учебника
        </p>
      </div>

      <div className={styles.options}>
        <Button
          variant="outline"
          size="lg"
          isFullWidth
          onClick={handleFileSelect}
          leftIcon="📁"
        >
          <div className={styles.optionContent}>
            <span className={styles.optionTitle}>Выбрать изображение</span>
            <span className={styles.optionSubtitle}>JPG, PNG</span>
          </div>
        </Button>

        <Button
          variant="outline"
          size="lg"
          isFullWidth
          onClick={handleCameraClick}
          leftIcon="📷"
        >
          <div className={styles.optionContent}>
            <span className={styles.optionTitle}>Сфотографировать</span>
            <span className={styles.optionSubtitle}>Открыть камеру</span>
          </div>
        </Button>

        <Button
          variant="outline"
          size="lg"
          isFullWidth
          onClick={handleClipboardPaste}
          leftIcon="📋"
        >
          <div className={styles.optionContent}>
            <span className={styles.optionTitle}>Вставить из буфера</span>
            <span className={styles.optionSubtitle}>Ctrl+V / ⌘+V</span>
          </div>
        </Button>
      </div>

      <DropZone onDrop={handleDragDrop} />
    </Container>
  );
}

function getFileSizeBucket(size: number): string {
  if (size < 1024 * 1024) return '<1MB';
  if (size < 5 * 1024 * 1024) return '1-5MB';
  if (size < 10 * 1024 * 1024) return '5-10MB';
  return '>10MB';
}

// DropZone component for drag&drop
interface DropZoneProps {
  onDrop: (file: File) => void;
}

function DropZone({ onDrop }: DropZoneProps) {
  const [isDragging, setIsDragging] = useState(false);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const file = e.dataTransfer.files[0];
    if (file && file.type.startsWith('image/')) {
      onDrop(file);
    }
  };

  return (
    <div
      className={clsx(styles.dropZone, {
        [styles.dragging]: isDragging,
      })}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      <p>Или перетащи файл сюда</p>
    </div>
  );
}
```

---

### 3.3. ImageUpload Component

**Файл:** `src/pages/Help/ImageUpload.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { ProgressBar } from '@/components/ui/ProgressBar';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';
import imageCompression from 'browser-image-compression';
import styles from './ImageUpload.module.css';

export default function ImageUpload() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const [uploadProgress, setUploadProgress] = useState(0);
  const [status, setStatus] = useState<'compressing' | 'uploading' | 'success' | 'error'>('compressing');
  const [attemptId, setAttemptId] = useState<string | null>(null);

  const file = location.state?.file as Blob;
  const source = location.state?.source as string;

  useEffect(() => {
    if (!file || !profile?.id) {
      navigate(ROUTES.HELP_UPLOAD);
      return;
    }

    uploadImage();
  }, [file, profile]);

  const uploadImage = async () => {
    try {
      // 1. Create attempt
      const attempt = await helpAPI.createAttempt(profile!.id);
      setAttemptId(attempt.id);

      // Analytics: help_image_upload_started
      analytics.trackEvent('help_image_upload_started', {
        child_profile_id: profile!.id,
        upload_source: source,
        attempt_id: attempt.id,
      });

      // 2. Compress image
      setStatus('compressing');
      const maxSizeMB = Number(import.meta.env.VITE_MAX_IMAGE_SIZE_MB) || 10;
      const quality = Number(import.meta.env.VITE_IMAGE_COMPRESSION_QUALITY) || 0.8;

      const compressedFile = await imageCompression(file as File, {
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

      // Analytics: help_image_upload_completed (sent by backend)
      // Backend sends this event when upload succeeds

      setStatus('success');

      // Navigate to quality check
      setTimeout(() => {
        navigate(ROUTES.HELP_QUALITY, {
          state: {
            attemptId: attempt.id,
            imageUrl: result.imageUrl,
            thumbnailUrl: result.thumbnailUrl,
          },
        });
      }, 500);
    } catch (error) {
      console.error('[ImageUpload] Upload failed:', error);
      setStatus('error');

      // Analytics: help_image_upload_failed (sent by backend)
      // Backend sends this event when upload fails
    }
  };

  const handleRetry = () => {
    navigate(ROUTES.HELP_UPLOAD);
  };

  return (
    <Container className={styles.container}>
      <div className={styles.content}>
        {status === 'compressing' && (
          <>
            <div className={styles.spinner} />
            <h2>Сжимаем изображение...</h2>
            <p>Это займет несколько секунд</p>
          </>
        )}

        {status === 'uploading' && (
          <>
            <h2>Загружаем изображение...</h2>
            <ProgressBar
              value={uploadProgress}
              variant="primary"
              showLabel
              label={`${uploadProgress}%`}
            />
          </>
        )}

        {status === 'success' && (
          <>
            <div className={styles.checkmark}>✓</div>
            <h2>Изображение загружено!</h2>
          </>
        )}

        {status === 'error' && (
          <>
            <div className={styles.errorIcon}>✕</div>
            <h2>Ошибка загрузки</h2>
            <p>Попробуйте еще раз</p>
            <Button variant="primary" onClick={handleRetry}>
              Повторить
            </Button>
          </>
        )}
      </div>
    </Container>
  );
}
```

---

### 3.4. QualityCheck Component

**Файл:** `src/pages/Help/QualityCheck.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Button } from '@/components/ui/Button';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';
import styles from './QualityCheck.module.css';

export default function QualityCheck() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const attemptId = location.state?.attemptId as string;
  const imageUrl = location.state?.imageUrl as string;

  useEffect(() => {
    if (!attemptId || !imageUrl) {
      navigate(ROUTES.HELP_UPLOAD);
      return;
    }

    // Analytics: help_quality_check_shown
    analytics.trackEvent('help_quality_check_shown', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });
  }, [attemptId, imageUrl]);

  const handleConfirm = async () => {
    // Analytics: help_quality_confirm_clicked
    analytics.trackEvent('help_quality_confirm_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    try {
      await helpAPI.confirmQuality(attemptId);

      // Navigate to processing
      navigate(ROUTES.HELP_PROCESSING, {
        state: { attemptId },
      });
    } catch (error) {
      console.error('[QualityCheck] Confirm failed:', error);
    }
  };

  const handleReshoot = () => {
    // Analytics: help_quality_reshoot_clicked
    analytics.trackEvent('help_quality_reshoot_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    navigate(ROUTES.HELP_UPLOAD);
  };

  const handleCrop = () => {
    // Analytics: help_quality_manual_crop_clicked
    analytics.trackEvent('help_quality_manual_crop_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    navigate(ROUTES.HELP_CROP, {
      state: { attemptId, imageUrl },
    });
  };

  return (
    <Container className={styles.container}>
      <div className={styles.header}>
        <button
          className={styles.backButton}
          onClick={() => navigate(ROUTES.HOME)}
        >
          Назад
        </button>
        <h1 className={styles.title}>Всё ли видно?</h1>
        <p className={styles.subtitle}>
          Проверь, что задание хорошо видно на фото
        </p>
      </div>

      <div className={styles.imagePreview}>
        <img src={imageUrl} alt="Task preview" className={styles.image} />
      </div>

      <div className={styles.actions}>
        <Button variant="outline" isFullWidth onClick={handleCrop}>
          Обрезать
        </Button>
        <Button variant="primary" size="lg" isFullWidth onClick={handleConfirm}>
          Всё видно, продолжить
        </Button>
        <Button variant="ghost" isFullWidth onClick={handleReshoot}>
          Переснять
        </Button>
      </div>
    </Container>
  );
}
```

---

### 3.5. Processing Component

**Файл:** `src/pages/Help/Processing.tsx`

```typescript
import { useEffect, useState, useRef } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Button } from '@/components/ui/Button';
import { Spinner } from '@/components/ui/Spinner';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';
import styles from './Processing.module.css';

const LONG_WAIT_THRESHOLD = 30000; // 30 seconds
const POLLING_INTERVAL = 2000; // 2 seconds

export default function Processing() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const attemptId = location.state?.attemptId as string;

  const [status, setStatus] = useState<'processing' | 'long_wait' | 'error'>('processing');
  const [startTime] = useState(Date.now());
  const pollingRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    if (!attemptId) {
      navigate(ROUTES.HELP_UPLOAD);
      return;
    }

    // Analytics: help_processing_started (sent by backend when processing starts)

    // Start processing
    startProcessing();

    // Check for long wait
    const longWaitTimer = setTimeout(() => {
      setStatus('long_wait');

      // Analytics: help_long_wait_shown
      const duration = Math.round((Date.now() - startTime) / 1000);
      analytics.trackEvent('help_long_wait_shown', {
        child_profile_id: profile?.id,
        attempt_id: attemptId,
        duration_seconds: duration,
      });
    }, LONG_WAIT_THRESHOLD);

    return () => {
      clearTimeout(longWaitTimer);
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
      }
    };
  }, [attemptId]);

  const startProcessing = async () => {
    try {
      // Trigger processing
      await helpAPI.processAttempt(attemptId);

      // Start polling for result
      pollingRef.current = setInterval(async () => {
        try {
          const result = await helpAPI.getResult(attemptId);

          // Processing complete
          if (pollingRef.current) {
            clearInterval(pollingRef.current);
          }

          // Navigate to result
          navigate(ROUTES.HELP_RESULT, {
            state: { attemptId, result },
          });
        } catch (error: any) {
          // Still processing - continue polling
          if (error.response?.status !== 404) {
            // Actual error
            console.error('[Processing] Polling error:', error);
            setStatus('error');
            if (pollingRef.current) {
              clearInterval(pollingRef.current);
            }
          }
        }
      }, POLLING_INTERVAL);
    } catch (error) {
      console.error('[Processing] Failed to start processing:', error);
      setStatus('error');
    }
  };

  const handleSaveAndWait = async () => {
    // Analytics: help_save_and_wait_clicked
    analytics.trackEvent('help_save_and_wait_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    try {
      await helpAPI.saveAndWait(attemptId);
      navigate(ROUTES.HOME);
    } catch (error) {
      console.error('[Processing] Save and wait failed:', error);
    }
  };

  const handleRetry = () => {
    // Analytics: help_retry_clicked
    analytics.trackEvent('help_retry_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    setStatus('processing');
    startProcessing();
  };

  const handleCancel = async () => {
    // Analytics: help_cancel_clicked
    analytics.trackEvent('help_cancel_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    try {
      await helpAPI.cancelAttempt(attemptId);
      navigate(ROUTES.HOME);
    } catch (error) {
      console.error('[Processing] Cancel failed:', error);
    }
  };

  return (
    <Container className={styles.container}>
      {status === 'processing' && (
        <div className={styles.content}>
          <Spinner size="lg" />
          <h2 className={styles.title}>Анализируем задание...</h2>
          <p className={styles.subtitle}>Это займет несколько секунд</p>
        </div>
      )}

      {status === 'long_wait' && (
        <div className={styles.content}>
          <Spinner size="lg" />
          <h2 className={styles.title}>Обработка занимает больше времени</h2>
          <p className={styles.subtitle}>
            Можешь продолжить ждать или сохранить и вернуться позже
          </p>
          <div className={styles.actions}>
            <Button variant="primary" onClick={handleSaveAndWait}>
              Сохранить и подождать
            </Button>
            <Button variant="ghost" onClick={handleCancel}>
              Отменить
            </Button>
          </div>
        </div>
      )}

      {status === 'error' && (
        <div className={styles.content}>
          <div className={styles.errorIcon}>✕</div>
          <h2 className={styles.title}>Ошибка обработки</h2>
          <p className={styles.subtitle}>Попробуйте еще раз</p>
          <div className={styles.actions}>
            <Button variant="primary" onClick={handleRetry}>
              Повторить
            </Button>
            <Button variant="ghost" onClick={() => navigate(ROUTES.HOME)}>
              На главную
            </Button>
          </div>
        </div>
      )}
    </Container>
  );
}
```

---

### 3.6. Result Component (с подсказками)

**Файл:** `src/pages/Help/Result.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';
import type { HelpResult, Hint } from '@/types/help';
import styles from './Result.module.css';

export default function Result() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const attemptId = location.state?.attemptId as string;
  const initialResult = location.state?.result as HelpResult;

  const [currentHintLevel, setCurrentHintLevel] = useState(1);
  const [hints, setHints] = useState<Hint[]>(initialResult?.hints || []);
  const [answer, setAnswer] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!attemptId || !initialResult) {
      navigate(ROUTES.HELP_UPLOAD);
      return;
    }

    // Analytics: help_result_opened
    analytics.trackEvent('help_result_opened', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
      used_hints_count: currentHintLevel,
    });
  }, [attemptId, initialResult]);

  const currentHint = hints.find((h) => h.level === currentHintLevel);

  useEffect(() => {
    if (currentHint) {
      // Analytics: hint_opened
      analytics.trackEvent('hint_opened', {
        child_profile_id: profile?.id,
        attempt_id: attemptId,
        hint_level: currentHintLevel,
      });
    }
  }, [currentHint]);

  const handleNextHint = async () => {
    if (currentHintLevel >= 3) return;

    // Analytics: next_hint_clicked
    analytics.trackEvent('next_hint_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
      from_hint_level: currentHintLevel,
    });

    try {
      const nextHint = await helpAPI.getNextHint(attemptId, currentHintLevel);
      setHints([...hints, nextHint]);
      setCurrentHintLevel(currentHintLevel + 1);
    } catch (error) {
      console.error('[Result] Failed to get next hint:', error);
    }
  };

  const handleSubmitAnswer = async () => {
    if (!answer.trim()) return;

    // Analytics: answer_submit_clicked
    analytics.trackEvent('answer_submit_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
      used_hints_count: currentHintLevel,
    });

    setIsSubmitting(true);

    try {
      const result = await helpAPI.submitAnswer(attemptId, answer);

      // Show success message
      // TODO: Show coins earned and damage dealt

      // Navigate to home
      setTimeout(() => {
        navigate(ROUTES.HOME);
      }, 2000);
    } catch (error) {
      console.error('[Result] Failed to submit answer:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleNewTask = () => {
    // Analytics: new_task_clicked_from_help
    analytics.trackEvent('new_task_clicked_from_help', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    navigate(ROUTES.HELP_UPLOAD);
  };

  return (
    <Container className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>Подсказки</h1>
        <p className={styles.subtitle}>Уровень {currentHintLevel}</p>
      </div>

      {currentHint && (
        <Card className={styles.hintCard} variant="elevated">
          <h3 className={styles.hintTitle}>{currentHint.title}</h3>
          <div
            className={styles.hintContent}
            dangerouslySetInnerHTML={{ __html: currentHint.content }}
          />
        </Card>
      )}

      {currentHintLevel < 3 && (
        <Button variant="outline" isFullWidth onClick={handleNextHint}>
          Следующая подсказка (уровень {currentHintLevel + 1})
        </Button>
      )}

      <div className={styles.answerSection}>
        <h3 className={styles.answerTitle}>Твой ответ</h3>
        <Input
          placeholder="Введи ответ..."
          value={answer}
          onChange={(e) => setAnswer(e.target.value)}
          isFullWidth
        />
        <Button
          variant="primary"
          size="lg"
          isFullWidth
          onClick={handleSubmitAnswer}
          isLoading={isSubmitting}
          disabled={!answer.trim()}
        >
          Отправить ответ
        </Button>
      </div>

      <Button variant="ghost" isFullWidth onClick={handleNewTask}>
        Новое задание
      </Button>
    </Container>
  );
}
```

---

## Часть 4: Аналитические события

### События Help Flow

| Event Name | Когда | Параметры |
|------------|-------|-----------|
| `help_flow_started` | Начат поток помощи | `child_profile_id`, `mode`, `entry_point` |
| `help_source_picker_opened` | Открыт выбор источника | `child_profile_id` |
| `help_choose_file_clicked` | Выбран файл | `child_profile_id` |
| `help_camera_clicked` | Выбрана камера | `child_profile_id` |
| `help_clipboard_clicked` | Выбран буфер | `child_profile_id` |
| `help_dragdrop_used` | Использован drag&drop | `child_profile_id` |
| `help_image_selected` | Изображение выбрано | `child_profile_id`, `upload_source`, `file_size_bucket`, `mime_type` |
| `help_image_upload_started` | Старт загрузки | `child_profile_id`, `upload_source`, `attempt_id` |
| `help_image_upload_completed` | Загрузка завершена (backend) | `child_profile_id`, `attempt_id`, `attempt_image_id`, `file_size_bucket` |
| `help_image_upload_failed` | Ошибка загрузки (backend) | `child_profile_id`, `error_code`, `attempt_id` |
| `help_quality_check_shown` | Показан экран проверки качества | `child_profile_id`, `attempt_id` |
| `help_quality_confirm_clicked` | Подтверждено качество | `child_profile_id`, `attempt_id` |
| `help_quality_reshoot_clicked` | Выбран пересъём | `child_profile_id`, `attempt_id` |
| `help_quality_manual_crop_clicked` | Выбран crop | `child_profile_id`, `attempt_id` |
| `help_processing_started` | Старт обработки (backend) | `child_profile_id`, `attempt_id` |
| `help_long_wait_shown` | Показан long-wait | `child_profile_id`, `attempt_id`, `duration_seconds` |
| `help_save_and_wait_clicked` | Нажато "Сохранить и подождать" | `child_profile_id`, `attempt_id` |
| `help_retry_clicked` | Нажато "Повторить" | `child_profile_id`, `attempt_id` |
| `help_cancel_clicked` | Нажато "Отменить" | `child_profile_id`, `attempt_id` |
| `help_result_opened` | Открыт результат | `child_profile_id`, `attempt_id`, `used_hints_count` |
| `hint_opened` | Открыта подсказка | `child_profile_id`, `attempt_id`, `hint_level` |
| `next_hint_clicked` | Следующая подсказка | `child_profile_id`, `attempt_id`, `from_hint_level` |
| `answer_submit_clicked` | Отправка ответа | `child_profile_id`, `attempt_id`, `used_hints_count` |
| `new_task_clicked_from_help` | Новое задание | `child_profile_id`, `attempt_id` |

---

## Чеклист задач

### Компоненты
- [ ] Создать SourcePicker с выбором источника
- [ ] Создать ImageUpload с прогрессом
- [ ] Создать QualityCheck с предпросмотром
- [ ] Создать Crop компонент (опционально)
- [ ] Создать Processing с polling
- [ ] Создать Result с подсказками
- [ ] Создать Answer input форму

### API
- [ ] Реализовать helpAPI.createAttempt()
- [ ] Реализовать helpAPI.uploadImage()
- [ ] Реализовать helpAPI.processAttempt()
- [ ] Реализовать helpAPI.getResult() с polling
- [ ] Реализовать helpAPI.getNextHint()
- [ ] Реализовать helpAPI.submitAnswer()

### Features
- [ ] Добавить сжатие изображений
- [ ] Добавить drag&drop
- [ ] Добавить вставку из буфера
- [ ] Добавить long-wait экран
- [ ] Добавить многоуровневые подсказки

### Аналитика
- [ ] Добавить все события из таблицы
- [ ] Отправлять параметры файла (размер, тип)
- [ ] Трекать использование подсказок

### Testing
- [ ] Протестировать загрузку изображений
- [ ] Протестировать обработку и polling
- [ ] Протестировать подсказки и ответы
- [ ] Протестировать обработку ошибок

---

## Следующий этап

После завершения Help Flow переходи к **06_CHECK.md** для создания потока проверки ДЗ.
