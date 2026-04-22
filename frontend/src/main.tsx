import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';
import { ErrorBoundary } from '@/components/ErrorBoundary';

// VK Bridge инициализируется автоматически в @/lib/platform/bridge
// Mock будет использоваться автоматически в dev режиме вне VK

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ErrorBoundary>
      <App />
    </ErrorBoundary>
  </React.StrictMode>
);