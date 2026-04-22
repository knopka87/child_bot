// src/api/profile.ts
import { apiClient } from './client';
import type {
  ProfileData,
  HistoryAttempt,
  HistoryFilters,
  ReportSettings,
  WeeklyReport,
} from '@/types/profile';

// Тип ответа от бэкенда
interface BackendProfileResponse {
  id: string;
  display_name: string;
  avatar_id: string;
  avatar_url: string;
  grade: number;
  subscription: {
    status: 'trial' | 'active' | 'expired' | 'cancelled';
    plan_id?: string;
    plan_name?: string;
    trial_days_remaining?: number;
    expires_at?: string;
  };
}

export const profileAPI = {
  /**
   * Получить профиль
   */
  async getProfile(_childProfileId: string): Promise<ProfileData> {
    const response = await apiClient.get<BackendProfileResponse>('/profile');

    return {
      id: response.id,
      displayName: response.display_name,
      avatarId: response.avatar_id,
      avatarUrl: response.avatar_url,
      grade: response.grade,
      subscription: {
        status: response.subscription.status,
        planId: response.subscription.plan_id,
        planName: response.subscription.plan_name,
        trialDaysRemaining: response.subscription.trial_days_remaining,
        expiresAt: response.subscription.expires_at,
      },
    };
  },

  /**
   * Получить историю попыток
   * Backend эндпоинт: GET /profile/history
   * child_profile_id берётся из middleware автоматически
   */
  async getHistory(
    _childProfileId: string,
    filters?: HistoryFilters
  ): Promise<HistoryAttempt[]> {
    interface BackendHistoryAttempt {
      id: string;
      mode: 'help' | 'check';
      status: 'success' | 'error' | 'in_progress';
      scenario_type?: 'single_photo' | 'two_photo';
      created_at: string;
      completed_at?: string;
      images: Array<{
        id: string;
        role: 'task' | 'answer' | 'single';
        url: string;
        thumbnail_url: string;
      }>;
      result?: {
        status: 'correct' | 'has_errors' | 'processing';
        error_count?: number;
        feedback?: Array<{
          id: string;
          step_number?: number;
          line_reference?: string;
          description: string;
          location_type: 'line' | 'step' | 'general';
        }>;
        summary?: string;
      };
      hints_used?: number;
    }

    const params: Record<string, string> = {};

    if (filters?.mode && filters.mode !== 'all') {
      params.mode = filters.mode;
    }

    if (filters?.status && filters.status !== 'all') {
      params.status = filters.status;
    }

    if (filters?.dateFrom) {
      params.date_from = filters.dateFrom;
    }

    if (filters?.dateTo) {
      params.date_to = filters.dateTo;
    }

    const response = await apiClient.get<BackendHistoryAttempt[]>(
      '/profile/history',
      { params }
    );

    console.log('[profileAPI] History response:', response);

    return response.map(
      (attempt): HistoryAttempt => ({
        id: attempt.id,
        mode: attempt.mode,
        status: attempt.status,
        scenarioType: attempt.scenario_type,
        createdAt: attempt.created_at,
        completedAt: attempt.completed_at,
        images: attempt.images.map((img) => ({
          id: img.id,
          role: img.role,
          url: img.url,
          thumbnailUrl: img.thumbnail_url,
        })),
        result: attempt.result
          ? {
              status: attempt.result.status,
              errorCount: attempt.result.error_count,
              feedback: attempt.result.feedback?.map((fb) => ({
                id: fb.id,
                stepNumber: fb.step_number,
                lineReference: fb.line_reference,
                description: fb.description,
                locationType: fb.location_type,
              })),
              summary: attempt.result.summary,
            }
          : undefined,
        hintsUsed: attempt.hints_used,
      })
    );
  },

  /**
   * Получить детали попытки
   */
  async getHistoryDetail(
    _childProfileId: string,
    attemptId: string
  ): Promise<HistoryAttempt> {
    const history = await this.getHistory(_childProfileId);
    const attempt = history.find((a) => a.id === attemptId);

    if (!attempt) {
      throw new Error(`Attempt ${attemptId} not found in history`);
    }

    return attempt;
  },

  /**
   * Получить настройки отчётов
   */
  async getReportSettings(childProfileId: string): Promise<ReportSettings> {
    return apiClient.get<ReportSettings>(`/reports/${childProfileId}/settings`);
  },

  /**
   * Обновить настройки отчётов
   */
  async updateReportSettings(
    childProfileId: string,
    settings: Partial<ReportSettings>
  ): Promise<void> {
    return apiClient.patch<void>(`/reports/${childProfileId}/settings`, settings);
  },

  /**
   * Получить архив отчётов
   */
  async getReportArchive(childProfileId: string): Promise<WeeklyReport[]> {
    return apiClient.get<WeeklyReport[]>(`/reports/${childProfileId}/archive`);
  },

  /**
   * Скачать отчёт
   */
  async downloadReport(
    childProfileId: string,
    reportId: string
  ): Promise<Blob> {
    const response = await apiClient.get<Blob>(
      `/reports/${childProfileId}/${reportId}/download`,
      { responseType: 'blob' }
    );
    return response;
  },

  /**
   * Отправить тестовый отчёт
   */
  async sendTestReport(childProfileId: string): Promise<void> {
    return apiClient.post<void>(`/reports/${childProfileId}/send-test`);
  },

  /**
   * Отправить сообщение в поддержку
   */
  async sendSupportMessage(
    childProfileId: string,
    message: string
  ): Promise<void> {
    return apiClient.post<void>(`/support/messages`, {
      childProfileId,
      message,
    });
  },
};
