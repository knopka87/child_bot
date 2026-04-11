// src/pages/Profile/History/components/HistoryDetailModal.tsx
import { motion, AnimatePresence } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import { RefreshCw, Pencil, CheckCircle, XCircle, Clock, Lightbulb } from 'lucide-react';
import { useAnalytics } from '@/hooks/useAnalytics';
import { ROUTES } from '@/config/routes';
import type { HistoryAttempt } from '@/types/profile';
import styles from './HistoryDetailModal.module.css';

interface HistoryDetailModalProps {
  attempt: HistoryAttempt;
  isOpen: boolean;
  onClose: () => void;
  childProfileId: string | null;
}

export function HistoryDetailModal({
  attempt,
  isOpen,
  onClose,
  childProfileId,
}: HistoryDetailModalProps) {
  const navigate = useNavigate();
  const analytics = useAnalytics();

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const day = date.getDate();
    const month = date.toLocaleDateString('ru-RU', { month: 'long' });
    return `${day} ${month}`;
  };

  const getStatusConfig = () => {
    // Для режима help используем статус из attempt, а не result
    if (attempt.mode === 'help') {
      if (attempt.status === 'completed') {
        return {
          label: 'Завершено',
          icon: <CheckCircle size={14} />,
          color: '#00B894',
          bg: '#E8FFF8',
        };
      } else if (attempt.status === 'in_progress') {
        return {
          label: 'В процессе',
          icon: <Clock size={14} />,
          color: '#FDCB6E',
          bg: '#FFF9E8',
        };
      }
    }

    // Для режима check используем result
    if (!attempt.result) {
      return {
        label: 'В процессе',
        icon: <Clock size={14} />,
        color: '#FDCB6E',
        bg: '#FFF9E8',
      };
    }

    switch (attempt.result.status) {
      case 'correct':
        return {
          label: 'Решено верно',
          icon: <CheckCircle size={14} />,
          color: '#00B894',
          bg: '#E8FFF8',
        };
      case 'has_errors':
        return {
          label: `Есть ${attempt.result.errorCount || 0} ошибок`,
          icon: <XCircle size={14} />,
          color: '#DC3545',
          bg: '#FFE8E8',
        };
      case 'processing':
        return {
          label: 'В обработке',
          icon: <Clock size={14} />,
          color: '#6C757D',
          bg: '#F8F9FA',
        };
      default:
        return {
          label: 'Неизвестно',
          icon: <Clock size={14} />,
          color: '#6C757D',
          bg: '#F8F9FA',
        };
    }
  };

  const handleContinue = () => {
    if (childProfileId) {
      analytics.trackEvent('unfinished_attempt_continue_clicked', {
        child_profile_id: childProfileId,
        attempt_id: attempt.id,
        mode: attempt.mode,
        status: attempt.status,
      });
    }

    onClose();

    // Navigate to appropriate flow (используем тот же подход, что и на главной)
    const route =
      attempt.mode === 'help' ? ROUTES.HELP_PROCESSING : ROUTES.CHECK_PROCESSING;

    navigate(route, {
      state: { attemptId: attempt.id },
    });
  };

  const handleFixErrors = () => {
    if (childProfileId) {
      analytics.trackEvent('history_fix_errors_clicked', {
        child_profile_id: childProfileId,
        attempt_id: attempt.id,
      });
    }

    onClose();
    navigate(ROUTES.CHECK);
  };

  const statusConfig = getStatusConfig();
  const hasErrors = attempt.result?.status === 'has_errors' && (attempt.result.errorCount || 0) > 0;
  const isInProgress = attempt.status === 'in_progress';

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className={styles.overlay}
          onClick={onClose}
        >
          <motion.div
            initial={{ y: '100%' }}
            animate={{ y: 0 }}
            exit={{ y: '100%' }}
            transition={{ type: 'spring', damping: 30, stiffness: 300 }}
            className={styles.modal}
            onClick={(e) => e.stopPropagation()}
          >
            {/* Drag handle */}
            <div className={styles.dragHandle} />

            {/* Header */}
            <div className={styles.header}>
              <h3 className={styles.title}>
                {attempt.mode === 'help' ? 'Помощь' : 'Проверка'}
                {attempt.scenarioType === 'two_photo' && ' — 2 фото'}
              </h3>
              <p className={styles.date}>{formatDate(attempt.createdAt)}</p>
            </div>

            {/* Status badge */}
            <div
              className={styles.statusBadge}
              style={{
                color: statusConfig.color,
                background: statusConfig.bg,
              }}
            >
              {statusConfig.icon}
              <span>{statusConfig.label}</span>
            </div>

            {/* Hints used */}
            {attempt.hintsUsed && attempt.hintsUsed > 0 && (
              <div className={styles.hintsInfo}>
                <Lightbulb size={16} color="#6C5CE7" />
                <span>Использовано подсказок: {attempt.hintsUsed}</span>
              </div>
            )}

            {/* Images */}
            {attempt.images && attempt.images.length > 0 && (
              <div className={styles.imagesSection}>
                <h4 className={styles.sectionTitle}>Изображения</h4>
                <div className={styles.imageGrid}>
                  {attempt.images.map((image) => (
                    <div key={image.id} className={styles.imageCard}>
                      <img
                        src={image.thumbnailUrl || image.url}
                        alt={`${image.role}`}
                        className={styles.image}
                      />
                      <p className={styles.imageLabel}>
                        {image.role === 'task' && 'Условие'}
                        {image.role === 'answer' && 'Решение'}
                        {image.role === 'single' && 'Фото'}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Result summary */}
            {attempt.result?.summary && (
              <div className={styles.summarySection}>
                <h4 className={styles.sectionTitle}>Результат</h4>
                <p className={styles.summaryText}>{attempt.result.summary}</p>
              </div>
            )}

            {/* Errors */}
            {attempt.result?.feedback && attempt.result.feedback.length > 0 && (
              <div className={styles.errorsSection}>
                <h4 className={styles.sectionTitle}>Найденные ошибки</h4>
                <div className={styles.errorsList}>
                  {attempt.result.feedback.map((error, index) => (
                    <div key={error.id} className={styles.errorCard}>
                      <div className={styles.errorHeader}>
                        <span className={styles.errorNumber}>#{index + 1}</span>
                        {error.stepNumber && (
                          <span className={styles.errorStep}>Шаг {error.stepNumber}</span>
                        )}
                        {error.lineReference && (
                          <span className={styles.errorLine}>{error.lineReference}</span>
                        )}
                      </div>
                      <p className={styles.errorText}>{error.description}</p>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Actions */}
            <div className={styles.actions}>
              {/* Кнопка "Продолжить" только для незавершенных попыток */}
              {isInProgress && (
                <button onClick={handleContinue} className={styles.primaryButton}>
                  <RefreshCw size={16} />
                  Продолжить
                </button>
              )}

              {hasErrors && (
                <button onClick={handleFixErrors} className={styles.primaryButton}>
                  <Pencil size={16} />
                  Исправить и проверить
                </button>
              )}

              <button onClick={onClose} className={styles.primaryButton}>
                Закрыть
              </button>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
