# Phase 5: Поток "Проверка ДЗ" (Check Flow)

**Длительность:** 5-7 дней
**Приоритет:** Критический
**Зависимости:** 02_CORE.md, 04_HOME.md, 05_HELP.md

---

## Цель

Реализовать полный поток проверки домашнего задания: выбор сценария (1 или 2 фото), загрузка изображений, обработка, отображение результата с ошибками, исправление и повторная отправка.

---

## Архитектура потока

### Этапы Check Flow

```
1. Scenario Selection  → Выбор сценария (1 фото / 2 фото)
2. Image Upload        → Загрузка изображений (task_image, answer_image)
3. Quality Check       → Проверка качества каждого фото
4. Crop (optional)     → Обрезка изображений
5. Processing          → Обработка (loading, long-wait, save-and-wait)
6. Result              → Показ результата (успех / ошибки)
7. Error Display       → Отображение ошибок с локализацией
8. Fix & Resubmit      → Исправление и повторная проверка
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/check.ts`

```typescript
export type CheckScenario = 'single_photo' | 'two_photo';

export type ImageRole = 'task' | 'answer';

export interface CheckAttempt {
  id: string;
  childProfileId: string;
  mode: 'check';
  scenario: CheckScenario;
  status: CheckStatus;
  taskImageUrl?: string;
  answerImageUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export type CheckStatus =
  | 'created'
  | 'uploading_task'
  | 'uploading_answer'
  | 'uploaded'
  | 'processing'
  | 'long_wait'
  | 'completed'
  | 'failed';

export interface CheckResult {
  attemptId: string;
  status: 'success' | 'error';
  errors?: CheckError[];
  coinsEarned: number;
  damageDealt: number;
}

export interface CheckError {
  id: string;
  stepNumber: number | null;
  lineReference: string | null;
  description: string;
  locationType: 'step' | 'line' | 'general';
  severity: 'error' | 'warning';
}

export interface ErrorLocalization {
  type: 'step' | 'line' | 'general';
  reference: string | null;
  description: string;
}
```

---

## Часть 2: API Integration

### 2.1. Check API

**Файл:** `src/api/check.ts`

```typescript
import { apiClient } from './client';
import type {
  CheckAttempt,
  CheckResult,
  CheckScenario,
  ImageRole,
} from '@/types/check';

export const checkAPI = {
  /**
   * Создать новую попытку проверки
   */
  async createAttempt(
    childProfileId: string,
    scenario: CheckScenario
  ): Promise<CheckAttempt> {
    return apiClient.post<CheckAttempt>('/attempts', {
      childProfileId,
      mode: 'check',
      scenarioType: scenario,
    });
  },

  /**
   * Загрузить изображение (задание или ответ)
   */
  async uploadImage(
    attemptId: string,
    imageRole: ImageRole,
    file: Blob,
    onProgress?: (progress: number) => void
  ): Promise<{ imageUrl: string; thumbnailUrl: string }> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('imageRole', imageRole);

    return apiClient.upload<{ imageUrl: string; thumbnailUrl: string }>(
      `/attempts/${attemptId}/images`,
      file,
      onProgress
    );
  },

  /**
   * Подтвердить качество изображения
   */
  async confirmQuality(
    attemptId: string,
    imageRole: ImageRole
  ): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/confirm-quality`, {
      imageRole,
    });
  },

  /**
   * Обрезать изображение
   */
  async cropImage(
    attemptId: string,
    imageRole: ImageRole,
    cropArea: { x: number; y: number; width: number; height: number }
  ): Promise<{ imageUrl: string }> {
    return apiClient.post<{ imageUrl: string }>(
      `/attempts/${attemptId}/crop`,
      { imageRole, cropArea }
    );
  },

  /**
   * Начать обработку проверки
   */
  async processAttempt(attemptId: string): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/process`);
  },

  /**
   * Получить результат проверки (polling)
   */
  async getResult(attemptId: string): Promise<CheckResult> {
    return apiClient.get<CheckResult>(`/attempts/${attemptId}/result`);
  },

  /**
   * Сохранить и подождать
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

### 3.1. ScenarioSelector Component

**Файл:** `src/pages/Check/ScenarioSelector.tsx`

```typescript
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Card } from '@/components/ui/Card';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';
import type { CheckScenario } from '@/types/check';
import styles from './ScenarioSelector.module.css';

export default function ScenarioSelector() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  useEffect(() => {
    // Analytics: check_flow_started
    analytics.trackEvent('check_flow_started', {
      child_profile_id: profile?.id,
      mode: 'check',
      entry_point: 'home_button',
    });

    // Analytics: check_scenario_picker_opened
    analytics.trackEvent('check_scenario_picker_opened', {
      child_profile_id: profile?.id,
    });
  }, []);

  const handleScenarioSelect = (scenario: CheckScenario) => {
    // Analytics: check_scenario_selected
    analytics.trackEvent('check_scenario_selected', {
      child_profile_id: profile?.id,
      check_scenario: scenario,
    });

    if (scenario === 'single_photo') {
      // Analytics: check_single_photo_selected
      analytics.trackEvent('check_single_photo_selected', {
        child_profile_id: profile?.id,
        check_scenario: scenario,
      });
    } else {
      // Analytics: check_two_photo_selected
      analytics.trackEvent('check_two_photo_selected', {
        child_profile_id: profile?.id,
        check_scenario: scenario,
      });
    }

    // Navigate to upload with scenario
    navigate(ROUTES.CHECK_UPLOAD, {
      state: { scenario },
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
        <h1 className={styles.title}>Проверка ДЗ</h1>
        <p className={styles.subtitle}>Выбери, как проверить работу</p>
      </div>

      <div className={styles.scenarios}>
        <Card
          className={styles.scenarioCard}
          variant="bordered"
          onClick={() => handleScenarioSelect('single_photo')}
        >
          <div className={styles.scenarioIcon}>📸</div>
          <h3 className={styles.scenarioTitle}>Одно фото</h3>
          <p className={styles.scenarioDescription}>
            Задание и ответ на одном фото
          </p>
        </Card>

        <Card
          className={styles.scenarioCard}
          variant="bordered"
          onClick={() => handleScenarioSelect('two_photo')}
        >
          <div className={styles.scenarioIcon}>📷📷</div>
          <h3 className={styles.scenarioTitle}>Два фото</h3>
          <p className={styles.scenarioDescription}>
            Отдельно фото задания и фото ответа
          </p>
        </Card>
      </div>
    </Container>
  );
}
```

---

### 3.2. CheckImageUpload Component

**Файл:** `src/pages/Check/CheckImageUpload.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Button } from '@/components/ui/Button';
import { ProgressBar } from '@/components/ui/ProgressBar';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { checkAPI } from '@/api/check';
import { ROUTES } from '@/config/routes';
import imageCompression from 'browser-image-compression';
import type { CheckScenario, ImageRole } from '@/types/check';
import styles from './CheckImageUpload.module.css';

export default function CheckImageUpload() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const scenario = location.state?.scenario as CheckScenario;
  const attemptId = location.state?.attemptId as string | undefined;

  const [currentRole, setCurrentRole] = useState<ImageRole>('task');
  const [taskImageUrl, setTaskImageUrl] = useState<string | null>(null);
  const [answerImageUrl, setAnswerImageUrl] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);

  useEffect(() => {
    if (!scenario) {
      navigate(ROUTES.CHECK_SCENARIO);
    }
  }, [scenario]);

  const openFilePicker = (role: ImageRole) => {
    setCurrentRole(role);

    // Analytics: check_source_picker_opened
    analytics.trackEvent('check_source_picker_opened', {
      child_profile_id: profile?.id,
      check_scenario: scenario,
    });

    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/jpeg,image/png';
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        // Analytics: check_choose_file_clicked
        analytics.trackEvent('check_choose_file_clicked', {
          child_profile_id: profile?.id,
          check_scenario: scenario,
        });

        handleFileSelected(file, role, 'file');
      }
    };
    input.click();
  };

  const handleFileSelected = async (
    file: File,
    role: ImageRole,
    source: string
  ) => {
    // Analytics: check_task_image_selected / check_answer_image_selected
    const eventName =
      role === 'task'
        ? 'check_task_image_selected'
        : 'check_answer_image_selected';

    analytics.trackEvent(eventName, {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
      image_role: role,
      upload_source: source,
    });

    setIsUploading(true);
    setUploadProgress(0);

    try {
      // Create attempt if not exists
      let currentAttemptId = attemptId;
      if (!currentAttemptId) {
        const attempt = await checkAPI.createAttempt(profile!.id, scenario);
        currentAttemptId = attempt.id;
      }

      // Analytics: check_image_upload_started
      analytics.trackEvent('check_image_upload_started', {
        child_profile_id: profile?.id,
        attempt_id: currentAttemptId,
        image_role: role,
        upload_source: source,
      });

      // Compress image
      const compressedFile = await imageCompression(file, {
        maxSizeMB: 10,
        maxWidthOrHeight: 1920,
        useWebWorker: true,
        initialQuality: 0.8,
      });

      // Upload image
      const result = await checkAPI.uploadImage(
        currentAttemptId,
        role,
        compressedFile,
        (progress) => setUploadProgress(progress)
      );

      // Analytics: check_image_upload_completed (sent by backend)

      // Save image URL
      if (role === 'task') {
        setTaskImageUrl(result.imageUrl);
      } else {
        setAnswerImageUrl(result.imageUrl);
      }

      setIsUploading(false);

      // Navigate to quality check for this image
      navigate(ROUTES.CHECK_QUALITY, {
        state: {
          scenario,
          attemptId: currentAttemptId,
          imageRole: role,
          imageUrl: result.imageUrl,
          taskImageUrl,
          answerImageUrl: role === 'answer' ? result.imageUrl : answerImageUrl,
        },
      });
    } catch (error) {
      console.error('[CheckImageUpload] Upload failed:', error);
      setIsUploading(false);

      // Analytics: check_image_upload_failed (sent by backend)
    }
  };

  const needsTaskImage = !taskImageUrl;
  const needsAnswerImage = scenario === 'two_photo' && !answerImageUrl;

  return (
    <Container className={styles.container}>
      <div className={styles.header}>
        <button
          className={styles.backButton}
          onClick={() => navigate(ROUTES.CHECK_SCENARIO)}
        >
          Назад
        </button>
        <h1 className={styles.title}>Загрузи фото</h1>
        <p className={styles.subtitle}>
          {scenario === 'single_photo'
            ? 'Загрузи фото с заданием и ответом'
            : 'Загрузи фото задания и ответа'}
        </p>
      </div>

      <div className={styles.uploadArea}>
        {isUploading ? (
          <div className={styles.uploadProgress}>
            <h3>Загружаем изображение...</h3>
            <ProgressBar
              value={uploadProgress}
              variant="primary"
              showLabel
              label={`${uploadProgress}%`}
            />
          </div>
        ) : (
          <>
            {needsTaskImage && (
              <UploadButton
                title={
                  scenario === 'single_photo'
                    ? 'Загрузить фото'
                    : 'Загрузить фото задания'
                }
                icon="📷"
                onClick={() => openFilePicker('task')}
              />
            )}

            {taskImageUrl && needsAnswerImage && (
              <>
                <div className={styles.uploadedIndicator}>
                  <span className={styles.checkmark}>✓</span>
                  <span>Фото задания загружено</span>
                </div>

                <UploadButton
                  title="Загрузить фото ответа"
                  icon="📷"
                  onClick={() => openFilePicker('answer')}
                />
              </>
            )}
          </>
        )}
      </div>
    </Container>
  );
}

interface UploadButtonProps {
  title: string;
  icon: string;
  onClick: () => void;
}

function UploadButton({ title, icon, onClick }: UploadButtonProps) {
  return (
    <Button
      variant="primary"
      size="lg"
      isFullWidth
      onClick={onClick}
      leftIcon={icon}
      className={styles.uploadButton}
    >
      {title}
    </Button>
  );
}
```

---

### 3.3. CheckResult Component

**Файл:** `src/pages/Check/CheckResult.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';
import type { CheckResult as CheckResultType, CheckError } from '@/types/check';
import styles from './CheckResult.module.css';

export default function CheckResult() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const attemptId = location.state?.attemptId as string;
  const result = location.state?.result as CheckResultType;

  const [selectedError, setSelectedError] = useState<CheckError | null>(null);

  useEffect(() => {
    if (!attemptId || !result) {
      navigate(ROUTES.CHECK_SCENARIO);
      return;
    }

    // Analytics: check_result_opened
    analytics.trackEvent('check_result_opened', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
      result_status: result.status,
      error_count: result.errors?.length || 0,
    });

    // Analytics: check_error_feedback_viewed (if has errors)
    if (result.errors && result.errors.length > 0) {
      analytics.trackEvent('check_error_feedback_viewed', {
        child_profile_id: profile?.id,
        attempt_id: attemptId,
        error_count: result.errors.length,
      });

      // Analytics: soft_error_message_shown
      analytics.trackEvent('soft_error_message_shown', {
        child_profile_id: profile?.id,
        attempt_id: attemptId,
        error_count: result.errors.length,
      });
    }
  }, [attemptId, result]);

  const handleErrorClick = (error: CheckError) => {
    setSelectedError(error);

    // Analytics: error_hint_block_opened
    analytics.trackEvent('error_hint_block_opened', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
      error_block_id: error.id,
      step_number: error.stepNumber,
      line_reference: error.lineReference,
    });

    // Analytics: error_location_viewed
    analytics.trackEvent('error_location_viewed', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
      location_type: error.locationType,
    });
  };

  const handleRetry = () => {
    // Analytics: retry_after_errors_clicked
    analytics.trackEvent('retry_after_errors_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    // Analytics: user_retries_after_error
    analytics.trackEvent('user_retries_after_error', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    navigate(ROUTES.CHECK_SCENARIO);
  };

  const handleFixAndResubmit = () => {
    // Analytics: fixed_and_resubmit_clicked
    analytics.trackEvent('fixed_and_resubmit_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    // Go back to upload with same scenario
    navigate(ROUTES.CHECK_UPLOAD, {
      state: { scenario: location.state?.scenario },
    });
  };

  const handleNewTask = () => {
    // Analytics: new_task_clicked_from_check
    analytics.trackEvent('new_task_clicked_from_check', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    navigate(ROUTES.CHECK_SCENARIO);
  };

  const handleAbandon = () => {
    // Analytics: user_abandons_after_error
    analytics.trackEvent('user_abandons_after_error', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    navigate(ROUTES.HOME);
  };

  const hasErrors = result.errors && result.errors.length > 0;

  return (
    <Container className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>
          {hasErrors ? 'Нашли ошибки' : 'Отлично!'}
        </h1>
        <p className={styles.subtitle}>
          {hasErrors
            ? `Найдено ошибок: ${result.errors!.length}`
            : 'Все правильно!'}
        </p>
      </div>

      {/* Success state */}
      {!hasErrors && (
        <div className={styles.successContent}>
          <div className={styles.successIcon}>✅</div>
          <h2 className={styles.successTitle}>Задание выполнено правильно!</h2>
          <div className={styles.rewards}>
            <div className={styles.reward}>
              <span className={styles.rewardIcon}>💰</span>
              <span className={styles.rewardValue}>
                +{result.coinsEarned} монет
              </span>
            </div>
            {result.damageDealt > 0 && (
              <div className={styles.reward}>
                <span className={styles.rewardIcon}>⚔️</span>
                <span className={styles.rewardValue}>
                  -{result.damageDealt} HP злодею
                </span>
              </div>
            )}
          </div>

          <Button
            variant="primary"
            size="lg"
            isFullWidth
            onClick={() => navigate(ROUTES.HOME)}
          >
            На главную
          </Button>
          <Button variant="ghost" isFullWidth onClick={handleNewTask}>
            Проверить ещё
          </Button>
        </div>
      )}

      {/* Error state */}
      {hasErrors && (
        <div className={styles.errorContent}>
          <div className={styles.errorList}>
            {result.errors!.map((error, index) => (
              <ErrorCard
                key={error.id}
                error={error}
                index={index}
                onClick={() => handleErrorClick(error)}
                isSelected={selectedError?.id === error.id}
              />
            ))}
          </div>

          <div className={styles.actions}>
            <Button variant="primary" isFullWidth onClick={handleFixAndResubmit}>
              Исправил, проверить снова
            </Button>
            <Button variant="outline" isFullWidth onClick={handleRetry}>
              Попробовать другое задание
            </Button>
            <Button variant="ghost" isFullWidth onClick={handleAbandon}>
              На главную
            </Button>
          </div>
        </div>
      )}

      {/* Error detail modal */}
      {selectedError && (
        <ErrorDetailModal
          error={selectedError}
          onClose={() => setSelectedError(null)}
        />
      )}
    </Container>
  );
}

// ErrorCard component
interface ErrorCardProps {
  error: CheckError;
  index: number;
  onClick: () => void;
  isSelected: boolean;
}

function ErrorCard({ error, index, onClick, isSelected }: ErrorCardProps) {
  return (
    <Card
      className={clsx(styles.errorCard, {
        [styles.selected]: isSelected,
      })}
      variant="bordered"
      onClick={onClick}
    >
      <div className={styles.errorHeader}>
        <span className={styles.errorNumber}>Ошибка {index + 1}</span>
        {error.stepNumber && (
          <span className={styles.errorLocation}>Шаг {error.stepNumber}</span>
        )}
        {error.lineReference && (
          <span className={styles.errorLocation}>
            Строка {error.lineReference}
          </span>
        )}
      </div>
      <p className={styles.errorDescription}>{error.description}</p>
    </Card>
  );
}

// ErrorDetailModal component
interface ErrorDetailModalProps {
  error: CheckError;
  onClose: () => void;
}

function ErrorDetailModal({ error, onClose }: ErrorDetailModalProps) {
  return (
    <Modal isOpen onClose={onClose} title="Детали ошибки">
      <div className={styles.errorDetail}>
        {error.stepNumber && (
          <div className={styles.detailSection}>
            <h4>Шаг</h4>
            <p>{error.stepNumber}</p>
          </div>
        )}
        {error.lineReference && (
          <div className={styles.detailSection}>
            <h4>Строка</h4>
            <p>{error.lineReference}</p>
          </div>
        )}
        <div className={styles.detailSection}>
          <h4>Описание</h4>
          <p>{error.description}</p>
        </div>
        <div className={styles.detailSection}>
          <h4>Тип</h4>
          <p>{getLocationTypeLabel(error.locationType)}</p>
        </div>
      </div>
    </Modal>
  );
}

function getLocationTypeLabel(type: string): string {
  switch (type) {
    case 'step':
      return 'Ошибка в шаге решения';
    case 'line':
      return 'Ошибка в строке';
    case 'general':
      return 'Общая ошибка';
    default:
      return type;
  }
}
```

---

## Часть 4: Аналитические события

### События Check Flow

| Event Name | Когда | Параметры |
|------------|-------|-----------|
| `check_flow_started` | Начат поток проверки | `child_profile_id`, `mode`, `entry_point` |
| `check_scenario_picker_opened` | Открыт выбор сценария | `child_profile_id` |
| `check_scenario_selected` | Выбран сценарий | `child_profile_id`, `check_scenario` |
| `check_single_photo_selected` | Выбран сценарий 1 фото | `child_profile_id`, `check_scenario` |
| `check_two_photo_selected` | Выбран сценарий 2 фото | `child_profile_id`, `check_scenario` |
| `check_task_image_selected` | Выбрано фото задания | `child_profile_id`, `attempt_id`, `image_role`, `upload_source` |
| `check_answer_image_selected` | Выбрано фото ответа | `child_profile_id`, `attempt_id`, `image_role`, `upload_source` |
| `check_source_picker_opened` | Открыт выбор источника | `child_profile_id`, `check_scenario` |
| `check_choose_file_clicked` | Выбран файл | `child_profile_id`, `check_scenario` |
| `check_camera_clicked` | Выбрана камера | `child_profile_id`, `check_scenario` |
| `check_clipboard_clicked` | Выбран буфер | `child_profile_id`, `check_scenario` |
| `check_dragdrop_used` | Использован drag&drop | `child_profile_id`, `check_scenario` |
| `check_image_upload_started` | Старт загрузки | `child_profile_id`, `attempt_id`, `image_role`, `upload_source` |
| `check_image_upload_completed` | Загрузка завершена (backend) | `child_profile_id`, `attempt_id`, `attempt_image_id`, `image_role`, `file_size_bucket` |
| `check_image_upload_failed` | Ошибка загрузки (backend) | `child_profile_id`, `attempt_id`, `image_role`, `error_code` |
| `upload_more_clicked` | Нажата "Загрузить ещё" | `child_profile_id`, `mode`, `image_role_expected`, `attempt_id` |
| `check_quality_check_shown` | Показан экран проверки качества | `child_profile_id`, `attempt_id`, `attempt_image_id`, `image_role` |
| `check_quality_confirm_clicked` | Подтверждено качество | `child_profile_id`, `attempt_id`, `attempt_image_id`, `image_role` |
| `check_quality_reshoot_clicked` | Выбран пересъём | `child_profile_id`, `attempt_id`, `attempt_image_id`, `image_role` |
| `check_quality_manual_crop_clicked` | Выбран crop | `child_profile_id`, `attempt_id`, `attempt_image_id`, `image_role` |
| `check_processing_started` | Старт обработки (backend) | `child_profile_id`, `attempt_id`, `check_scenario` |
| `check_long_wait_shown` | Показан long-wait | `child_profile_id`, `attempt_id`, `duration_seconds` |
| `check_save_and_wait_clicked` | Нажато "Сохранить и подождать" | `child_profile_id`, `attempt_id` |
| `check_retry_clicked` | Нажато "Повторить" | `child_profile_id`, `attempt_id` |
| `check_cancel_clicked` | Нажато "Отменить" | `child_profile_id`, `attempt_id` |
| `check_result_opened` | Открыт результат | `child_profile_id`, `attempt_id`, `result_status`, `error_count` |
| `check_error_feedback_viewed` | Показан блок ошибок | `child_profile_id`, `attempt_id`, `error_count` |
| `soft_error_message_shown` | Показана мягкая формулировка | `child_profile_id`, `attempt_id`, `error_count` |
| `error_hint_block_opened` | Открыт блок ошибки | `child_profile_id`, `attempt_id`, `error_block_id`, `step_number`, `line_reference` |
| `error_location_viewed` | Просмотрена локализация | `child_profile_id`, `attempt_id`, `location_type` |
| `retry_after_errors_clicked` | Нажато повторить после ошибки | `child_profile_id`, `attempt_id` |
| `fixed_and_resubmit_clicked` | Нажато "Исправил" | `child_profile_id`, `attempt_id` |
| `new_task_clicked_from_check` | Новое задание | `child_profile_id`, `attempt_id` |
| `user_retries_after_error` | Пользователь пытается ещё раз | `child_profile_id`, `attempt_id` |
| `user_abandons_after_error` | Пользователь бросает | `child_profile_id`, `attempt_id` |

---

## Часть 5: API Endpoints

### Эндпоинты для Check Flow

```typescript
// POST /api/v1/attempts
// Request:
{
  childProfileId: string;
  mode: 'check';
  scenarioType: 'single_photo' | 'two_photo';
}
// Response:
{
  id: string;
  childProfileId: string;
  mode: 'check';
  scenario: 'single_photo' | 'two_photo';
  status: 'created';
  createdAt: string;
  updatedAt: string;
}

// POST /api/v1/attempts/:attemptId/images
// Form Data:
{
  file: Blob;
  imageRole: 'task' | 'answer';
}
// Response:
{
  imageUrl: string;
  thumbnailUrl: string;
}

// POST /api/v1/attempts/:attemptId/process
// Response: 204 No Content

// GET /api/v1/attempts/:attemptId/result
// Response:
{
  attemptId: string;
  status: 'success' | 'error';
  errors?: Array<{
    id: string;
    stepNumber: number | null;
    lineReference: string | null;
    description: string;
    locationType: 'step' | 'line' | 'general';
    severity: 'error' | 'warning';
  }>;
  coinsEarned: number;
  damageDealt: number;
}
```

---

## Часть 6: Error Localization

### Типы локализации ошибок

```typescript
// 1. По шагу решения (step)
{
  locationType: 'step',
  stepNumber: 3,
  description: 'Ошибка в третьем шаге: неправильно применена формула'
}

// 2. По строке (line)
{
  locationType: 'line',
  lineReference: '5',
  description: 'Ошибка в строке 5: арифметическая ошибка'
}

// 3. Общая ошибка (general)
{
  locationType: 'general',
  description: 'Неправильный подход к решению задачи'
}
```

---

## Чеклист задач

### Компоненты
- [ ] Создать ScenarioSelector с выбором сценария
- [ ] Создать CheckImageUpload для двух изображений
- [ ] Создать QualityCheck для каждого изображения
- [ ] Создать Processing с polling
- [ ] Создать CheckResult с отображением ошибок
- [ ] Создать ErrorCard компоненты
- [ ] Создать ErrorDetailModal

### API
- [ ] Реализовать checkAPI.createAttempt()
- [ ] Реализовать checkAPI.uploadImage() для task/answer
- [ ] Реализовать checkAPI.processAttempt()
- [ ] Реализовать checkAPI.getResult()
- [ ] Обработка разных типов ошибок

### Features
- [ ] Поддержка двух сценариев (1 фото / 2 фото)
- [ ] Отображение ошибок с локализацией
- [ ] Мягкие формулировки об ошибках
- [ ] Возможность исправления и повторной отправки

### Аналитика
- [ ] Добавить все события из таблицы
- [ ] Трекать выбор сценария
- [ ] Трекать загрузку обоих изображений
- [ ] Трекать просмотр и взаимодействие с ошибками

### Testing
- [ ] Протестировать оба сценария
- [ ] Протестировать отображение ошибок
- [ ] Протестировать исправление и повтор
- [ ] Протестировать обработку ошибок

---

## Следующий этап

После завершения Check Flow основные критические потоки готовы. Далее можно переходить к второстепенным экранам: Достижения, Друзья, Профиль.
