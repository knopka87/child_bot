# Phase 2: Онбординг и регистрация

**Длительность:** 5-6 дней
**Приоритет:** Критический
**Зависимости:** 01_SETUP.md, 02_CORE.md

---

## Цель

Создать полный онбординг пользователя: экраны приветствия, выбор класса, выбор аватара, ввод имени, согласия, email верификация, и сохранение профиля с полной аналитикой.

---

## Архитектура онбординга

### Flow диаграмма

```
Start
  ↓
Onboarding Screen (welcome)
  ↓
Grade Selection (1-11 класс)
  ↓
Avatar Selection (grid 4x3)
  ↓
Display Name Input
  ↓
Consent Screen (adult_consent checkbox)
  ↓
Email Input
  ↓
Email Verification Waiting
  ↓
Onboarding Complete → Home
```

### Структура компонентов

```
OnboardingFlow
├── WelcomeScreen
├── GradeSelection
├── AvatarSelection
├── DisplayNameInput
├── ConsentScreen
│   ├── AdultConsentCheckbox
│   ├── PrivacyPolicyLink
│   └── TermsOfServiceLink
├── EmailInput
├── EmailVerificationWaiting
└── OnboardingCompletedScreen
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/onboarding.ts`

```typescript
export type OnboardingStep =
  | 'welcome'
  | 'grade'
  | 'avatar'
  | 'display_name'
  | 'consent'
  | 'email'
  | 'email_verification'
  | 'completed';

export interface OnboardingState {
  currentStep: OnboardingStep;
  grade: number | null;
  avatarId: string | null;
  displayName: string;
  adultConsent: boolean;
  privacyAccepted: boolean;
  termsAccepted: boolean;
  email: string;
  emailVerified: boolean;
}

export interface Avatar {
  id: string;
  imageUrl: string;
  name: string;
  isPremium: boolean;
}

export interface ConsentDocument {
  type: 'privacy_policy' | 'terms_of_service';
  version: string;
  url: string;
  acceptedAt?: string;
}

export interface OnboardingProgress {
  step: OnboardingStep;
  completedSteps: OnboardingStep[];
  totalSteps: number;
  progressPercent: number;
}
```

---

## Часть 2: API Integration

### 2.1. Onboarding API

**Файл:** `src/api/onboarding.ts`

```typescript
import { apiClient } from './client';
import type { Avatar, ConsentDocument } from '@/types/onboarding';

export const onboardingAPI = {
  /**
   * Получить список аватаров
   */
  async getAvatars(): Promise<Avatar[]> {
    return apiClient.get<Avatar[]>('/avatars');
  },

  /**
   * Создать профиль ребёнка
   */
  async createChildProfile(data: {
    parentUserId: string;
    grade: number;
    avatarId: string;
    displayName: string;
  }): Promise<{ childProfileId: string }> {
    return apiClient.post<{ childProfileId: string }>('/profiles/child', data);
  },

  /**
   * Отправить email для верификации
   */
  async sendEmailVerification(email: string): Promise<void> {
    return apiClient.post<void>('/auth/email/send-verification', { email });
  },

  /**
   * Проверить статус верификации email
   */
  async checkEmailVerification(email: string): Promise<{ verified: boolean }> {
    return apiClient.get<{ verified: boolean }>('/auth/email/check-verification', {
      params: { email },
    });
  },

  /**
   * Сохранить согласие на обработку данных
   */
  async saveConsent(data: {
    parentUserId: string;
    privacyPolicyVersion: string;
    termsVersion: string;
    adultConsent: boolean;
  }): Promise<void> {
    return apiClient.post<void>('/consent', data);
  },

  /**
   * Завершить онбординг
   */
  async completeOnboarding(data: {
    parentUserId: string;
    childProfileId: string;
  }): Promise<void> {
    return apiClient.post<void>('/onboarding/complete', data);
  },
};
```

---

## Часть 3: Компоненты

### 3.1. OnboardingFlow Container

**Файл:** `src/pages/Onboarding/OnboardingFlow.tsx`

```typescript
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { WelcomeScreen } from './screens/WelcomeScreen';
import { GradeSelection } from './screens/GradeSelection';
import { AvatarSelection } from './screens/AvatarSelection';
import { DisplayNameInput } from './screens/DisplayNameInput';
import { ConsentScreen } from './screens/ConsentScreen';
import { EmailInput } from './screens/EmailInput';
import { EmailVerificationWaiting } from './screens/EmailVerificationWaiting';
import { OnboardingCompletedScreen } from './screens/OnboardingCompletedScreen';
import { ProgressBar } from '@/components/ui/ProgressBar';
import { useAnalytics } from '@/hooks/useAnalytics';
import { usePlatform } from '@/hooks/usePlatform';
import { useProfileStore } from '@/stores/profileStore';
import type { OnboardingStep, OnboardingState } from '@/types/onboarding';
import { ROUTES } from '@/config/routes';
import styles from './OnboardingFlow.module.css';

export default function OnboardingFlow() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const platform = usePlatform();
  const setProfile = useProfileStore((state) => state.setProfile);

  const [state, setState] = useState<OnboardingState>({
    currentStep: 'welcome',
    grade: null,
    avatarId: null,
    displayName: '',
    adultConsent: false,
    privacyAccepted: false,
    termsAccepted: false,
    email: '',
    emailVerified: false,
  });

  useEffect(() => {
    // Analytics: onboarding_opened
    analytics.trackEvent('onboarding_opened', {
      platform_type: platform.type,
      session_id: analytics.sessionId,
      entry_point: 'app_start',
    });
  }, []);

  const handleStepComplete = (updates: Partial<OnboardingState>) => {
    setState((prev) => ({ ...prev, ...updates }));
  };

  const goToNextStep = () => {
    const steps: OnboardingStep[] = [
      'welcome',
      'grade',
      'avatar',
      'display_name',
      'consent',
      'email',
      'email_verification',
      'completed',
    ];

    const currentIndex = steps.indexOf(state.currentStep);
    if (currentIndex < steps.length - 1) {
      setState((prev) => ({
        ...prev,
        currentStep: steps[currentIndex + 1],
      }));
    }
  };

  const goToPreviousStep = () => {
    const steps: OnboardingStep[] = [
      'welcome',
      'grade',
      'avatar',
      'display_name',
      'consent',
      'email',
      'email_verification',
    ];

    const currentIndex = steps.indexOf(state.currentStep);
    if (currentIndex > 0) {
      setState((prev) => ({
        ...prev,
        currentStep: steps[currentIndex - 1],
      }));
    }
  };

  const handleComplete = async () => {
    // Analytics: onboarding_completed
    analytics.trackEvent('onboarding_completed', {
      parent_user_id: platform.userId,
      child_profile_id: state.childProfileId,
      grade: state.grade,
    });

    // Update profile store
    setProfile({
      id: state.childProfileId!,
      displayName: state.displayName,
      grade: state.grade!,
      avatarId: state.avatarId!,
    });

    // Navigate to home
    navigate(ROUTES.HOME);
  };

  const progress = getProgress(state.currentStep);

  return (
    <div className={styles.container}>
      <ProgressBar
        value={progress}
        className={styles.progressBar}
        size="sm"
      />

      <div className={styles.content}>
        {state.currentStep === 'welcome' && (
          <WelcomeScreen onContinue={goToNextStep} />
        )}

        {state.currentStep === 'grade' && (
          <GradeSelection
            selectedGrade={state.grade}
            onSelect={(grade) => {
              handleStepComplete({ grade });
              goToNextStep();
            }}
            onBack={goToPreviousStep}
          />
        )}

        {state.currentStep === 'avatar' && (
          <AvatarSelection
            selectedAvatarId={state.avatarId}
            onSelect={(avatarId) => {
              handleStepComplete({ avatarId });
              goToNextStep();
            }}
            onBack={goToPreviousStep}
          />
        )}

        {state.currentStep === 'display_name' && (
          <DisplayNameInput
            value={state.displayName}
            onSubmit={(displayName) => {
              handleStepComplete({ displayName });
              goToNextStep();
            }}
            onBack={goToPreviousStep}
          />
        )}

        {state.currentStep === 'consent' && (
          <ConsentScreen
            adultConsent={state.adultConsent}
            privacyAccepted={state.privacyAccepted}
            termsAccepted={state.termsAccepted}
            onSubmit={(data) => {
              handleStepComplete(data);
              goToNextStep();
            }}
            onBack={goToPreviousStep}
          />
        )}

        {state.currentStep === 'email' && (
          <EmailInput
            email={state.email}
            onSubmit={(email) => {
              handleStepComplete({ email });
              goToNextStep();
            }}
            onBack={goToPreviousStep}
          />
        )}

        {state.currentStep === 'email_verification' && (
          <EmailVerificationWaiting
            email={state.email}
            onVerified={() => {
              handleStepComplete({ emailVerified: true });
              goToNextStep();
            }}
            onResend={() => {
              // Resend verification email
            }}
          />
        )}

        {state.currentStep === 'completed' && (
          <OnboardingCompletedScreen onContinue={handleComplete} />
        )}
      </div>
    </div>
  );
}

function getProgress(step: OnboardingStep): number {
  const steps: OnboardingStep[] = [
    'welcome',
    'grade',
    'avatar',
    'display_name',
    'consent',
    'email',
    'email_verification',
    'completed',
  ];

  const index = steps.indexOf(step);
  return ((index + 1) / steps.length) * 100;
}
```

---

### 3.2. GradeSelection Screen

**Файл:** `src/pages/Onboarding/screens/GradeSelection.tsx`

```typescript
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { useAnalytics } from '@/hooks/useAnalytics';
import styles from './GradeSelection.module.css';

interface GradeSelectionProps {
  selectedGrade: number | null;
  onSelect: (grade: number) => void;
  onBack: () => void;
}

export function GradeSelection({
  selectedGrade,
  onSelect,
  onBack,
}: GradeSelectionProps) {
  const analytics = useAnalytics();
  const grades = Array.from({ length: 11 }, (_, i) => i + 1);

  const handleSelect = (grade: number) => {
    // Analytics: grade_selected
    analytics.trackEvent('grade_selected', {
      grade,
    });

    onSelect(grade);
  };

  return (
    <div className={styles.container}>
      <button className={styles.backButton} onClick={onBack}>
        ← Назад
      </button>

      <h1 className={styles.title}>В каком ты классе?</h1>
      <p className={styles.subtitle}>
        Мы подберём задания под твой уровень
      </p>

      <div className={styles.grid}>
        {grades.map((grade) => (
          <Card
            key={grade}
            className={styles.gradeCard}
            variant={selectedGrade === grade ? 'primary' : 'bordered'}
            onClick={() => handleSelect(grade)}
          >
            <span className={styles.gradeNumber}>{grade}</span>
            <span className={styles.gradeLabel}>класс</span>
          </Card>
        ))}
      </div>
    </div>
  );
}
```

**Файл:** `src/pages/Onboarding/screens/GradeSelection.module.css`

```css
.container {
  display: flex;
  flex-direction: column;
  padding: 20px;
  min-height: 100vh;
}

.backButton {
  align-self: flex-start;
  background: none;
  border: none;
  color: #666;
  font-size: 16px;
  cursor: pointer;
  padding: 8px;
  margin-bottom: 16px;
}

.title {
  font-size: 28px;
  font-weight: 700;
  color: #000;
  margin: 0 0 8px 0;
  text-align: center;
}

.subtitle {
  font-size: 16px;
  color: #666;
  margin: 0 0 32px 0;
  text-align: center;
}

.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.gradeCard {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
  min-height: 100px;
}

.gradeCard:active {
  transform: scale(0.95);
}

.gradeNumber {
  font-size: 32px;
  font-weight: 700;
  color: #000;
  line-height: 1;
}

.gradeLabel {
  font-size: 14px;
  color: #666;
  margin-top: 4px;
}
```

---

### 3.3. AvatarSelection Screen

**Файл:** `src/pages/Onboarding/screens/AvatarSelection.tsx`

```typescript
import { useEffect, useState } from 'react';
import { Card } from '@/components/ui/Card';
import { Spinner } from '@/components/ui/Spinner';
import { useAnalytics } from '@/hooks/useAnalytics';
import { onboardingAPI } from '@/api/onboarding';
import type { Avatar } from '@/types/onboarding';
import styles from './AvatarSelection.module.css';

interface AvatarSelectionProps {
  selectedAvatarId: string | null;
  onSelect: (avatarId: string) => void;
  onBack: () => void;
}

export function AvatarSelection({
  selectedAvatarId,
  onSelect,
  onBack,
}: AvatarSelectionProps) {
  const analytics = useAnalytics();
  const [avatars, setAvatars] = useState<Avatar[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadAvatars();
  }, []);

  const loadAvatars = async () => {
    try {
      const data = await onboardingAPI.getAvatars();
      setAvatars(data);
    } catch (error) {
      console.error('[AvatarSelection] Failed to load avatars:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSelect = (avatar: Avatar) => {
    // Analytics: avatar_selected
    analytics.trackEvent('avatar_selected', {
      avatar_id: avatar.id,
    });

    onSelect(avatar.id);
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <button className={styles.backButton} onClick={onBack}>
        ← Назад
      </button>

      <h1 className={styles.title}>Выбери своего героя</h1>
      <p className={styles.subtitle}>
        Он будет помогать тебе в обучении
      </p>

      <div className={styles.grid}>
        {avatars.map((avatar) => (
          <Card
            key={avatar.id}
            className={styles.avatarCard}
            variant={selectedAvatarId === avatar.id ? 'primary' : 'bordered'}
            onClick={() => handleSelect(avatar)}
          >
            <img
              src={avatar.imageUrl}
              alt={avatar.name}
              className={styles.avatarImage}
            />
            {avatar.isPremium && (
              <span className={styles.premiumBadge}>⭐</span>
            )}
          </Card>
        ))}
      </div>
    </div>
  );
}
```

---

### 3.4. DisplayNameInput Screen

**Файл:** `src/pages/Onboarding/screens/DisplayNameInput.tsx`

```typescript
import { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { useAnalytics } from '@/hooks/useAnalytics';
import styles from './DisplayNameInput.module.css';

interface DisplayNameInputProps {
  value: string;
  onSubmit: (name: string) => void;
  onBack: () => void;
}

export function DisplayNameInput({
  value,
  onSubmit,
  onBack,
}: DisplayNameInputProps) {
  const analytics = useAnalytics();
  const [name, setName] = useState(value);
  const [error, setError] = useState('');

  const handleSubmit = () => {
    if (!name.trim()) {
      setError('Введи своё имя');
      return;
    }

    if (name.length < 2) {
      setError('Имя слишком короткое');
      return;
    }

    if (name.length > 20) {
      setError('Имя слишком длинное (максимум 20 символов)');
      return;
    }

    // Analytics: display_name_entered
    analytics.trackEvent('display_name_entered', {
      name_length: name.length,
    });

    onSubmit(name.trim());
  };

  return (
    <div className={styles.container}>
      <button className={styles.backButton} onClick={onBack}>
        ← Назад
      </button>

      <h1 className={styles.title}>Как тебя зовут?</h1>
      <p className={styles.subtitle}>
        Твоё имя будет видно только тебе и родителям
      </p>

      <div className={styles.form}>
        <Input
          value={name}
          onChange={(e) => {
            setName(e.target.value);
            setError('');
          }}
          placeholder="Твоё имя"
          maxLength={20}
          error={error}
          autoFocus
          onKeyPress={(e) => {
            if (e.key === 'Enter') {
              handleSubmit();
            }
          }}
        />

        <Button
          variant="primary"
          size="lg"
          isFullWidth
          onClick={handleSubmit}
          disabled={!name.trim()}
        >
          Продолжить
        </Button>
      </div>
    </div>
  );
}
```

---

### 3.5. ConsentScreen

**Файл:** `src/pages/Onboarding/screens/ConsentScreen.tsx`

```typescript
import { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Checkbox } from '@/components/ui/Checkbox';
import { useAnalytics } from '@/hooks/useAnalytics';
import { CONFIG } from '@/config/constants';
import styles from './ConsentScreen.module.css';

interface ConsentScreenProps {
  adultConsent: boolean;
  privacyAccepted: boolean;
  termsAccepted: boolean;
  onSubmit: (data: {
    adultConsent: boolean;
    privacyAccepted: boolean;
    termsAccepted: boolean;
  }) => void;
  onBack: () => void;
}

export function ConsentScreen({
  adultConsent: initialConsent,
  privacyAccepted: initialPrivacy,
  termsAccepted: initialTerms,
  onSubmit,
  onBack,
}: ConsentScreenProps) {
  const analytics = useAnalytics();
  const [adultConsent, setAdultConsent] = useState(initialConsent);
  const [privacyAccepted, setPrivacyAccepted] = useState(initialPrivacy);
  const [termsAccepted, setTermsAccepted] = useState(initialTerms);

  const handleAdultConsentChange = (checked: boolean) => {
    setAdultConsent(checked);

    // Analytics: adult_consent_checked
    analytics.trackEvent('adult_consent_checked', {
      checked,
    });
  };

  const handlePrivacyClick = () => {
    // Analytics: privacy_policy_opened
    analytics.trackEvent('privacy_policy_opened', {
      platform_type: 'web',
      session_id: analytics.sessionId,
      policy_version: CONFIG.PRIVACY_POLICY_VERSION,
    });

    // Open privacy policy
    window.open(CONFIG.PRIVACY_POLICY_URL, '_blank');
  };

  const handleTermsClick = () => {
    // Analytics: terms_opened
    analytics.trackEvent('terms_opened', {
      platform_type: 'web',
      session_id: analytics.sessionId,
      terms_version: CONFIG.TERMS_VERSION,
    });

    // Open terms
    window.open(CONFIG.TERMS_URL, '_blank');
  };

  const handleSubmit = () => {
    if (privacyAccepted) {
      // Analytics: privacy_policy_accepted
      analytics.trackEvent('privacy_policy_accepted', {
        policy_version: CONFIG.PRIVACY_POLICY_VERSION,
      });
    }

    if (termsAccepted) {
      // Analytics: terms_accepted
      analytics.trackEvent('terms_accepted', {
        terms_version: CONFIG.TERMS_VERSION,
      });
    }

    onSubmit({ adultConsent, privacyAccepted, termsAccepted });
  };

  const canSubmit = adultConsent && privacyAccepted && termsAccepted;

  return (
    <div className={styles.container}>
      <button className={styles.backButton} onClick={onBack}>
        ← Назад
      </button>

      <h1 className={styles.title}>Согласие на обработку данных</h1>
      <p className={styles.subtitle}>
        Для использования приложения необходимо согласие взрослого
      </p>

      <div className={styles.content}>
        <Checkbox
          checked={adultConsent}
          onChange={handleAdultConsentChange}
          label="Я являюсь родителем/законным представителем и даю согласие на использование приложения"
        />

        <Checkbox
          checked={privacyAccepted}
          onChange={setPrivacyAccepted}
          label={
            <span>
              Я принимаю{' '}
              <button
                className={styles.link}
                onClick={handlePrivacyClick}
              >
                политику конфиденциальности
              </button>
            </span>
          }
        />

        <Checkbox
          checked={termsAccepted}
          onChange={setTermsAccepted}
          label={
            <span>
              Я принимаю{' '}
              <button
                className={styles.link}
                onClick={handleTermsClick}
              >
                пользовательское соглашение
              </button>
            </span>
          }
        />
      </div>

      <Button
        variant="primary"
        size="lg"
        isFullWidth
        onClick={handleSubmit}
        disabled={!canSubmit}
      >
        Продолжить
      </Button>
    </div>
  );
}
```

---

### 3.6. EmailInput Screen

**Файл:** `src/pages/Onboarding/screens/EmailInput.tsx`

```typescript
import { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { useAnalytics } from '@/hooks/useAnalytics';
import { onboardingAPI } from '@/api/onboarding';
import styles from './EmailInput.module.css';

interface EmailInputProps {
  email: string;
  onSubmit: (email: string) => void;
  onBack: () => void;
}

export function EmailInput({ email, onSubmit, onBack }: EmailInputProps) {
  const analytics = useAnalytics();
  const [value, setValue] = useState(email);
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const validateEmail = (email: string): boolean => {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
  };

  const handleSubmit = async () => {
    if (!value.trim()) {
      setError('Введите email');
      return;
    }

    if (!validateEmail(value)) {
      setError('Некорректный email');
      return;
    }

    const emailDomain = value.split('@')[1];

    // Analytics: email_entered
    analytics.trackEvent('email_entered', {
      email_domain: emailDomain,
    });

    setIsSubmitting(true);

    try {
      await onboardingAPI.sendEmailVerification(value);

      // Analytics: email_verification_sent (backend should send this)
      // But we track it here for redundancy
      analytics.trackEvent('email_verification_sent', {
        email_domain: emailDomain,
      });

      onSubmit(value);
    } catch (err) {
      setError('Не удалось отправить письмо. Попробуйте ещё раз');
      console.error('[EmailInput] Failed to send verification:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className={styles.container}>
      <button className={styles.backButton} onClick={onBack}>
        ← Назад
      </button>

      <h1 className={styles.title}>Email родителя</h1>
      <p className={styles.subtitle}>
        Мы отправим письмо для подтверждения
      </p>

      <div className={styles.form}>
        <Input
          type="email"
          value={value}
          onChange={(e) => {
            setValue(e.target.value);
            setError('');
          }}
          placeholder="parent@example.com"
          error={error}
          autoFocus
          disabled={isSubmitting}
          onKeyPress={(e) => {
            if (e.key === 'Enter') {
              handleSubmit();
            }
          }}
        />

        <Button
          variant="primary"
          size="lg"
          isFullWidth
          onClick={handleSubmit}
          disabled={!value.trim() || isSubmitting}
          isLoading={isSubmitting}
        >
          Отправить письмо
        </Button>
      </div>
    </div>
  );
}
```

---

### 3.7. EmailVerificationWaiting Screen

**Файл:** `src/pages/Onboarding/screens/EmailVerificationWaiting.tsx`

```typescript
import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Spinner } from '@/components/ui/Spinner';
import { useAnalytics } from '@/hooks/useAnalytics';
import { onboardingAPI } from '@/api/onboarding';
import styles from './EmailVerificationWaiting.module.css';

interface EmailVerificationWaitingProps {
  email: string;
  onVerified: () => void;
  onResend: () => void;
}

export function EmailVerificationWaiting({
  email,
  onVerified,
  onResend,
}: EmailVerificationWaitingProps) {
  const analytics = useAnalytics();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    // Poll verification status every 3 seconds
    const interval = setInterval(async () => {
      try {
        const { verified } = await onboardingAPI.checkEmailVerification(email);

        if (verified) {
          clearInterval(interval);

          // Analytics: email_verification_success (backend should send this)
          analytics.trackEvent('email_verification_success', {
            email_domain: email.split('@')[1],
          });

          onVerified();
        }
      } catch (error) {
        console.error('[EmailVerification] Check failed:', error);
      }
    }, 3000);

    return () => clearInterval(interval);
  }, [email, onVerified, analytics]);

  return (
    <div className={styles.container}>
      <div className={styles.icon}>
        <Spinner size="lg" />
      </div>

      <h1 className={styles.title}>Проверьте почту</h1>
      <p className={styles.subtitle}>
        Мы отправили письмо на <strong>{email}</strong>
      </p>
      <p className={styles.instruction}>
        Перейдите по ссылке в письме для подтверждения
      </p>

      <div className={styles.resendSection}>
        <p className={styles.resendText}>Не пришло письмо?</p>
        <Button
          variant="outline"
          size="md"
          onClick={onResend}
        >
          Отправить ещё раз
        </Button>
      </div>
    </div>
  );
}
```

---

## Часть 4: Аналитические события

### События онбординга

| Event Name | Когда срабатывает | Параметры |
|------------|-------------------|-----------|
| `onboarding_opened` | Открыт первый экран онбординга | `platform_type`, `session_id`, `entry_point` |
| `registration_opened` | Открыт экран регистрации | `platform_type`, `session_id` |
| `consent_screen_opened` | Открыт экран согласий | `platform_type`, `session_id` |
| `grade_selected` | Выбран класс | `grade` |
| `avatar_selected` | Выбран аватар | `avatar_id` |
| `display_name_entered` | Введено имя | `name_length` |
| `adult_consent_checked` | Изменено согласие взрослого | `checked` |
| `privacy_policy_opened` | Открыта политика | `platform_type`, `session_id`, `policy_version` |
| `privacy_policy_accepted` | Принята политика | `policy_version` |
| `terms_opened` | Открыто соглашение | `platform_type`, `session_id`, `terms_version` |
| `terms_accepted` | Принято соглашение | `terms_version` |
| `email_entered` | Введён email | `email_domain` |
| `email_verification_sent` | Отправлено письмо | `email_domain` |
| `email_verification_success` | Email подтверждён | `parent_user_id`, `email_domain` |
| `onboarding_completed` | Онбординг завершён | `parent_user_id`, `child_profile_id`, `grade` |

---

## Часть 5: API Endpoints

### Backend эндпоинты

```typescript
// GET /api/v1/avatars
// Response:
interface AvatarsResponse {
  avatars: Array<{
    id: string;
    imageUrl: string;
    name: string;
    isPremium: boolean;
  }>;
}

// POST /api/v1/profiles/child
// Request:
interface CreateChildProfileRequest {
  parentUserId: string;
  grade: number;
  avatarId: string;
  displayName: string;
}
// Response:
interface CreateChildProfileResponse {
  childProfileId: string;
}

// POST /api/v1/auth/email/send-verification
// Request:
interface SendVerificationRequest {
  email: string;
}
// Response: 204 No Content

// GET /api/v1/auth/email/check-verification?email=xxx
// Response:
interface CheckVerificationResponse {
  verified: boolean;
}

// POST /api/v1/consent
// Request:
interface SaveConsentRequest {
  parentUserId: string;
  privacyPolicyVersion: string;
  termsVersion: string;
  adultConsent: boolean;
}
// Response: 204 No Content

// POST /api/v1/onboarding/complete
// Request:
interface CompleteOnboardingRequest {
  parentUserId: string;
  childProfileId: string;
}
// Response: 204 No Content
```

---

## Часть 6: Конфигурация

### Constants

**Файл:** `src/config/constants.ts`

```typescript
export const CONFIG = {
  PRIVACY_POLICY_URL: 'https://example.com/privacy',
  PRIVACY_POLICY_VERSION: '1.0.0',
  TERMS_URL: 'https://example.com/terms',
  TERMS_VERSION: '1.0.0',
  MAX_DISPLAY_NAME_LENGTH: 20,
  MIN_DISPLAY_NAME_LENGTH: 2,
  EMAIL_VERIFICATION_POLL_INTERVAL: 3000, // ms
};
```

---

## Чеклист задач

### Типы и API
- [ ] Создать типы onboarding.ts
- [ ] Реализовать onboardingAPI.getAvatars()
- [ ] Реализовать onboardingAPI.createChildProfile()
- [ ] Реализовать onboardingAPI.sendEmailVerification()
- [ ] Реализовать onboardingAPI.checkEmailVerification()
- [ ] Реализовать onboardingAPI.saveConsent()
- [ ] Реализовать onboardingAPI.completeOnboarding()

### Компоненты
- [ ] Создать OnboardingFlow container
- [ ] Создать WelcomeScreen
- [ ] Создать GradeSelection (grid 3x11)
- [ ] Создать AvatarSelection (grid 4x3)
- [ ] Создать DisplayNameInput
- [ ] Создать ConsentScreen с checkbox и ссылками
- [ ] Создать EmailInput
- [ ] Создать EmailVerificationWaiting (polling)
- [ ] Создать OnboardingCompletedScreen

### Аналитика
- [ ] Добавить onboarding_opened
- [ ] Добавить grade_selected
- [ ] Добавить avatar_selected
- [ ] Добавить display_name_entered
- [ ] Добавить adult_consent_checked
- [ ] Добавить privacy_policy_opened/accepted
- [ ] Добавить terms_opened/accepted
- [ ] Добавить email_entered
- [ ] Добавить email_verification_sent
- [ ] Добавить onboarding_completed

### UI/UX
- [ ] Добавить прогресс-бар вверху
- [ ] Добавить кнопку "Назад" на каждом экране
- [ ] Добавить валидацию имени (2-20 символов)
- [ ] Добавить валидацию email
- [ ] Добавить анимации переходов между экранами
- [ ] Добавить обработку ошибок сети

### Тестирование
- [ ] Протестировать весь flow от начала до конца
- [ ] Протестировать валидацию полей
- [ ] Протестировать email верификацию
- [ ] Протестировать возврат назад на каждом шаге
- [ ] Протестировать на разных устройствах

---

## Следующий этап

После завершения онбординга переходи к **04_HOME.md** для создания главного экрана.
