// src/pages/Friends/hooks/useReferralData.ts
import { useEffect, useState } from 'react';
import { referralAPI } from '@/api/referral';
import type { ReferralData } from '@/types/referral';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';

export function useReferralData() {
  const [data, setData] = useState<ReferralData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [childProfileId, setChildProfileId] = useState<string | null>(null);

  // Загружаем child_profile_id из storage
  useEffect(() => {
    const loadProfileId = async () => {
      try {
        const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
        console.log('[useReferralData] Loaded child_profile_id from storage:', profileId);
        setChildProfileId(profileId);
      } catch (err) {
        console.error('[useReferralData] Failed to load profile ID from storage:', err);
        setError(err as Error);
        setIsLoading(false);
      }
    };

    loadProfileId();
  }, []);

  const fetchData = async () => {
    if (!childProfileId) {
      console.warn('[useReferralData] No child_profile_id, skipping fetch');
      setIsLoading(false);
      return;
    }

    console.log('[useReferralData] Starting fetch for profile:', childProfileId);
    setIsLoading(true);
    setError(null);

    try {
      const referralData = await referralAPI.getReferralData(childProfileId);
      console.log('[useReferralData] Data fetched successfully:', referralData);
      setData(referralData);
    } catch (err) {
      setError(err as Error);
      console.error('[useReferralData] Failed to fetch referral data:', err);
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
