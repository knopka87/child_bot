// src/pages/Achievements/hooks/useAchievements.ts
import { useEffect, useState } from 'react';
import { achievementsAPI } from '@/api/achievements';
import type { Achievement, AchievementsStats } from '@/types/achievements';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';

export function useAchievements() {
  const [achievements, setAchievements] = useState<Achievement[]>([]);
  const [stats, setStats] = useState<AchievementsStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [childProfileId, setChildProfileId] = useState<string | null>(null);

  // Загружаем child_profile_id из storage
  useEffect(() => {
    const loadProfileId = async () => {
      try {
        const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
        console.log('[useAchievements] Loaded child_profile_id from storage:', profileId);
        setChildProfileId(profileId);
      } catch (err) {
        console.error('[useAchievements] Failed to load profile ID from storage:', err);
        setError(err as Error);
        setIsLoading(false);
      }
    };

    loadProfileId();
  }, []);

  const fetchAchievements = async () => {
    if (!childProfileId) {
      console.warn('[useAchievements] No child_profile_id, skipping fetch');
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      console.log('[useAchievements] Fetching achievements from backend');

      const [achievementsData, statsData] = await Promise.all([
        achievementsAPI.getAchievements(),
        achievementsAPI.getAchievementsStats(),
      ]);

      console.log('[useAchievements] Achievements loaded:', achievementsData);
      console.log('[useAchievements] Stats loaded:', statsData);

      setAchievements(achievementsData);
      setStats({
        totalAchievements: statsData.totalCount || statsData.total_count || 0,
        unlockedAchievements: statsData.unlockedCount || statsData.unlocked_count || 0,
        totalCoinsEarned: 0, // TODO: добавить на backend
      });
    } catch (err) {
      setError(err as Error);
      console.error('[useAchievements] Failed to fetch achievements:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (childProfileId) {
      fetchAchievements();
    }
  }, [childProfileId]);

  return {
    achievements,
    stats,
    isLoading,
    error,
    refetch: fetchAchievements,
  };
}
