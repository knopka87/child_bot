// src/pages/Onboarding/screens/EmailInput.tsx
import { Div, Title, Text, FormItem, Input } from '@vkontakte/vkui';

interface EmailInputProps {
  email: string;
  onChange: (email: string) => void;
}

export function EmailInput({ email, onChange }: EmailInputProps) {
  const isValidEmail = email.includes('@') && email.includes('.');

  return (
    <Div>
      <Div style={{ marginBottom: '24px' }}>
        <Title level="2" weight="2" style={{ marginBottom: '8px' }}>
          Email для связи
        </Title>
        <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
          На этот email мы отправим отчёты о прогрессе
        </Text>
      </Div>

      <FormItem top="Email" status={email && !isValidEmail ? 'error' : 'default'}>
        <Input
          type="email"
          placeholder="example@mail.ru"
          value={email}
          onChange={(e) => onChange(e.target.value)}
        />
      </FormItem>

      {email && !isValidEmail && (
        <Text style={{ color: 'var(--vkui--color_text_negative)', fontSize: '14px' }}>
          Введите корректный email
        </Text>
      )}
    </Div>
  );
}
