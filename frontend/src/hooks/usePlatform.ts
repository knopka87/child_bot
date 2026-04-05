// src/hooks/usePlatform.ts
import { useState, useEffect } from 'react';
import type { PlatformInfo } from '@/lib/platform/platform-detection';
import { detectPlatform } from '@/lib/platform/platform-detection';

export function usePlatform() {
  const [platformInfo, setPlatformInfo] = useState<PlatformInfo | null>(null);

  useEffect(() => {
    detectPlatform().then(setPlatformInfo);
  }, []);

  return {
    platform: platformInfo?.platform || 'android',
    isIOS: platformInfo?.isIOS || false,
    isAndroid: platformInfo?.isAndroid || false,
    isDesktop: platformInfo?.isDesktop || false,
    requestPhotoAccess: async () => {
      // Stub for photo access request
      return true;
    },
  };
}
