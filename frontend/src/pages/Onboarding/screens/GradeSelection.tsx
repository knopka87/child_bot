// src/pages/Onboarding/screens/GradeSelection.tsx
import { Div, Title, Text, SimpleCell } from '@vkontakte/vkui';
import { Icon24CheckCircleOn } from '@vkontakte/icons';

interface GradeSelectionProps {
  selectedGrade: number | null;
  onSelect: (grade: number) => void;
}

export function GradeSelection({ selectedGrade, onSelect }: GradeSelectionProps) {
  const grades = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11];

  return (
    <Div>
      <Div style={{ marginBottom: '24px' }}>
        <Title level="2" weight="2" style={{ marginBottom: '8px' }}>
          В каком ты классе?
        </Title>
        <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
          Выбери свой класс
        </Text>
      </Div>

      <Div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '12px' }}>
        {grades.map((grade) => (
          <SimpleCell
            key={grade}
            onClick={() => onSelect(grade)}
            after={selectedGrade === grade ? <Icon24CheckCircleOn /> : null}
            style={{
              cursor: 'pointer',
              border: '1px solid var(--vkui--color_separator_primary)',
              borderRadius: '8px',
              background:
                selectedGrade === grade
                  ? 'var(--vkui--color_background_accent)'
                  : 'var(--vkui--color_background_content)',
            }}
          >
            <Title level="3" weight="2">
              {grade}
            </Title>
          </SimpleCell>
        ))}
      </Div>
    </Div>
  );
}
