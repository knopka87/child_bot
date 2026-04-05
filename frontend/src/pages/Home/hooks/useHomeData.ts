// src/pages/Home/hooks/useHomeData.ts
import { useEffect, useState } from 'react';
import { homeAPI } from '@/api/home';
import type { HomeData } from '@/types/home';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';

export function useHomeData() {
  const [data, setData] = useState<HomeData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [childProfileId, setChildProfileId] = useState<string | null>(null);

  // Загружаем child_profile_id из storage
  useEffect(() => {
    const loadProfileId = async () => {
      try {
        const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
        console.log('[useHomeData] Loaded child_profile_id from storage:', profileId);
        setChildProfileId(profileId);
      } catch (err) {
        console.error('[useHomeData] Failed to load profile ID from storage:', err);
        setError(err as Error);
        setIsLoading(false);
      }
    };

    loadProfileId();
  }, []);

  const fetchData = async () => {
    if (!childProfileId) {
      console.warn('[useHomeData] No child_profile_id, skipping fetch');
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      console.log('[useHomeData] Fetching home data for profile:', childProfileId);
      const homeData = await homeAPI.getHomeData(childProfileId);
      console.log('[useHomeData] Home data loaded successfully:', homeData);

      // ВРЕМЕННО: Добавляем мок-данные для villain, если их нет
      if (!homeData.villain) {
        console.log('[useHomeData] Adding mock villain data for development');
        homeData.villain = {
          id: 'villain-1',
          name: 'Кракозябра',
          imageUrl: '/images/villain.png',
          healthPercent: 66, // 2 из 3 полосок
          isDefeated: false,
        };
      }

      setData(homeData);
    } catch (err) {
      setError(err as Error);
      console.error('[useHomeData] Failed to fetch home data:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (childProfileId) {
      fetchData();
    }
  }, [childProfileId]);

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  };
}
