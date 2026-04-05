// src/lib/platform/platform-detection.ts
import bridge from './bridge';
import { useState, useEffect } from 'react';

export type Platform = 'android' | 'ios' | 'vkcom';

export type VKPlatform =
  | 'mobile_iphone'
  | 'mobile_ipad'
  | 'mobile_android'
  | 'desktop_web'
  | 'mobile_web';

export interface PlatformInfo {
  platform: Platform; // для VKUI
  vkPlatform: VKPlatform;
  isIOS: boolean;
  isAndroid: boolean;
  isDesktop: boolean;
  isMobile: boolean;
}

let cachedPlatformInfo: PlatformInfo | null = null;

/**
 * Определение платформы через VK Bridge
 */
export async function detectPlatform(): Promise<PlatformInfo> {
  if (cachedPlatformInfo) {
    return cachedPlatformInfo;
  }

  try {
    // Добавляем таймаут для VK Bridge запроса
    const launchParamsPromise = bridge.send('VKWebAppGetLaunchParams');
    const timeoutPromise = new Promise<never>((_, reject) =>
      setTimeout(() => reject(new Error('VK Bridge timeout')), 2000)
    );

    const launchParams = await Promise.race([launchParamsPromise, timeoutPromise]);
    const vkPlatform = launchParams.vk_platform as VKPlatform;

    const isIOS = vkPlatform.includes('iphone') || vkPlatform.includes('ipad');
    const isAndroid = vkPlatform.includes('android');
    const isDesktop = vkPlatform.includes('desktop');
    const isMobile = !isDesktop;

    let platform: Platform;
    if (isIOS) {
      platform = 'ios';
    } else if (isAndroid) {
      platform = 'android';
    } else {
      platform = 'vkcom'; // desktop
    }

    cachedPlatformInfo = {
      platform,
      vkPlatform,
      isIOS,
      isAndroid,
      isDesktop,
      isMobile,
    };

    console.log('[Platform] Detected:', cachedPlatformInfo);
    return cachedPlatformInfo;
  } catch (error) {
    console.error('[Platform] Detection failed, using fallback:', error);

    // Fallback based on user agent
    const ua = navigator.userAgent.toLowerCase();
    const isIOS = /iphone|ipad|ipod/.test(ua);
    const isAndroid = /android/.test(ua);
    const isMobile = /mobile/.test(ua);

    cachedPlatformInfo = {
      platform: isIOS ? 'ios' : isAndroid ? 'android' : 'vkcom',
      vkPlatform: 'desktop_web',
      isIOS,
      isAndroid,
      isDesktop: !isMobile,
      isMobile,
    };

    return cachedPlatformInfo;
  }
}

/**
 * React hook для использования в компонентах
 */
export function usePlatformDetection() {
  const [platformInfo, setPlatformInfo] = useState<PlatformInfo | null>(
    cachedPlatformInfo
  );

  useEffect(() => {
    detectPlatform().then(setPlatformInfo);
  }, []);

  return platformInfo;
}
