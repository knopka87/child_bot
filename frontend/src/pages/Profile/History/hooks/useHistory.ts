// src/pages/Profile/History/hooks/useHistory.ts
import { useState, useEffect } from 'react';
import { profileAPI } from '@/api/profile';
import { MOCK_HISTORY_DATA } from '@/api/mockHistoryData';
import type { HistoryAttempt, HistoryFilters } from '@/types/profile';

interface UseHistoryResult {
  data: HistoryAttempt[];
  isLoading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
}

// Флаг для переключения между mock и реальными данными
const USE_MOCK_DATA = false;

export function useHistory(
  childProfileId: string | null,
  filters?: HistoryFilters
): UseHistoryResult {
  const [data, setData] = useState<HistoryAttempt[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchHistory = async () => {
    if (!childProfileId && !USE_MOCK_DATA) {
      console.log('[useHistory] No childProfileId, skipping fetch');
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      setError(null);

      let history: HistoryAttempt[];

      if (USE_MOCK_DATA) {
        console.log('[useHistory] Using MOCK data');
        // Имитация задержки API
        await new Promise((resolve) => setTimeout(resolve, 500));
        history = MOCK_HISTORY_DATA;
      } else {
        console.log('[useHistory] Fetching history for:', childProfileId, 'filters:', filters);
        history = await profileAPI.getHistory(childProfileId!, filters);
      }

      // Применяем фильтры к mock данным
      let filteredHistory = [...history];

      if (filters?.mode && filters.mode !== 'all') {
        filteredHistory = filteredHistory.filter((h) => h.mode === filters.mode);
      }

      if (filters?.status && filters.status !== 'all') {
        filteredHistory = filteredHistory.filter((h) => h.status === filters.status);
      }

      console.log('[useHistory] Filtered history:', filteredHistory);
      setData(filteredHistory);
    } catch (err) {
      console.error('[useHistory] Failed to fetch history:', err);
      setError(err instanceof Error ? err : new Error('Failed to fetch history'));
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchHistory();
  }, [childProfileId, filters?.mode, filters?.status, filters?.dateFrom, filters?.dateTo]);

  return {
    data,
    isLoading,
    error,
    refetch: fetchHistory,
  };
}
