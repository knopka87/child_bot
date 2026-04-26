// src/pages/Onboarding/OnboardingPageNew.tsx
import { useState, useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Check, ArrowRight, Shield } from 'lucide-react';
import { ROUTES } from '@/config/routes';
import { onboardingAPI } from '@/api/onboarding';
import { PlatformBridge } from '@/services/platform/PlatformBridge';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { useAnalytics } from '@/hooks/useAnalytics';
import { LegalDocumentModal } from '@/components/LegalDocumentModal';

type OnboardingStep = 'consent' | 'profile';

const avatars = [
  { id: 'cat', emoji: '🐱', name: 'Кот' },
  { id: 'dog', emoji: '🐶', name: 'Пёс' },
  { id: 'panda', emoji: '🐼', name: 'Панда' },
  { id: 'fox', emoji: '🦊', name: 'Лиса' },
  { id: 'bear', emoji: '🐻', name: 'Медведь' },
  { id: 'lion', emoji: '🦁', name: 'Лев' },
  { id: 'tiger', emoji: '🐯', name: 'Тигр' },
  { id: 'unicorn', emoji: '🦄', name: 'Единорог' },
];

const grades = [1, 2, 3, 4];

export function OnboardingPageNew() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const analytics = useAnalytics();

  const [currentStep, setCurrentStep] = useState<OnboardingStep>('consent');
  const [grade, setGrade] = useState<number | null>(null);
  const [avatarId, setAvatarId] = useState<string | null>(null);
  const [displayName, setDisplayName] = useState<string>('');
  const [adultConsent, setAdultConsent] = useState(false);
  const [privacyAccepted, setPrivacyAccepted] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [referralCode, setReferralCode] = useState<string | null>(null);
  const [legalModalType, setLegalModalType] = useState<'privacy' | 'terms' | null>(null);
  const hasInitialized = useRef(false);

  // Инициализация при монтировании
  useEffect(() => {
    if (hasInitialized.current) return;
    hasInitialized.current = true;

    const initOnboarding = async () => {
      try {
        // Загружаем данные от VK
        const platformBridge = new PlatformBridge();
        await platformBridge.init();
        const user = await platformBridge.getUser();

        // Автоматически заполняем имя из VK
        const vkDisplayName = user.firstName || 'Ученик';
        setDisplayName(vkDisplayName);

        console.log('[Onboarding] VK user data loaded:', { firstName: user.firstName });

        // Извлекаем реферальный код из URL
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

        analytics.trackEvent('onboarding_opened', {});
      } catch (error) {
        console.error('[Onboarding] Failed to initialize:', error);
        setDisplayName('Ученик');
      }
    };

    initOnboarding();
  }, [searchParams, analytics]);

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

      // Важно: устанавливаем профиль в глобальное хранилище аутентификации
      // Это позволяет всем последующим запросам использовать этот профиль
      import('@/lib/auth').then(({ setCurrentChildProfileId }) => {
        setCurrentChildProfileId(childProfileId);
      });

      await onboardingAPI.saveConsent({
        parentUserId,
        privacyPolicyVersion: '1.0',
        termsVersion: '1.0',
        adultConsent: adultConsent!,
      });

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

  const canProceedConsent = adultConsent && privacyAccepted;
  const canFinish = displayName.trim().length > 0 && avatarId !== null && grade !== null;

  const handleLegalLinkClick = (type: 'privacy' | 'terms') => {
    setLegalModalType(type);
  };

  // Consent Screen
  if (currentStep === 'consent') {
    return (
      <>
        <LegalDocumentModal
          type={legalModalType}
          isOpen={legalModalType !== null}
          onClose={() => setLegalModalType(null)}
        />

        <div className="flex flex-col min-h-screen px-6 py-8 bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF]">
          <div className="flex-1 flex flex-col items-center justify-center text-center">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: 'spring' }}
              className="w-20 h-20 bg-[#6C5CE7]/10 rounded-full flex items-center justify-center mb-6"
            >
              <Shield size={36} className="text-[#6C5CE7]" />
            </motion.div>

            <h1 className="text-[32px] font-bold text-[#6C5CE7] mb-2">Добро пожаловать!</h1>
            <p className="text-[#636e72] text-[14px] mb-8">
              Перед началом работы, пожалуйста, подтвердите согласие
            </p>

            <div className="w-full flex flex-col gap-4 text-left">
              <button
                onClick={() => setAdultConsent(!adultConsent)}
                className="flex items-start gap-3 bg-white rounded-2xl p-4 shadow-sm active:scale-[0.98] transition-transform"
              >
                <div
                  className={`w-6 h-6 mt-0.5 rounded-md flex items-center justify-center flex-shrink-0 transition-all ${
                    adultConsent ? 'bg-[#6C5CE7]' : 'bg-white border-2 border-[#DFE6E9]'
                  }`}
                >
                  {adultConsent && <Check size={14} className="text-white" />}
                </div>
                <span className="text-[14px] text-[#2D3436]">
                  Я подтверждаю, что являюсь взрослым (родителем или опекуном) и даю согласие на
                  использование сервиса ребёнком
                </span>
              </button>

              <button
                onClick={() => setPrivacyAccepted(!privacyAccepted)}
                className="flex items-start gap-3 bg-white rounded-2xl p-4 shadow-sm active:scale-[0.98] transition-transform"
              >
                <div
                  className={`w-6 h-6 mt-0.5 rounded-md flex items-center justify-center flex-shrink-0 transition-all ${
                    privacyAccepted ? 'bg-[#6C5CE7]' : 'bg-white border-2 border-[#DFE6E9]'
                  }`}
                >
                  {privacyAccepted && <Check size={14} className="text-white" />}
                </div>
                <span className="text-[13px] text-[#2D3436] leading-relaxed">
                  Я согласен с{' '}
                  <span
                    onClick={(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                      analytics.trackEvent('privacy_policy_opened', {});
                      handleLegalLinkClick('privacy');
                    }}
                    className="text-[#6C5CE7] underline cursor-pointer"
                  >
                    Политикой конфиденциальности
                  </span>{' '}
                  и{' '}
                  <span
                    onClick={(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                      analytics.trackEvent('terms_opened', {});
                      handleLegalLinkClick('terms');
                    }}
                    className="text-[#6C5CE7] underline cursor-pointer"
                  >
                    Условиями использования
                  </span>
                </span>
              </button>
            </div>
          </div>

          <button
            disabled={!canProceedConsent}
            onClick={() => {
              analytics.trackEvent('adult_consent_checked', {});
              analytics.trackEvent('privacy_policy_accepted', {});
              setCurrentStep('profile');
            }}
            className={`w-full py-4 rounded-2xl transition-all text-white flex items-center justify-center gap-2 ${
              canProceedConsent
                ? 'bg-[#6C5CE7] shadow-lg shadow-[#6C5CE7]/30 active:scale-[0.98]'
                : 'bg-[#B2BEC3] cursor-not-allowed'
            }`}
          >
            Продолжить
            <ArrowRight size={18} />
          </button>
        </div>
      </>
    );
  }

  // Profile Screen
  return (
    <div className="flex flex-col min-h-screen px-6 py-8 bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF]">
      <div className="text-center mb-6">
        <h1 className="text-[32px] font-bold text-[#6C5CE7] mb-1">Создание профиля</h1>
        <p className="text-[#636e72] text-[14px]">Расскажи о себе!</p>
      </div>

      {/* Avatar selection */}
      <div className="mb-5">
        <label className="text-[14px] text-[#636e72] mb-2 block">Выбери аватар</label>
        <div className="grid grid-cols-4 gap-3">
          {avatars.map((a) => (
            <button
              key={a.id}
              onClick={() => {
                setAvatarId(a.id);
                analytics.trackEvent('avatar_selected', { avatar_id: a.id });
              }}
              className={`w-full aspect-square rounded-2xl flex items-center justify-center text-[32px] transition-all ${
                avatarId === a.id
                  ? 'bg-[#6C5CE7] shadow-lg scale-105 ring-3 ring-[#6C5CE7]/30'
                  : 'bg-white shadow-sm active:scale-95'
              }`}
            >
              {a.emoji}
            </button>
          ))}
        </div>
      </div>

      {/* Name input */}
      <div className="mb-5">
        <label className="text-[14px] text-[#636e72] mb-1.5 block">Как тебя зовут?</label>
        <input
          type="text"
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          placeholder="Введи имя"
          className="w-full bg-white rounded-xl px-4 py-3 border border-[#DFE6E9] focus:ring-2 focus:ring-[#6C5CE7]/30 focus:border-[#6C5CE7] outline-none transition-all text-[#2D3436]"
        />
      </div>

      {/* Grade selector */}
      <div className="mb-5">
        <label className="text-[14px] text-[#636e72] mb-1.5 block">В каком классе учишься?</label>
        <div className="flex gap-2 flex-wrap">
          {grades.map((g) => (
            <button
              key={g}
              onClick={() => {
                setGrade(g);
                analytics.trackEvent('grade_selected', { grade: g });
              }}
              className={`px-4 py-2 rounded-xl transition-all text-[14px] ${
                grade === g
                  ? 'bg-[#6C5CE7] text-white shadow'
                  : 'bg-white text-[#2D3436] border border-[#DFE6E9] active:scale-95'
              }`}
            >
              {g} класс
            </button>
          ))}
        </div>
      </div>

      <div className="mt-auto">
        <button
          disabled={!canFinish || isSubmitting}
          onClick={handleComplete}
          className={`w-full py-4 rounded-2xl transition-all text-white ${
            canFinish && !isSubmitting
              ? 'bg-[#6C5CE7] shadow-lg shadow-[#6C5CE7]/30 active:scale-[0.98]'
              : 'bg-[#B2BEC3] cursor-not-allowed'
          }`}
        >
          {isSubmitting ? 'Загрузка...' : 'Начать!'}
        </button>
      </div>
    </div>
  );
}
