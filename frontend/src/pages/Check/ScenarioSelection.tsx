// src/pages/Check/ScenarioSelection.tsx
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Panel, PanelHeader, PanelHeaderBack, Group, Div, Title, Text } from '@vkontakte/vkui';
import { Card } from '@/components/ui/Card';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';
import type { CheckScenario } from '@/types/check';
import styles from './ScenarioSelection.module.css';

export default function ScenarioSelection() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  useEffect(() => {
    analytics.trackEvent('check_scenario_selection_opened', {
      child_profile_id: profile?.child_profile_id,
    });
  }, []);

  const handleScenarioSelect = (scenario: CheckScenario) => {
    analytics.trackEvent('check_scenario_selected', {
      child_profile_id: profile?.child_profile_id,
      scenario_type: scenario,
    });

    navigate('/check/upload', {
      state: { scenario },
    });
  };

  return (
    <Panel id="scenario-selection">
      <PanelHeader before={<PanelHeaderBack onClick={() => navigate(ROUTES.HOME)} />}>
        Проверка
      </PanelHeader>

      <Group>
        <Div>
          <Title level="1" weight="1" style={{ marginBottom: '8px' }}>
            Проверка ДЗ
          </Title>
          <Text style={{ color: 'var(--vkui--color_text_secondary)', marginBottom: '16px' }}>
            Выбери сценарий проверки
          </Text>

          <Card
            className={styles.scenarioCard}
            variant="bordered"
            onClick={() => handleScenarioSelect('single_photo')}
          >
            <div className={styles.scenarioIcon}>📷</div>
            <Title level="2" weight="2" style={{ marginBottom: '8px' }}>
              Одно фото
            </Title>
            <Text>
              Задание и ответ на одном фото (например, решённая задача в тетради)
            </Text>
          </Card>

          <Card
            className={styles.scenarioCard}
            variant="bordered"
            onClick={() => handleScenarioSelect('two_photo')}
          >
            <div className={styles.scenarioIcon}>📷 📷</div>
            <Title level="2" weight="2" style={{ marginBottom: '8px' }}>
              Два фото
            </Title>
            <Text>
              Отдельно задание из учебника и твой ответ
            </Text>
          </Card>
        </Div>
      </Group>
    </Panel>
  );
}
