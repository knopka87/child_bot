// src/components/VKOnlyAccess.tsx
import { Shield, ExternalLink } from 'lucide-react';

export function VKOnlyAccess() {
  const vkAppUrl = 'https://vk.com/app54517931';

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF] flex items-center justify-center p-6">
      <div className="bg-white rounded-3xl p-8 shadow-lg max-w-md w-full">
        <div className="flex flex-col items-center text-center">
          <div className="w-20 h-20 bg-[#6C5CE7]/10 rounded-full flex items-center justify-center mb-6">
            <Shield size={36} className="text-[#6C5CE7]" />
          </div>

          <h1 className="text-2xl font-bold text-[#2D3436] mb-3">
            Доступ только через VK
          </h1>

          <p className="text-[#636e72] text-sm mb-6 leading-relaxed">
            Приложение «Объяснятель ДЗ» работает только внутри социальной сети VK.
            Это необходимо для вашей безопасности и защиты данных.
          </p>

          <div className="bg-[#E3F2FD] border border-[#90CAF9] rounded-xl p-4 mb-6 w-full">
            <p className="text-[#1565C0] text-sm font-medium mb-2">
              💡 Как открыть приложение:
            </p>
            <ol className="text-[#1565C0] text-sm text-left space-y-2">
              <li>1. Откройте приложение VK на телефоне</li>
              <li>2. Перейдите в раздел «Сервисы» или «Игры»</li>
              <li>3. Найдите «Объяснятель ДЗ»</li>
              <li>4. Или используйте прямую ссылку ниже</li>
            </ol>
          </div>

          <a
            href={vkAppUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="w-full py-4 bg-[#6C5CE7] text-white rounded-2xl font-medium shadow-lg shadow-[#6C5CE7]/30 active:scale-[0.98] transition-transform flex items-center justify-center gap-2"
          >
            Открыть в VK
            <ExternalLink size={18} />
          </a>

          <div className="mt-6 p-4 bg-[#FFF3CD] border border-[#FFE69C] rounded-xl w-full">
            <p className="text-[#856404] text-xs leading-relaxed">
              <strong>⚠️ Важно:</strong> Убедитесь, что VPN выключен.
              VK блокирует доступ к приложениям при использовании VPN или прокси.
            </p>
          </div>

          <details className="mt-4 w-full">
            <summary className="text-xs text-[#B2BEC3] cursor-pointer hover:text-[#636e72] transition-colors">
              Почему только VK?
            </summary>
            <div className="mt-2 text-xs text-left text-[#636e72] bg-[#F5F6FA] p-3 rounded-lg">
              <p className="mb-2">
                Приложение использует VK Mini Apps API для:
              </p>
              <ul className="space-y-1 list-disc list-inside">
                <li>Безопасной авторизации без пароля</li>
                <li>Защиты персональных данных</li>
                <li>Интеграции с профилем VK</li>
                <li>Получения родительского согласия</li>
              </ul>
              <p className="mt-2">
                Открытие через браузер не поддерживается из соображений безопасности.
              </p>
            </div>
          </details>
        </div>
      </div>
    </div>
  );
}
