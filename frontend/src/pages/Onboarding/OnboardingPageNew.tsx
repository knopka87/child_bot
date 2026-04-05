// src/pages/Onboarding/OnboardingPageNew.tsx
import { useState, useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronLeft, Check } from 'lucide-react';
import { ROUTES } from '@/config/routes';
import { onboardingAPI } from '@/api/onboarding';
import { PlatformBridge } from '@/services/platform/PlatformBridge';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { useAnalytics } from '@/hooks/useAnalytics';
import { LegalDocumentModal } from '@/components/LegalDocumentModal';

type OnboardingStep = 'grade' | 'avatar' | 'email' | 'email_verification' | 'consent' | 'completed';

const avatars = [
  { id: 'cat', emoji: '🐱', name: 'Кот' },
  { id: 'dog', emoji: '🐶', name: 'Пёс' },
  { id: 'panda', emoji: '🐼', name: 'Панда' },
  { id: 'fox', emoji: '🦊', name: 'Лиса' },
  { id: 'bear', emoji: '🐻', name: 'Медведь' },
  { id: 'lion', emoji: '🦁', name: 'Лев' },
  { id: 'tiger', emoji: '🐯', name: 'Тигр' },
  { id: 'unicorn', emoji: '🦄', name: 'Единорог' },
  { id: 'robot', emoji: '🤖', name: 'Робот' },
  { id: 'alien', emoji: '👽', name: 'Инопланетянин' },
];

const grades = [1, 2, 3, 4];

export function OnboardingPageNew() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const analytics = useAnalytics();

  const [currentStep, setCurrentStep] = useState<OnboardingStep>('grade');
  const [grade, setGrade] = useState<number | null>(null);
  const [avatarId, setAvatarId] = useState<string | null>(null);
  const [adultConsent, setAdultConsent] = useState(false);
  const [privacyAccepted, setPrivacyAccepted] = useState(false);
  const [termsAccepted, setTermsAccepted] = useState(false);
  const [displayName, setDisplayName] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [verificationCode, setVerificationCode] = useState<string>('');
  const [devCode, setDevCode] = useState<string>(''); // Код для разработки (показываем в dev режиме)
  const [isEmailVerified, setIsEmailVerified] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [referralCode, setReferralCode] = useState<string | null>(null);
  const [isInitialized, setIsInitialized] = useState(false); // Флаг завершения инициализации
  const [legalModalType, setLegalModalType] = useState<'privacy' | 'terms' | null>(null); // Тип открытого legal-документа
  const hasInitialized = useRef(false); // Ref для гарантии что инициализация произойдёт только один раз

  // Загружаем данные от VK и восстанавливаем прогресс онбординга при монтировании
  useEffect(() => {
    // Гарантируем что инициализация происходит только один раз
    if (hasInitialized.current) {
      console.log('[Onboarding] Skipping init - already initialized');
      return;
    }
    hasInitialized.current = true;

    const initOnboarding = async () => {
      try {
        // 1. Загружаем данные от VK
        const platformBridge = new PlatformBridge();
        await platformBridge.init();
        const user = await platformBridge.getUser();

        // Автоматически заполняем имя из VK
        const vkDisplayName = user.firstName || 'Ученик';

        console.log('[Onboarding] VK user data loaded:', { firstName: user.firstName });

        // 2. Восстанавливаем сохранённый прогресс онбординга
        const savedStep = await vkStorage.getItem(storageKeys.ONBOARDING_STEP);
        const savedGrade = await vkStorage.getItem(storageKeys.ONBOARDING_GRADE);
        const savedAvatar = await vkStorage.getItem(storageKeys.ONBOARDING_AVATAR);
        const savedEmail = await vkStorage.getItem(storageKeys.ONBOARDING_EMAIL);
        const savedEmailVerified = await vkStorage.getItem(storageKeys.ONBOARDING_EMAIL_VERIFIED);
        const savedDisplayName = await vkStorage.getItem(storageKeys.ONBOARDING_DISPLAY_NAME);
        const savedConsents = await vkStorage.getItem(storageKeys.ONBOARDING_CONSENTS);

        if (savedStep) {
          console.log('[Onboarding] Restoring progress from storage:', {
            step: savedStep,
            grade: savedGrade,
            avatar: savedAvatar,
            email: savedEmail,
            emailVerified: savedEmailVerified,
          });

          // Восстанавливаем состояние БАТЧЕМ (используем setTimeout для группировки обновлений)
          setTimeout(() => {
            setCurrentStep(savedStep as OnboardingStep);
            if (savedGrade) setGrade(parseInt(savedGrade, 10));
            if (savedAvatar) setAvatarId(savedAvatar);
            if (savedEmail) setEmail(savedEmail);
            if (savedEmailVerified === 'true') setIsEmailVerified(true);
            setDisplayName(savedDisplayName || vkDisplayName);

            // Восстанавливаем согласия
            if (savedConsents) {
              try {
                const consents = JSON.parse(savedConsents);
                setAdultConsent(consents.adultConsent || false);
                setPrivacyAccepted(consents.privacyAccepted || false);
                setTermsAccepted(consents.termsAccepted || false);
              } catch (e) {
                console.error('[Onboarding] Failed to parse saved consents:', e);
              }
            }

            // Отмечаем что инициализация завершена ПОСЛЕ восстановления
            setIsInitialized(true);
          }, 0);
        } else {
          console.log('[Onboarding] No saved progress found, starting fresh');
          setDisplayName(vkDisplayName);
          // Инициализация завершена сразу
          setIsInitialized(true);
        }

        // 3. Извлекаем реферальный код из URL
        const refCode = searchParams.get('ref');
        if (refCode) {
          console.log('[Onboarding] Referral code detected:', refCode);
          setReferralCode(refCode);
          vkStorage.setItem(storageKeys.REFERRAL_CODE, refCode);
        } else {
          const savedCode = await vkStorage.getItem(storageKeys.REFERRAL_CODE);
          if (savedCode) {
            console.log('[Onboarding] Referral code loaded from storage:', savedCode);
            setReferralCode(savedCode);
          }
        }

        // Analytics
        analytics.trackEvent('onboarding_opened', {});
      } catch (error) {
        console.error('[Onboarding] Failed to initialize:', error);
        setDisplayName('Ученик'); // fallback
        setIsInitialized(true);
      }
    };

    initOnboarding();
  }, [searchParams, analytics]);

  // Автоматическое сохранение прогресса при изменении любого поля
  useEffect(() => {
    // НЕ сохраняем пока не завершилась инициализация
    if (!isInitialized) {
      console.log('[Onboarding] Skipping auto-save - not initialized yet');
      return;
    }

    const saveProgress = async () => {
      try {
        // Сохраняем текущий шаг
        await vkStorage.setItem(storageKeys.ONBOARDING_STEP, currentStep);

        // Сохраняем данные профиля
        if (grade !== null) {
          await vkStorage.setItem(storageKeys.ONBOARDING_GRADE, grade.toString());
        }
        if (avatarId) {
          await vkStorage.setItem(storageKeys.ONBOARDING_AVATAR, avatarId);
        }
        if (email) {
          await vkStorage.setItem(storageKeys.ONBOARDING_EMAIL, email);
        }
        if (displayName) {
          await vkStorage.setItem(storageKeys.ONBOARDING_DISPLAY_NAME, displayName);
        }

        // Сохраняем статус верификации email
        await vkStorage.setItem(
          storageKeys.ONBOARDING_EMAIL_VERIFIED,
          isEmailVerified.toString()
        );

        // Сохраняем согласия
        const consents = {
          adultConsent,
          privacyAccepted,
          termsAccepted,
        };
        await vkStorage.setItem(storageKeys.ONBOARDING_CONSENTS, JSON.stringify(consents));

        console.log('[Onboarding] Progress saved:', {
          step: currentStep,
          grade,
          avatarId,
          email,
          emailVerified: isEmailVerified,
        });
      } catch (error) {
        console.error('[Onboarding] Failed to save progress:', error);
      }
    };

    // Сохраняем только если не на шаге completed
    if (currentStep !== 'completed') {
      saveProgress();
    }
  }, [currentStep, grade, avatarId, email, displayName, isEmailVerified, adultConsent, privacyAccepted, termsAccepted, isInitialized]);

  const handleBack = () => {
    const steps: OnboardingStep[] = ['grade', 'avatar', 'email', 'email_verification', 'consent'];
    const currentIndex = steps.indexOf(currentStep);
    if (currentIndex > 0) {
      setCurrentStep(steps[currentIndex - 1]);
    }
  };

  const handleNext = async () => {
    console.log('[Onboarding] handleNext called, currentStep:', currentStep);

    const steps: OnboardingStep[] = ['grade', 'avatar', 'email', 'email_verification', 'consent', 'completed'];
    const currentIndex = steps.indexOf(currentStep);

    console.log('[Onboarding] currentIndex:', currentIndex, 'steps.length:', steps.length);

    // Analytics для текущего шага
    if (currentStep === 'grade' && grade) {
      analytics.trackEvent('grade_selected', { grade });
    } else if (currentStep === 'avatar' && avatarId) {
      analytics.trackEvent('avatar_selected', { avatar_id: avatarId });
    } else if (currentStep === 'email' && email) {
      analytics.trackEvent('email_entered', { email_domain: email.split('@')[1] });
      // Отправляем код верификации
      await sendVerificationCode();
    } else if (currentStep === 'email_verification' && isEmailVerified) {
      analytics.trackEvent('email_verification_success', {});
    } else if (currentStep === 'consent') {
      console.log('[Onboarding] On consent step, tracking analytics...');
      analytics.trackEvent('adult_consent_checked', {});
      if (privacyAccepted) {
        analytics.trackEvent('privacy_policy_accepted', {});
      }
      if (termsAccepted) {
        analytics.trackEvent('terms_accepted', {});
      }
    }

    if (currentIndex < steps.length - 1) {
      const nextStep = steps[currentIndex + 1];
      console.log('[Onboarding] Moving to next step:', nextStep);
      setCurrentStep(nextStep);
    } else {
      console.log('[Onboarding] Already on last step, not moving');
    }
  };

  const sendVerificationCode = async () => {
    try {
      const platformBridge = new PlatformBridge();
      await platformBridge.init();
      const user = await platformBridge.getUser();

      const result = await onboardingAPI.sendEmailVerification({
        email,
        parentUserId: user.id,
      });

      // В dev режиме сохраняем код для отображения
      if (result.devCode) {
        setDevCode(result.devCode);
        console.log('[Onboarding] Dev code received:', result.devCode);
      }

      analytics.trackEvent('email_verification_sent', {});
    } catch (error) {
      console.error('[Onboarding] Failed to send verification:', error);
      alert('Не удалось отправить код подтверждения. Попробуйте ещё раз.');
    }
  };

  const verifyCode = async () => {
    try {
      const result = await onboardingAPI.verifyEmailCode({
        email,
        code: verificationCode,
      });

      if (result.verified) {
        setIsEmailVerified(true);
        analytics.trackEvent('email_verification_success', {});
        // Автоматически переходим к следующему шагу
        setTimeout(() => {
          setCurrentStep('consent');
        }, 1000);
      } else {
        alert('Неверный код. Проверьте правильность ввода.');
      }
    } catch (error) {
      console.error('[Onboarding] Failed to verify code:', error);
      alert('Не удалось проверить код. Попробуйте ещё раз.');
    }
  };

  const handleComplete = async () => {
    if (isSubmitting) return;

    try {
      setIsSubmitting(true);
      console.log('[Onboarding] Starting completion...');

      const platformBridge = new PlatformBridge();
      await platformBridge.init();
      const user = await platformBridge.getUser();
      const parentUserId = user.id;

      console.log('[Onboarding] Creating child profile...');
      const { childProfileId } = await onboardingAPI.createChildProfile({
        parentUserId,
        grade: grade!,
        avatarId: avatarId!,
        displayName: displayName!,
        referralCode: referralCode || undefined,
      });

      console.log('[Onboarding] Child profile created:', childProfileId);

      if (referralCode) {
        await vkStorage.removeItem(storageKeys.REFERRAL_CODE);
      }

      await vkStorage.setItem(storageKeys.USER_ID, parentUserId);
      await vkStorage.setItem(storageKeys.PROFILE_ID, childProfileId);
      await vkStorage.setItem(storageKeys.ONBOARDING_COMPLETED, 'true');

      await onboardingAPI.saveConsent({
        parentUserId,
        privacyPolicyVersion: '1.0',
        termsVersion: '1.0',
        adultConsent: adultConsent!,
      });

      // Очищаем временные данные онбординга после успешного завершения
      await vkStorage.removeItem(storageKeys.ONBOARDING_STEP);
      await vkStorage.removeItem(storageKeys.ONBOARDING_GRADE);
      await vkStorage.removeItem(storageKeys.ONBOARDING_AVATAR);
      await vkStorage.removeItem(storageKeys.ONBOARDING_EMAIL);
      await vkStorage.removeItem(storageKeys.ONBOARDING_EMAIL_VERIFIED);
      await vkStorage.removeItem(storageKeys.ONBOARDING_DISPLAY_NAME);
      await vkStorage.removeItem(storageKeys.ONBOARDING_CONSENTS);
      console.log('[Onboarding] Temporary onboarding data cleared');

      analytics.trackEvent('onboarding_completed', {
        child_profile_id: childProfileId,
        grade: grade!,
        avatar_id: avatarId!,
        has_referral: !!referralCode,
      });

      console.log('[Onboarding] Completed successfully!');
      navigate(ROUTES.HOME, { replace: true });
    } catch (error) {
      console.error('[Onboarding] Failed to complete:', error);
      alert('Не удалось завершить регистрацию. Попробуйте ещё раз.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const canProceed = () => {
    let result = false;
    switch (currentStep) {
      case 'grade':
        result = grade !== null;
        break;
      case 'avatar':
        result = avatarId !== null;
        break;
      case 'email':
        result = email.includes('@') && email.includes('.');
        break;
      case 'email_verification':
        result = isEmailVerified;
        break;
      case 'consent':
        result = adultConsent && privacyAccepted && termsAccepted;
        console.log('[Onboarding] canProceed for consent:', {
          adultConsent,
          privacyAccepted,
          termsAccepted,
          result,
        });
        break;
      case 'completed':
        result = true;
        break;
      default:
        result = false;
    }
    return result;
  };

  const getProgressInfo = () => {
    const steps: OnboardingStep[] = ['grade', 'avatar', 'email', 'email_verification', 'consent', 'completed'];
    const currentIndex = steps.indexOf(currentStep);
    const currentStepNumber = currentIndex + 1;
    const totalSteps = steps.length;
    const percent = (currentStepNumber / totalSteps) * 100;

    return {
      currentStep: currentStepNumber,
      totalSteps,
      percent,
    };
  };

  const handleLegalLinkClick = (type: 'privacy' | 'terms') => {
    setLegalModalType(type);
  };

  return (
    <>
      <LegalDocumentModal
        type={legalModalType}
        isOpen={legalModalType !== null}
        onClose={() => setLegalModalType(null)}
      />

      <div className="min-h-screen bg-gradient-to-br from-[#E8E4FF] to-[#F0ECFF] flex flex-col">
        {/* Header с прогрессом */}
      <div className="bg-white shadow-sm">
        <div className="flex items-center justify-between px-4 py-3">
          {currentStep !== 'grade' && currentStep !== 'completed' && (
            <button
              onClick={handleBack}
              className="p-2 -ml-2 text-[#6C5CE7] active:scale-95 transition-transform"
            >
              <ChevronLeft size={24} />
            </button>
          )}
          <div className="flex-1 mx-4">
            <div className="h-1.5 bg-gray-200 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: `${getProgressInfo().percent}%` }}
                transition={{ duration: 0.3 }}
                className="h-full bg-gradient-to-r from-[#6C5CE7] to-[#5B4FDB]"
              />
            </div>
          </div>
          <div className="text-sm text-gray-500 font-medium min-w-[70px] text-right">
            Шаг {getProgressInfo().currentStep} из {getProgressInfo().totalSteps}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto px-4 py-6">
        <AnimatePresence mode="wait">
          {/* Grade Selection */}
          {currentStep === 'grade' && (
            <motion.div
              key="grade"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.3 }}
            >
              <h2 className="text-2xl font-bold text-[#2D3436] mb-2">
                В каком классе учится ребёнок?
              </h2>
              <p className="text-gray-500 mb-6">
                Это поможет нам подобрать подходящие задания
              </p>

              <div className="grid grid-cols-3 gap-3">
                {grades.map((g) => (
                  <button
                    key={g}
                    onClick={() => setGrade(g)}
                    className={`py-4 rounded-2xl font-semibold text-lg transition-all ${
                      grade === g
                        ? 'bg-[#6C5CE7] text-white shadow-lg scale-105'
                        : 'bg-white text-gray-700 hover:bg-gray-50 active:scale-95'
                    }`}
                  >
                    {g}
                  </button>
                ))}
              </div>
            </motion.div>
          )}

          {/* Email Input */}
          {currentStep === 'email' && (
            <motion.div
              key="email"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.3 }}
            >
              <h2 className="text-2xl font-bold text-[#2D3436] mb-2">
                Email родителя
              </h2>
              <p className="text-gray-500 mb-6">
                Для отправки отчётов об успеваемости
              </p>

              <div className="space-y-4">
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="parent@example.com"
                  className="w-full px-4 py-4 bg-white rounded-2xl border-2 border-gray-200 focus:border-[#6C5CE7] focus:outline-none text-lg"
                  autoComplete="email"
                />

                <div className="bg-blue-50 rounded-2xl p-4 mt-4">
                  <p className="text-sm text-gray-600">
                    💡 На этот email мы будем отправлять отчёты о прогрессе ребёнка и достижениях
                  </p>
                </div>
              </div>
            </motion.div>
          )}

          {/* Email Verification */}
          {currentStep === 'email_verification' && (
            <motion.div
              key="email_verification"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.3 }}
            >
              {!isEmailVerified ? (
                <>
                  <h2 className="text-2xl font-bold text-[#2D3436] mb-2">
                    Проверка email
                  </h2>
                  <p className="text-gray-500 mb-6">
                    Мы отправили код подтверждения на <strong>{email}</strong>
                  </p>

                  <div className="space-y-4">
                    <input
                      type="text"
                      value={verificationCode}
                      onChange={(e) => setVerificationCode(e.target.value.slice(0, 6))}
                      placeholder="000000"
                      maxLength={6}
                      className="w-full px-4 py-4 bg-white rounded-2xl border-2 border-gray-200 focus:border-[#6C5CE7] focus:outline-none text-2xl text-center tracking-[0.5em] font-mono"
                      autoComplete="one-time-code"
                    />

                    {verificationCode.length === 6 && (
                      <button
                        onClick={verifyCode}
                        className="w-full py-4 bg-[#6C5CE7] text-white rounded-2xl font-semibold active:scale-[0.98] transition-transform"
                      >
                        Проверить код
                      </button>
                    )}

                    <button
                      onClick={sendVerificationCode}
                      className="w-full py-3 text-[#6C5CE7] font-medium active:scale-95 transition-transform"
                    >
                      Отправить код повторно
                    </button>

                    <div className="bg-yellow-50 rounded-2xl p-4 mt-4">
                      <p className="text-sm text-gray-600">
                        ⏱️ Код действителен 15 минут. Проверьте папку «Спам», если письмо не пришло.
                      </p>
                    </div>

                    {/* Dev mode: показываем код для тестирования */}
                    {devCode && (
                      <div className="bg-green-50 border-2 border-green-200 rounded-2xl p-4 mt-4">
                        <p className="text-sm text-green-800 font-semibold mb-2">
                          🔧 Режим разработки
                        </p>
                        <p className="text-sm text-green-700 mb-2">
                          Письма пока не отправляются. Используйте код:
                        </p>
                        <div className="bg-white rounded-xl p-3 text-center">
                          <p className="text-3xl font-mono font-bold text-green-600 tracking-[0.5em]">
                            {devCode}
                          </p>
                        </div>
                      </div>
                    )}
                  </div>
                </>
              ) : (
                <div className="text-center py-8">
                  <div className="w-20 h-20 mx-auto mb-4 bg-green-100 rounded-full flex items-center justify-center">
                    <Check size={40} className="text-green-600" />
                  </div>
                  <h2 className="text-2xl font-bold text-[#2D3436] mb-2">
                    Email подтверждён!
                  </h2>
                  <p className="text-gray-500">
                    Переходим к следующему шагу...
                  </p>
                </div>
              )}
            </motion.div>
          )}

          {/* Avatar Selection */}
          {currentStep === 'avatar' && (
            <motion.div
              key="avatar"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.3 }}
            >
              <h2 className="text-2xl font-bold text-[#2D3436] mb-2">
                Выбери аватар для {displayName}
              </h2>
              <p className="text-gray-500 mb-6">
                Ты сможешь изменить его в любое время
              </p>

              <div className="grid grid-cols-3 gap-3">
                {avatars.map((avatar) => (
                  <button
                    key={avatar.id}
                    onClick={() => setAvatarId(avatar.id)}
                    className={`p-4 rounded-2xl transition-all ${
                      avatarId === avatar.id
                        ? 'bg-[#6C5CE7] shadow-lg scale-105'
                        : 'bg-white hover:bg-gray-50 active:scale-95'
                    }`}
                  >
                    <div className="text-5xl mb-2">{avatar.emoji}</div>
                    <div
                      className={`text-sm font-medium ${
                        avatarId === avatar.id ? 'text-white' : 'text-gray-700'
                      }`}
                    >
                      {avatar.name}
                    </div>
                  </button>
                ))}
              </div>
            </motion.div>
          )}

          {/* Consent Screen */}
          {currentStep === 'consent' && (
            <motion.div
              key="consent"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.3 }}
            >
              <h2 className="text-2xl font-bold text-[#2D3436] mb-2">
                Согласия
              </h2>
              <p className="text-gray-500 mb-6">
                Необходимо согласие взрослого для использования приложения
              </p>

              <div className="space-y-4">
                {/* Adult Consent */}
                <label className="flex items-start gap-3 p-4 bg-white rounded-2xl cursor-pointer active:scale-[0.98] transition-transform">
                  <div className="relative flex-shrink-0 mt-0.5">
                    <input
                      type="checkbox"
                      checked={adultConsent}
                      onChange={(e) => setAdultConsent(e.target.checked)}
                      className="sr-only peer"
                    />
                    <div className="w-6 h-6 border-2 border-gray-300 rounded-md peer-checked:bg-[#6C5CE7] peer-checked:border-[#6C5CE7] flex items-center justify-center transition-all">
                      {adultConsent && <Check size={16} className="text-white" />}
                    </div>
                  </div>
                  <span className="text-gray-700 flex-1">
                    Я являюсь родителем или законным представителем ребёнка
                  </span>
                </label>

                {/* Privacy Policy */}
                <label className="flex items-start gap-3 p-4 bg-white rounded-2xl cursor-pointer active:scale-[0.98] transition-transform">
                  <div className="relative flex-shrink-0 mt-0.5">
                    <input
                      type="checkbox"
                      checked={privacyAccepted}
                      onChange={(e) => setPrivacyAccepted(e.target.checked)}
                      className="sr-only peer"
                    />
                    <div className="w-6 h-6 border-2 border-gray-300 rounded-md peer-checked:bg-[#6C5CE7] peer-checked:border-[#6C5CE7] flex items-center justify-center transition-all">
                      {privacyAccepted && <Check size={16} className="text-white" />}
                    </div>
                  </div>
                  <span className="text-gray-700 flex-1">
                    Я принимаю{' '}
                    <button
                      type="button"
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation(); // Не даём клику всплыть на label
                        analytics.trackEvent('privacy_policy_opened', {});
                        handleLegalLinkClick('privacy');
                      }}
                      className="text-[#6C5CE7] underline"
                    >
                      политику конфиденциальности
                    </button>
                  </span>
                </label>

                {/* Terms of Service */}
                <label className="flex items-start gap-3 p-4 bg-white rounded-2xl cursor-pointer active:scale-[0.98] transition-transform">
                  <div className="relative flex-shrink-0 mt-0.5">
                    <input
                      type="checkbox"
                      checked={termsAccepted}
                      onChange={(e) => setTermsAccepted(e.target.checked)}
                      className="sr-only peer"
                    />
                    <div className="w-6 h-6 border-2 border-gray-300 rounded-md peer-checked:bg-[#6C5CE7] peer-checked:border-[#6C5CE7] flex items-center justify-center transition-all">
                      {termsAccepted && <Check size={16} className="text-white" />}
                    </div>
                  </div>
                  <span className="text-gray-700 flex-1">
                    Я принимаю{' '}
                    <button
                      type="button"
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation(); // Не даём клику всплыть на label
                        analytics.trackEvent('terms_opened', {});
                        handleLegalLinkClick('terms');
                      }}
                      className="text-[#6C5CE7] underline"
                    >
                      условия использования
                    </button>
                  </span>
                </label>
              </div>
            </motion.div>
          )}

          {/* Completed Screen */}
          {currentStep === 'completed' && (
            <motion.div
              key="completed"
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.4 }}
              className="text-center py-12"
            >
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: 0.2, type: 'spring', stiffness: 200 }}
                className="w-24 h-24 mx-auto mb-6 bg-green-100 rounded-full flex items-center justify-center"
              >
                <Check size={48} className="text-green-600" />
              </motion.div>

              <h2 className="text-2xl font-bold text-[#2D3436] mb-2">
                Всё готово, {displayName}!
              </h2>
              <p className="text-gray-500 mb-8">
                Профиль создан. Можно начинать учиться!
              </p>

              <button
                onClick={handleComplete}
                disabled={isSubmitting}
                className="w-full py-4 bg-[#6C5CE7] text-white rounded-2xl font-semibold active:scale-[0.98] transition-transform disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? 'Загрузка...' : 'Начать'}
              </button>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* Bottom Button */}
      {currentStep !== 'completed' && (
        <div className="p-4 bg-white shadow-lg">
          <button
            onClick={() => {
              console.log('[Onboarding] Button clicked!', { currentStep, canProceed: canProceed() });
              handleNext();
            }}
            disabled={!canProceed()}
            className="w-full py-4 bg-[#6C5CE7] text-white rounded-2xl font-semibold active:scale-[0.98] transition-transform disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100"
          >
            Далее
          </button>
        </div>
      )}
      </div>
    </>
  );
}
