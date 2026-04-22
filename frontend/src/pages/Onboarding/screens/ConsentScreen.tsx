// src/pages/Onboarding/screens/ConsentScreen.tsx
import { Div, Title, Text, Checkbox, Link } from '@vkontakte/vkui';

interface ConsentScreenProps {
  adultConsent: boolean;
  privacyAccepted: boolean;
  termsAccepted: boolean;
  onAdultConsentChange: (value: boolean) => void;
  onPrivacyChange: (value: boolean) => void;
  onTermsChange: (value: boolean) => void;
}

export function ConsentScreen({
  adultConsent,
  privacyAccepted,
  termsAccepted,
  onAdultConsentChange,
  onPrivacyChange,
  onTermsChange,
}: ConsentScreenProps) {
  return (
    <Div>
      <Div style={{ marginBottom: '24px' }}>
        <Title level="2" weight="2" style={{ marginBottom: '8px' }}>
          Согласия
        </Title>
        <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
          Необходимо согласие взрослого
        </Text>
      </Div>

      <Div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <Checkbox checked={adultConsent} onChange={(e) => onAdultConsentChange(e.target.checked)}>
          <Text>Я являюсь родителем или законным представителем</Text>
        </Checkbox>

        <Checkbox checked={privacyAccepted} onChange={(e) => onPrivacyChange(e.target.checked)}>
          <Text>
            Я принимаю{' '}
            <Link href="/privacy" target="_blank">
              политику конфиденциальности
            </Link>
          </Text>
        </Checkbox>

        <Checkbox checked={termsAccepted} onChange={(e) => onTermsChange(e.target.checked)}>
          <Text>
            Я принимаю{' '}
            <Link href="/terms" target="_blank">
              условия использования
            </Link>
          </Text>
        </Checkbox>
      </Div>
    </Div>
  );
}
