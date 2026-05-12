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
import { getCurrentChildProfileId } from '@/lib/auth';
import { VKBridgeError } from '@/components/VKBridgeError';
import { VKOnlyAccess } from '@/components/VKOnlyAccess';

type Appearance = 'light' | 'dark';

function AppInitializer() {
  const navigate = useNavigate();
  const [isInitialized, setIsInitialized] = useState(false);
  const [vkBridgeError, setVKBridgeError] = useState<Error | null>(null);
  const [requiresVKAccess, setRequiresVKAccess] = useState(false);
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
        // КРИТИЧЕСКИ ВАЖНО: Используем синхронную запись, чтобы избежать race condition
        // между async vkStorage.setItem и sync localStorage.getItem в PlatformBridge
        localStorage.setItem('platform_id', platformType);
        console.log('[App] Platform:', platformType);

        // Проверяем текущий путь - legal pages доступны без авторизации
        const currentPath = window.location.pathname;
        const isLegalPage = currentPath.startsWith('/legal/');

        if (isLegalPage) {
          console.log('[App] Legal page detected, skipping auth check');
          setIsInitialized(true);
          return;
        }

        // КРИТИЧЕСКИ ВАЖНО: Проверяем профиль через backend по VK user ID
        // getCurrentChildProfileId() делает запрос: GET /profiles/by-platform?platform_id=vk&platform_user_id=XXX
        // Если профиль найден в БД - значит онбординг пройден
        // Если профиль не найден (404) - требуем онбординг
        const childProfileId = await getCurrentChildProfileId();

        console.log('[App] Child Profile ID (validated from DB):', childProfileId || 'not found');

        if (!childProfileId) {
          console.log('[App] Profile not found in DB, redirecting to onboarding...');
          navigate('/onboarding', { replace: true });
        } else {
          console.log('[App] Profile found in DB, loading home...');
          // Обновляем флаг в локальном хранилище для совместимости
          await vkStorage.setItem(storageKeys.ONBOARDING_COMPLETED, 'true');
        }
      } catch (error) {
        console.error('[App] Initialization error:', error);

        // Проверяем не legal ли это страница
        const isLegalPage = window.location.pathname.startsWith('/legal/');

        // Проверяем тип ошибки
        if (error instanceof Error) {
          // Проверка на ошибку "Доступ только через VK"
          if (error.message.includes('VK_ONLY_ACCESS')) {
            console.error('[App] ❌ Application opened outside VK');
            setRequiresVKAccess(true);
            setIsInitialized(true);
            return;
          }

          // Проверка на ошибку VK Bridge (VPN, сеть)
          const isVKBridgeError =
            error.message.includes('VK Bridge') ||
            error.message.includes('timeout') ||
            error.message.includes('VKWebApp');

          if (isVKBridgeError && !isLegalPage) {
            // Показываем ошибку VK Bridge (возможно VPN)
            setVKBridgeError(error);
            setIsInitialized(true);
            return;
          }
        }

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

  // Показываем сообщение "Доступ только через VK"
  if (requiresVKAccess) {
    return <VKOnlyAccess />;
  }

  // Показываем ошибку VK Bridge (VPN или другие проблемы)
  if (vkBridgeError) {
    return <VKBridgeError error={vkBridgeError} />;
  }

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

        // КРИТИЧЕСКИ ВАЖНО: Используем синхронную запись для избежания race condition
        localStorage.setItem('platform_id', apiPlatformId);
        console.log('[App] platform_id saved to storage:', apiPlatformId);
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
        <ConfigProvider colorScheme={appearance} platform={platform}>
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
