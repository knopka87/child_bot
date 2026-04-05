// src/pages/Villain/hooks/useVillainData.ts
import { useEffect, useState } from 'react';
import { villainAPI } from '@/api/villain';
import type { Villain, VillainBattle } from '@/types/villain';
import { useProfileStore } from '@/stores/profileStore';

export function useVillainData() {
  const [villain, setVillain] = useState<Villain | null>(null);
  const [battle, setBattle] = useState<VillainBattle | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const profile = useProfileStore((state) => state.profile);

  const fetchData = async () => {
    if (!profile?.child_profile_id) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const villainData = await villainAPI.getActiveVillain();
      setVillain(villainData);

      if (villainData) {
        const battleData = await villainAPI.getBattleDetails(villainData.id);
        setBattle(battleData);
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
  }, [profile?.child_profile_id]);

  return {
    villain,
    battle,
    isLoading,
    error,
    refetch: fetchData,
  };
}
