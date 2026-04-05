// src/App.tsx
import { useEffect, useState, useRef } from 'react';
import { BrowserRouter, useNavigate } from 'react-router-dom';
import {
  ConfigProvider,
  AdaptivityProvider,
  AppRoot,
  SplitLayout,
  SplitCol,
} from '@vkontakte/vkui';
import '@vkontakte/vkui/dist/vkui.css';

import { detectPlatform, type Platform } from '@/lib/platform/platform-detection';
import { AppRoutes } from '@/routes';
import { AnalyticsProvider } from '@/contexts/AnalyticsContext';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { Spinner } from '@/components/ui/Spinner';
import { PlatformBridge } from '@/services/platform/PlatformBridge';

type Appearance = 'light' | 'dark';

function AppInitializer() {
  const navigate = useNavigate();
  const [isInitialized, setIsInitialized] = useState(false);
  const hasCheckedRef = useRef(false);

  useEffect(() => {
    // Выполняем проверку только один раз при монтировании
    if (hasCheckedRef.current) return;
    hasCheckedRef.current = true;

    const initializeApp = async () => {
      try {
        console.log('[App] Checking authentication...');

        // Определяем и сохраняем platform ID
        const platformBridge = new PlatformBridge();
        const platformType = platformBridge.getPlatformType();
        await vkStorage.setItem('platform_id', platformType);
        console.log('[App] Platform:', platformType);

        // Проверяем текущий путь - legal pages доступны без авторизации
        const currentPath = window.location.pathname;
        const isLegalPage = currentPath.startsWith('/legal/');

        if (isLegalPage) {
          console.log('[App] Legal page detected, skipping auth check');
          setIsInitialized(true);
          return;
        }

        // Проверяем наличие профиля и завершения онбординга
        const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
        const onboardingCompleted = await vkStorage.getItem(storageKeys.ONBOARDING_COMPLETED);

        console.log('[App] Child Profile ID:', childProfileId || 'missing');
        console.log('[App] Onboarding:', onboardingCompleted || 'not completed');

        if (!childProfileId || !onboardingCompleted) {
          console.log('[App] Redirecting to onboarding...');
          navigate('/onboarding', { replace: true });
        } else {
          console.log('[App] User authenticated, loading home...');
        }
      } catch (error) {
        console.error('[App] Initialization error:', error);
        // Проверяем не legal ли это страница
        const isLegalPage = window.location.pathname.startsWith('/legal/');
        if (!isLegalPage) {
          navigate('/onboarding', { replace: true });
        }
      } finally {
        setIsInitialized(true);
      }
    };

    // Таймаут для инициализации
    const timeout = setTimeout(() => {
      console.warn('[App] Initialization timeout');
      setIsInitialized(true);
      const isLegalPage = window.location.pathname.startsWith('/legal/');
      if (!isLegalPage) {
        navigate('/onboarding', { replace: true });
      }
    }, 3000);

    initializeApp().then(() => clearTimeout(timeout));

    return () => clearTimeout(timeout);
  }, [navigate]);

  if (!isInitialized) {
    return (
      <div className="app-desktop-wrapper">
        <div className="app-desktop-container">
          <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh'
          }}>
            <Spinner size="lg" />
          </div>
        </div>
      </div>
    );
  }

  return <AppRoutes />;
}

export default function App() {
  const [platform, setPlatform] = useState<Platform>('android');
  const [appearance, setAppearance] = useState<Appearance>('light');
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Определяем платформу с таймаутом
    const platformDetectionTimeout = setTimeout(() => {
      console.warn('[App] Platform detection timeout, using default');
      setIsLoading(false);
    }, 5000);

    detectPlatform()
      .then(async (info) => {
        clearTimeout(platformDetectionTimeout);
        setPlatform(info.platform);
        console.log('[App] Platform detected:', info.platform);

        // Сохраняем platform_id в storage для APIClient
        // Используем PlatformBridge для получения правильного API platform_id
        const platformBridge = new PlatformBridge();
        const apiPlatformId = platformBridge.getPlatformType();

        try {
          await vkStorage.setItem('platform_id', apiPlatformId);
          console.log('[App] platform_id saved to storage:', apiPlatformId);
        } catch (error) {
          console.error('[App] Failed to save platform_id:', error);
          // Fallback: сохраняем напрямую в localStorage
          localStorage.setItem('platform_id', apiPlatformId);
        }
      })
      .catch((error) => {
        clearTimeout(platformDetectionTimeout);
        console.error('[App] Platform detection failed:', error);
        // Продолжаем с дефолтной платформой, сохраняем 'web' в storage
        localStorage.setItem('platform_id', 'web');
      })
      .finally(() => {
        setIsLoading(false);
      });

    // Определяем тему
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    setAppearance(isDark ? 'dark' : 'light');

    return () => clearTimeout(platformDetectionTimeout);
  }, []);

  if (isLoading) {
    return (
      <div className="app-desktop-wrapper">
        <div className="app-desktop-container">
          <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
            background: 'var(--vkui--color_background)'
          }}>
            <Spinner size="lg" />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="app-desktop-wrapper">
      <div className="app-desktop-container">
        <ConfigProvider appearance={appearance} platform={platform}>
          <AdaptivityProvider>
            <AppRoot>
              <SplitLayout>
                <SplitCol>
                  <AnalyticsProvider>
                    <BrowserRouter>
                      <AppInitializer />
                    </BrowserRouter>
                  </AnalyticsProvider>
                </SplitCol>
              </SplitLayout>
            </AppRoot>
          </AdaptivityProvider>
        </ConfigProvider>
      </div>
    </div>
  );
}
