import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import { visualizer } from 'rollup-plugin-visualizer';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    // Bundle analyzer - анализ размера bundle
    visualizer({
      open: false,
      filename: 'dist/stats.html',
      gzipSize: true,
      brotliSize: true,
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@/api': path.resolve(__dirname, './src/api'),
      '@/components': path.resolve(__dirname, './src/components'),
      '@/hooks': path.resolve(__dirname, './src/hooks'),
      '@/lib': path.resolve(__dirname, './src/lib'),
      '@/pages': path.resolve(__dirname, './src/pages'),
      '@/stores': path.resolve(__dirname, './src/stores'),
      '@/types': path.resolve(__dirname, './src/types'),
      '@/utils': path.resolve(__dirname, './src/utils'),
    },
  },
  server: {
    port: 3000,
    host: true, // Необходимо для доступа с мобильных устройств
    // https отключен по умолчанию для локальной разработки (для production будет HTTPS через nginx)
    allowedHosts: [
      'localhost',
      '.ngrok-free.dev', // Разрешаем все ngrok домены
      '.ngrok.io',
      '.loca.lt', // localtunnel
      '.trycloudflare.com', // cloudflare tunnel
    ],
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    // Оптимизация для production
    minify: 'esbuild',
    target: 'es2020',
    rollupOptions: {
      output: {
        format: 'es',
        // Code splitting для оптимизации bundle size
        manualChunks: {
          // React ecosystem
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          // VK SDK
          'vk-vendor': [
            '@vkontakte/vk-bridge',
            '@vkontakte/vk-bridge-mock',
            '@vkontakte/vkui',
            '@vkontakte/icons',
          ],
          // State management & data fetching
          'ui-vendor': ['zustand', 'axios'],
        },
      },
    },
    // Warning при превышении 1MB на chunk
    chunkSizeWarningLimit: 1000,
  },
  // Оптимизация для production
  esbuild: {
    target: 'es2020',
    // Временно отключено для отладки
    drop: [], // process.env.NODE_ENV === 'production' ? ['console', 'debugger'] : [],
  },
});
