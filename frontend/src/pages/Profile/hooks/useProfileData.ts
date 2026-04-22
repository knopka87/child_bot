// src/pages/Profile/hooks/useProfileData.ts
import { useEffect, useState } from 'react';
import { profileAPI } from '@/api/profile';
import type { ProfileData } from '@/types/profile';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';

export function useProfileData() {
  const [data, setData] = useState<ProfileData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [childProfileId, setChildProfileId] = useState<string | null>(null);

  // Загружаем child_profile_id из storage
  useEffect(() => {
    const loadProfileId = async () => {
      try {
        const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
        console.log('[useProfileData] Loaded child_profile_id from storage:', profileId);
        setChildProfileId(profileId);
      } catch (err) {
        console.error('[useProfileData] Failed to load profile ID from storage:', err);
        setError(err as Error);
        setIsLoading(false);
      }
    };

    loadProfileId();
  }, []);

  const fetchData = async () => {
    if (!childProfileId) {
      console.warn('[useProfileData] No child_profile_id, skipping fetch');
      setIsLoading(false);
      return;
    }

    console.log('[useProfileData] Starting fetch for profile:', childProfileId);
    setIsLoading(true);
    setError(null);

    try {
      const profileData = await profileAPI.getProfile(childProfileId);
      console.log('[useProfileData] Data fetched successfully:', profileData);
      setData(profileData);
    } catch (err) {
      setError(err as Error);
      console.error('[useProfileData] Failed to fetch profile:', err);
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
