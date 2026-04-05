import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ChevronLeft } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import { apiClient } from '@/api/client';

interface LegalDocument {
  id: string;
  version: string;
  title: string;
  content: string;
  language: string;
  effectiveDate: string;
  lastUpdated: string;
}

export function PrivacyPolicyPage() {
  const navigate = useNavigate();
  const [document, setDocument] = useState<LegalDocument | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDocument = async () => {
      try {
        const data = await apiClient.get<LegalDocument>('/legal/privacy');
        setDocument(data);
      } catch (err) {
        console.error('[PrivacyPolicy] Failed to fetch:', err);
        setError('Не удалось загрузить политику конфиденциальности');
      } finally {
        setLoading(false);
      }
    };

    fetchDocument();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-[#E8E4FF] to-[#F0ECFF] flex items-center justify-center">
        <div className="text-gray-600">Загрузка...</div>
      </div>
    );
  }

  if (error || !document) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-[#E8E4FF] to-[#F0ECFF] flex flex-col items-center justify-center p-4">
        <div className="text-red-600 mb-4">{error || 'Документ не найден'}</div>
        <button
          onClick={() => navigate(-1)}
          className="px-6 py-3 bg-[#6C5CE7] text-white rounded-2xl font-semibold active:scale-95 transition-transform"
        >
          Назад
        </button>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#E8E4FF] to-[#F0ECFF]">
      {/* Header */}
      <div className="bg-white shadow-sm sticky top-0 z-10">
        <div className="flex items-center gap-4 px-4 py-3">
          <button
            onClick={() => navigate(-1)}
            className="p-2 -ml-2 text-[#6C5CE7] active:scale-95 transition-transform"
          >
            <ChevronLeft size={24} />
          </button>
          <div>
            <h1 className="text-lg font-bold text-[#2D3436]">{document.title}</h1>
            <p className="text-xs text-gray-500">Версия {document.version}</p>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="px-4 py-6 max-w-3xl mx-auto">
        <div className="bg-white rounded-2xl shadow-sm p-6">
          <div className="prose prose-sm max-w-none">
            <ReactMarkdown>{document.content}</ReactMarkdown>
          </div>

          <div className="mt-8 pt-6 border-t border-gray-200">
            <p className="text-sm text-gray-500">
              Дата вступления в силу: {new Date(document.effectiveDate).toLocaleDateString('ru-RU')}
            </p>
            <p className="text-sm text-gray-500 mt-1">
              Последнее обновление: {new Date(document.lastUpdated).toLocaleDateString('ru-RU')}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
