// src/pages/Onboarding/screens/WelcomeScreen.tsx
import { Div, Title, Text } from '@vkontakte/vkui';
import { PrimaryButton } from '@/components/ui/Button';

interface WelcomeScreenProps {
  onNext: () => void;
}

export function WelcomeScreen({ onNext }: WelcomeScreenProps) {
  return (
    <Div>
      <Div style={{ textAlign: 'center', marginBottom: '24px' }}>
        <Title level="1" weight="1" style={{ marginBottom: '16px' }}>
          👋 Привет!
        </Title>
        <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
          Давай настроим твой профиль
        </Text>
      </Div>

      <PrimaryButton onClick={onNext}>Начать</PrimaryButton>
    </Div>
  );
}
