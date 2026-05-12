// src/components/VKBridgeError.tsx
import { AlertTriangle } from 'lucide-react';

interface VKBridgeErrorProps {
  error?: Error;
}

export function VKBridgeError({ error }: VKBridgeErrorProps) {
  const isVPNRelated = error?.message?.includes('timeout') ||
                        error?.message?.includes('network') ||
                        error?.message?.includes('Connection');

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF] flex items-center justify-center p-6">
      <div className="bg-white rounded-3xl p-8 shadow-lg max-w-md w-full">
        <div className="flex flex-col items-center text-center">
          <div className="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center mb-6">
            <AlertTriangle size={36} className="text-red-500" />
          </div>

          <h1 className="text-2xl font-bold text-[#2D3436] mb-3">
            Не удалось подключиться к VK
          </h1>

          {isVPNRelated ? (
            <>
              <p className="text-[#636e72] text-sm mb-6 leading-relaxed">
                VK Mini Apps не работают с включённым VPN или прокси-сервером.
                Это сделано для защиты от мошенничества.
              </p>

              <div className="bg-[#FFF3CD] border border-[#FFE69C] rounded-xl p-4 mb-6 w-full">
                <p className="text-[#856404] text-sm font-medium mb-2">
                  💡 Что делать:
                </p>
                <ol className="text-[#856404] text-sm text-left space-y-1.5">
                  <li>1. Отключите VPN или прокси</li>
                  <li>2. Перезапустите приложение VK</li>
                  <li>3. Попробуйте снова</li>
                </ol>
              </div>
            </>
          ) : (
            <>
              <p className="text-[#636e72] text-sm mb-6 leading-relaxed">
                Возникла проблема при подключении к серверам VK.
                Проверьте подключение к интернету и попробуйте снова.
              </p>

              <div className="bg-[#E3F2FD] border border-[#90CAF9] rounded-xl p-4 mb-6 w-full">
                <p className="text-[#1565C0] text-sm font-medium mb-2">
                  🔍 Возможные причины:
                </p>
                <ul className="text-[#1565C0] text-sm text-left space-y-1.5">
                  <li>• Проблемы с интернет-соединением</li>
                  <li>• Приложение открыто не через VK</li>
                  <li>• Временные проблемы на серверах VK</li>
                </ul>
              </div>
            </>
          )}

          <button
            onClick={() => window.location.reload()}
            className="w-full py-4 bg-[#6C5CE7] text-white rounded-2xl font-medium shadow-lg shadow-[#6C5CE7]/30 active:scale-[0.98] transition-transform"
          >
            Попробовать снова
          </button>

          {error && (
            <details className="mt-4 w-full">
              <summary className="text-xs text-[#B2BEC3] cursor-pointer hover:text-[#636e72] transition-colors">
                Техническая информация
              </summary>
              <pre className="mt-2 text-xs text-left bg-[#F5F6FA] p-3 rounded-lg overflow-x-auto text-[#636e72]">
                {error.message}
              </pre>
            </details>
          )}
        </div>
      </div>
    </div>
  );
}
