// src/config/assets.ts
/**
 * Централизованная конфигурация путей к ассетам
 * Используйте эти константы вместо hardcoded путей
 */

export const ASSETS = {
  images: {
    villainDefeated: '/assets/villain-defeated.png',
    mascot: '/assets/mascot.png',
    placeholder: '/assets/placeholder.png',
  },
  icons: {
    coin: '/assets/icons/coin.svg',
    star: '/assets/icons/star.svg',
    trophy: '/assets/icons/trophy.svg',
  },
} as const;

// Type-safe helper для получения пути к ассету
export type AssetPath = typeof ASSETS;
