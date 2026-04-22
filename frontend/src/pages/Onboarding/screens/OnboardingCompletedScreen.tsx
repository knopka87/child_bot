// src/pages/Onboarding/screens/OnboardingCompletedScreen.tsx
import { Div, Title, Text } from '@vkontakte/vkui';
import { PrimaryButton } from '@/components/ui/Button';

interface OnboardingCompletedScreenProps {
  displayName: string;
  onComplete: () => void;
}

export function OnboardingCompletedScreen({
  displayName,
  onComplete,
}: OnboardingCompletedScreenProps) {
  return (
    <Div>
      <Div style={{ textAlign: 'center', marginBottom: '32px' }}>
        <div style={{ fontSize: '64px', marginBottom: '16px' }}>🎉</div>
        <Title level="1" weight="1" style={{ marginBottom: '16px' }}>
          Отлично, {displayName}!
        </Title>
        <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
          Теперь ты готов начать учиться!
        </Text>
      </Div>

      <PrimaryButton onClick={onComplete}>Начать</PrimaryButton>
    </Div>
  );
}
