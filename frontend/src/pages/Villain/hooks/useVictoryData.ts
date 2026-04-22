// src/pages/Villain/hooks/useVictoryData.ts
import { useEffect, useState } from 'react';
import { villainAPI } from '@/api/villain';
import type { VillainVictory } from '@/types/villain';
import { useProfileStore } from '@/stores/profileStore';

export function useVictoryData(villainId: string) {
  const [victory, setVictory] = useState<VillainVictory | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const profile = useProfileStore((state) => state.profile);

  useEffect(() => {
    const fetchVictory = async () => {
      if (!profile?.child_profile_id || !villainId) {
        setIsLoading(false);
        return;
      }

      setIsLoading(true);
      setError(null);

      try {
        const victoryData = await villainAPI.getVictoryDetails(villainId);
        setVictory(victoryData);
      } catch (err) {
        setError(err as Error);
        console.error('[useVictoryData] Failed to fetch victory data:', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchVictory();
  }, [profile?.child_profile_id, villainId]);

  return {
    victory,
    isLoading,
    error,
  };
}
