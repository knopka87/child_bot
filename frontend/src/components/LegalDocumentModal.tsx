// src/components/LegalDocumentModal.tsx
import { useState, useEffect } from 'react';
import { CustomModal } from '@/components/ui/Modal/CustomModal';
import { apiClient } from '@/api/client';
import ReactMarkdown from 'react-markdown';

interface LegalDocument {
  id: string;
  version: string;
  title: string;
  content: string;
  language: string;
  effectiveDate: string;
  lastUpdated: string;
}

interface LegalDocumentModalProps {
  type: 'privacy' | 'terms' | null;
  isOpen: boolean;
  onClose: () => void;
}

export function LegalDocumentModal({ type, isOpen, onClose }: LegalDocumentModalProps) {
  const [document, setDocument] = useState<LegalDocument | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isOpen || !type) {
      setDocument(null);
      setError(null);
      return;
    }

    const fetchDocument = async () => {
      setLoading(true);
      setError(null);
      try {
        const endpoint = type === 'privacy' ? '/legal/privacy' : '/legal/terms';
        const data = await apiClient.get<LegalDocument>(endpoint);
        setDocument(data);
      } catch (err) {
        console.error('[LegalDocumentModal] Failed to fetch:', err);
        setError('Не удалось загрузить документ');
      } finally {
        setLoading(false);
      }
    };

    fetchDocument();
  }, [isOpen, type]);

  return (
    <CustomModal isOpen={isOpen} onClose={onClose}>
      <div style={{ padding: '24px' }}>
        {loading && (
          <div style={{ textAlign: 'center', padding: '40px 0' }}>
            <p style={{ color: '#6C757D' }}>Загрузка...</p>
          </div>
        )}

        {error && (
          <div style={{ textAlign: 'center', padding: '40px 0' }}>
            <p style={{ color: '#DC3545', marginBottom: '16px' }}>{error}</p>
            <button
              onClick={onClose}
              style={{
                padding: '12px 24px',
                background: '#6C5CE7',
                color: 'white',
                border: 'none',
                borderRadius: '12px',
                fontSize: '16px',
                fontWeight: 600,
                cursor: 'pointer',
              }}
            >
              Закрыть
            </button>
          </div>
        )}

        {!loading && !error && document && (
          <>
            {/* Заголовок */}
            <h2
              style={{
                fontSize: '20px',
                fontWeight: 700,
                color: '#2D3436',
                marginBottom: '8px',
              }}
            >
              {document.title}
            </h2>

            {/* Версия */}
            <p
              style={{
                fontSize: '12px',
                color: '#6C757D',
                marginBottom: '24px',
              }}
            >
              Версия {document.version}
            </p>

            {/* Контент с прокруткой */}
            <div
              style={{
                maxHeight: '400px',
                overflowY: 'auto',
                fontSize: '14px',
                lineHeight: '1.6',
                color: '#495057',
                paddingRight: '8px',
              }}
            >
              <ReactMarkdown>{document.content}</ReactMarkdown>
            </div>

            {/* Даты */}
            <div
              style={{
                marginTop: '24px',
                paddingTop: '16px',
                borderTop: '1px solid #E9ECEF',
              }}
            >
              <p style={{ fontSize: '12px', color: '#6C757D', marginBottom: '4px' }}>
                Дата вступления в силу:{' '}
                {new Date(document.effectiveDate).toLocaleDateString('ru-RU')}
              </p>
              <p style={{ fontSize: '12px', color: '#6C757D' }}>
                Последнее обновление:{' '}
                {new Date(document.lastUpdated).toLocaleDateString('ru-RU')}
              </p>
            </div>

            {/* Кнопка закрытия */}
            <button
              onClick={onClose}
              style={{
                width: '100%',
                padding: '14px',
                background: '#6C5CE7',
                color: 'white',
                border: 'none',
                borderRadius: '16px',
                fontSize: '16px',
                fontWeight: 600,
                marginTop: '24px',
                cursor: 'pointer',
                transition: 'transform 0.2s',
              }}
              onMouseDown={(e) => {
                e.currentTarget.style.transform = 'scale(0.98)';
              }}
              onMouseUp={(e) => {
                e.currentTarget.style.transform = 'scale(1)';
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.transform = 'scale(1)';
              }}
            >
              Понятно
            </button>
          </>
        )}
      </div>
    </CustomModal>
  );
}
