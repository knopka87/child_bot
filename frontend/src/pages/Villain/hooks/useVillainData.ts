// src/pages/Villain/hooks/useVillainData.ts
import { useEffect, useState } from 'react';
import { villainAPI } from '@/api/villain';
import type { Villain, VillainBattle } from '@/types/villain';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';

export function useVillainData() {
  const [villain, setVillain] = useState<Villain | null>(null);
  const [battle, setBattle] = useState<VillainBattle | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchData = async () => {
    console.log('[useVillainData] fetchData called');

    setIsLoading(true);
    setError(null);

    try {
      // Получаем child_profile_id из storage
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      console.log('[useVillainData] childProfileId from storage:', childProfileId);

      if (!childProfileId) {
        console.log('[useVillainData] No child_profile_id in storage, skipping fetch');
        setIsLoading(false);
        return;
      }

      console.log('[useVillainData] Fetching active villain...');
      const villainData = await villainAPI.getActiveVillain();
      console.log('[useVillainData] Active villain data:', villainData);
      setVillain(villainData);

      if (villainData) {
        console.log('[useVillainData] Fetching battle details for villain:', villainData.id);
        const battleData = await villainAPI.getBattleDetails(villainData.id);
        console.log('[useVillainData] Battle data:', battleData);
        setBattle(battleData);
      } else {
        console.log('[useVillainData] No active villain found');
      }
    } catch (err) {
      setError(err as Error);
      console.error('[useVillainData] Failed to fetch villain data:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  return {
    villain,
    battle,
    isLoading,
    error,
    refetch: fetchData,
  };
}
