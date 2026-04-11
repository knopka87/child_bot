// src/hooks/useNewAchievements.ts
import { useEffect, useState, useCallback } from 'react';
import { achievementsAPI } from '@/api/achievements';

export function useNewAchievements() {
  const [hasNew, setHasNew] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchNewStatus = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const result = await achievementsAPI.hasNewAchievements();
      setHasNew(result.has_new);
    } catch (err) {
      setError(err as Error);
      console.error('[useNewAchievements] Failed to check new achievements:', err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchNewStatus();

    // Refresh every 30 seconds
    const interval = setInterval(fetchNewStatus, 30000);

    return () => clearInterval(interval);
  }, [fetchNewStatus]);

  const markAsViewed = useCallback(async () => {
    try {
      await achievementsAPI.markAchievementsViewed();
      setHasNew(false); // Update local state immediately
    } catch (err) {
      console.error('[useNewAchievements] Failed to mark as viewed:', err);
    }
  }, []);

  return {
    hasNew,
    isLoading,
    error,
    refetch: fetchNewStatus,
    markAsViewed,
  };
}
