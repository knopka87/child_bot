// src/pages/Onboarding/screens/DisplayNameInput.tsx
import { Div, Title, Text, FormItem, Input } from '@vkontakte/vkui';

interface DisplayNameInputProps {
  displayName: string;
  onChange: (displayName: string) => void;
}

export function DisplayNameInput({ displayName, onChange }: DisplayNameInputProps) {
  return (
    <Div>
      <Div style={{ marginBottom: '24px' }}>
        <Title level="2" weight="2" style={{ marginBottom: '8px' }}>
          Как тебя зовут?
        </Title>
        <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
          Введи своё имя
        </Text>
      </Div>

      <FormItem top="Имя">
        <Input
          type="text"
          placeholder="Например, Иван"
          value={displayName}
          onChange={(e) => onChange(e.target.value)}
        />
      </FormItem>

      {displayName.trim().length > 0 && displayName.trim().length < 2 && (
        <Text style={{ color: 'var(--vkui--color_text_negative)', fontSize: '14px' }}>
          Имя должно содержать минимум 2 символа
        </Text>
      )}
    </Div>
  );
}
