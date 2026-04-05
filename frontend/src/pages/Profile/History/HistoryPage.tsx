// src/pages/Profile/History/HistoryPage.tsx
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, BookOpen, ClipboardCheck, CheckCircle, XCircle, Clock, Lightbulb, Loader } from 'lucide-react';
import { motion } from 'framer-motion';
import { BottomNav } from '@/components/layout/BottomNav';
import { Spinner } from '@/components/ui/Spinner';
import { useHistory } from './hooks/useHistory';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { ROUTES } from '@/config/routes';
import { HistoryDetailModal } from './components/HistoryDetailModal';
import { MOCK_TASK_TITLES } from '@/api/mockHistoryData';
import type { HistoryAttempt, HistoryFilters } from '@/types/profile';
import styles from './HistoryPage.module.css';

type FilterMode = 'all' | 'help' | 'check';
type FilterStatus = 'all' | 'success' | 'error' | 'in_progress';

export function HistoryPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const [childProfileId, setChildProfileId] = useState<string | null>(null);
  const [filterMode, setFilterMode] = useState<FilterMode>('all');
  const [filterStatus, setFilterStatus] = useState<FilterStatus>('all');
  const [selectedAttempt, setSelectedAttempt] = useState<HistoryAttempt | null>(null);

  // Загружаем child_profile_id
  useEffect(() => {
    const loadProfileId = async () => {
      const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      console.log('[HistoryPage] Loaded child_profile_id:', profileId);
      setChildProfileId(profileId);
    };

    loadProfileId();
  }, []);

  // Формируем фильтры для API
  const filters: HistoryFilters | undefined = {
    mode: filterMode === 'all' ? undefined : filterMode,
    status: filterStatus === 'all' ? undefined : filterStatus,
  };

  const { data, isLoading, error, refetch } = useHistory(childProfileId, filters);

  // Трекаем открытие страницы
  useEffect(() => {
    if (childProfileId) {
      analytics.trackEvent('history_opened', {
        child_profile_id: childProfileId,
      });
    }
  }, [childProfileId, analytics]);

  const handleBack = () => {
    if (childProfileId) {
      analytics.trackEvent('history_back_clicked', {
        child_profile_id: childProfileId,
      });
    }
    navigate(ROUTES.PROFILE);
  };

  const handleCardClick = (attempt: HistoryAttempt) => {
    if (childProfileId) {
      analytics.trackEvent('history_item_clicked', {
        child_profile_id: childProfileId,
        attempt_id: attempt.id,
        mode: attempt.mode,
        status: attempt.status,
      });
    }
    setSelectedAttempt(attempt);
  };

  const handleFilterModeChange = (mode: FilterMode) => {
    setFilterMode(mode);
    if (childProfileId) {
      analytics.trackEvent('history_filter_changed', {
        child_profile_id: childProfileId,
        filter_type: 'mode',
        filter_value: mode,
      });
    }
  };

  const handleFilterStatusChange = (status: FilterStatus) => {
    setFilterStatus(status);
    if (childProfileId) {
      analytics.trackEvent('history_filter_changed', {
        child_profile_id: childProfileId,
        filter_type: 'status',
        filter_value: status,
      });
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const day = date.getDate();
    const month = date.toLocaleDateString('ru-RU', { month: 'long' });
    const year = date.getFullYear();
    return `${day} ${month} ${year}`;
  };

  const getStatusConfig = (attempt: HistoryAttempt) => {
    // Если использованы подсказки - приоритет этому статусу
    if (attempt.hintsUsed && attempt.hintsUsed > 0) {
      return {
        label: 'Использованы подсказки',
        icon: <Lightbulb size={16} />,
        className: styles.hints,
      };
    }

    switch (attempt.status) {
      case 'success':
        return {
          label: 'Решено верно',
          icon: <CheckCircle size={16} />,
          className: styles.success,
        };
      case 'error':
        return {
          label: 'Есть ошибки',
          icon: <XCircle size={16} />,
          className: styles.error,
        };
      case 'in_progress':
        return {
          label: 'В обработке',
          icon: <Loader size={16} />,
          className: styles.processing,
        };
      default:
        return {
          label: 'Незакончено',
          icon: <Clock size={16} />,
          className: styles.unfinished,
        };
    }
  };

  const getTaskTitle = (attempt: HistoryAttempt): string => {
    // Используем mock названия, в проде будет приходить с бэкенда
    return MOCK_TASK_TITLES[attempt.id] ||
           `${attempt.mode === 'help' ? 'Помощь' : 'Проверка'}${
             attempt.scenarioType === 'two_photo' ? ' — 2 фото' : ''
           }`;
  };

  if (isLoading && !childProfileId) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <Spinner size="lg" />
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <button onClick={handleBack} className={styles.backButton}>
          <ArrowLeft size={20} />
          <span>Профиль</span>
        </button>
      </div>

      <h1 className={styles.title}>История</h1>

      {/* Filters */}
      <div className={styles.filtersContainer}>
        {/* Фильтр по режиму */}
        <div className={styles.filterGroup}>
          <span className={styles.filterLabel}>Режим:</span>
          <div className={styles.filters}>
            <button
              onClick={() => handleFilterModeChange('all')}
              className={`${styles.filterButton} ${filterMode === 'all' ? styles.active : ''}`}
            >
              Все
            </button>
            <button
              onClick={() => handleFilterModeChange('help')}
              className={`${styles.filterButton} ${filterMode === 'help' ? styles.active : ''}`}
            >
              Помощь
            </button>
            <button
              onClick={() => handleFilterModeChange('check')}
              className={`${styles.filterButton} ${filterMode === 'check' ? styles.active : ''}`}
            >
              Проверка
            </button>
          </div>
        </div>

        {/* Фильтр по статусу */}
        <div className={styles.filterGroup}>
          <span className={styles.filterLabel}>Результат:</span>
          <div className={styles.filters}>
            <button
              onClick={() => handleFilterStatusChange('all')}
              className={`${styles.filterButton} ${filterStatus === 'all' ? styles.active : ''}`}
            >
              Все
            </button>
            <button
              onClick={() => handleFilterStatusChange('success')}
              className={`${styles.filterButton} ${filterStatus === 'success' ? styles.active : ''}`}
            >
              Верно
            </button>
            <button
              onClick={() => handleFilterStatusChange('error')}
              className={`${styles.filterButton} ${filterStatus === 'error' ? styles.active : ''}`}
            >
              Ошибки
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      {isLoading ? (
        <div className={styles.loadingContainer}>
          <Spinner size="lg" />
        </div>
      ) : error ? (
        <div className={styles.errorContainer}>
          <p className={styles.errorText}>Не удалось загрузить историю</p>
          <button onClick={refetch} className={styles.retryButton}>
            Повторить
          </button>
        </div>
      ) : data.length === 0 ? (
        <div className={styles.emptyContainer}>
          <p className={styles.emptyText}>История пуста</p>
          <p className={styles.emptyText}>Попробуйте решить первую задачу!</p>
        </div>
      ) : (
        <div className={styles.historyList}>
          {data.map((attempt, index) => {
            const statusConfig = getStatusConfig(attempt);
            const ModeIcon = attempt.mode === 'help' ? BookOpen : ClipboardCheck;
            const taskTitle = getTaskTitle(attempt);

            return (
              <motion.button
                key={attempt.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.05 }}
                onClick={() => handleCardClick(attempt)}
                className={styles.historyCard}
              >
                <div className={styles.cardHeader}>
                  <div className={styles.modeIcon}>
                    <ModeIcon size={20} />
                  </div>
                  <div className={styles.cardContent}>
                    <div className={styles.cardTitleRow}>
                      <span className={styles.cardTitle}>{taskTitle}</span>
                      <span className={styles.cardDate}>{formatDate(attempt.createdAt)}</span>
                    </div>
                    <div className={`${styles.statusBadge} ${statusConfig.className}`}>
                      {statusConfig.icon}
                      <span>{statusConfig.label}</span>
                    </div>
                  </div>
                </div>
              </motion.button>
            );
          })}
        </div>
      )}

      {/* Detail Modal */}
      {selectedAttempt && (
        <HistoryDetailModal
          attempt={selectedAttempt}
          isOpen={!!selectedAttempt}
          onClose={() => setSelectedAttempt(null)}
          childProfileId={childProfileId}
        />
      )}

      <BottomNav />
    </div>
  );
}
