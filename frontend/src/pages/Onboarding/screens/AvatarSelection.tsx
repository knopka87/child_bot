// src/pages/Onboarding/screens/AvatarSelection.tsx
import { Div, Title, Text } from '@vkontakte/vkui';
import { Icon24CheckCircleOn } from '@vkontakte/icons';
import { useState, useEffect } from 'react';
import { onboardingAPI } from '@/api/onboarding';
import type { Avatar } from '@/types/onboarding';

interface AvatarSelectionProps {
  selectedAvatarId: string | null;
  onSelect: (avatarId: string) => void;
}

// Mock avatars для разработки
const MOCK_AVATARS: Avatar[] = [
  { id: '1', imageUrl: '🐱', name: 'Кот', isPremium: false },
  { id: '2', imageUrl: '🐶', name: 'Пёс', isPremium: false },
  { id: '3', imageUrl: '🐼', name: 'Панда', isPremium: false },
  { id: '4', imageUrl: '🦊', name: 'Лиса', isPremium: false },
  { id: '5', imageUrl: '🐻', name: 'Медведь', isPremium: false },
  { id: '6', imageUrl: '🦁', name: 'Лев', isPremium: false },
  { id: '7', imageUrl: '🐯', name: 'Тигр', isPremium: true },
  { id: '8', imageUrl: '🦄', name: 'Единорог', isPremium: true },
];

export function AvatarSelection({ selectedAvatarId, onSelect }: AvatarSelectionProps) {
  console.log('=== [AvatarSelection] Component START ===');
  console.log('[AvatarSelection] Component rendered!', { selectedAvatarId });

  const [avatars, setAvatars] = useState<Avatar[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  console.log('[AvatarSelection] Before useEffect', { isLoading, avatarsCount: avatars.length });

  useEffect(() => {
    console.log('=== [AvatarSelection] useEffect START ===');
    console.log('[AvatarSelection] useEffect triggered - fetching avatars from API...');

    onboardingAPI
      .getAvatars()
      .then((data) => {
        console.log('[AvatarSelection] SUCCESS - Avatars loaded:', data);
        setAvatars(data);
      })
      .catch((error) => {
        // Fallback к mock данным
        console.error('[AvatarSelection] ERROR - Failed to load avatars, using mock data:', error);
        setAvatars(MOCK_AVATARS);
      })
      .finally(() => {
        console.log('[AvatarSelection] FINALLY - Loading complete, setting isLoading to false');
        setIsLoading(false);
      });

    console.log('[AvatarSelection] useEffect - API call initiated');
  }, []);

  if (isLoading) {
    return (
      <Div>
        <Text>Загружаю аватары из API... (v2)</Text>
        <Text style={{ marginTop: '8px', fontSize: '12px', color: '#999' }}>
          isLoading: {String(isLoading)}, avatars count: {avatars.length}
        </Text>
      </Div>
    );
  }

  if (avatars.length === 0) {
    return (
      <Div>
        <Text>Ошибка загрузки аватаров</Text>
      </Div>
    );
  }

  return (
    <Div>
      <Div style={{ marginBottom: '24px' }}>
        <Title level="2" weight="2" style={{ marginBottom: '8px' }}>
          Выбери своего маскота
        </Title>
        <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
          Он будет помогать тебе с домашкой
        </Text>
      </Div>

      <div
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(4, 1fr)',
          gap: '12px',
          padding: '0',
        }}
      >
        {avatars.map((avatar) => (
          <div
            key={avatar.id}
            onClick={() => onSelect(avatar.id)}
            style={{
              cursor: 'pointer',
              border: `2px solid ${
                selectedAvatarId === avatar.id
                  ? 'var(--vkui--color_accent_blue)'
                  : 'var(--vkui--color_separator_primary)'
              }`,
              borderRadius: '12px',
              background:
                selectedAvatarId === avatar.id
                  ? 'var(--vkui--color_background_accent_tint)'
                  : 'var(--vkui--color_background_content)',
              textAlign: 'center',
              padding: '12px 8px',
              position: 'relative',
              transition: 'all 0.2s ease',
            }}
          >
            <div style={{ fontSize: '48px', lineHeight: '1', marginBottom: '8px' }}>
              {avatar.imageUrl}
            </div>
            <Text style={{ fontSize: '11px', fontWeight: '500' }}>{avatar.name}</Text>
            {avatar.isPremium && (
              <Text style={{ fontSize: '10px', color: 'var(--vkui--color_text_accent)' }}>
                ⭐
              </Text>
            )}
            {selectedAvatarId === avatar.id && (
              <div
                style={{
                  position: 'absolute',
                  top: '4px',
                  right: '4px',
                  color: 'var(--vkui--color_accent_blue)',
                }}
              >
                <Icon24CheckCircleOn width={20} height={20} />
              </div>
            )}
          </div>
        ))}
      </div>
    </Div>
  );
}
