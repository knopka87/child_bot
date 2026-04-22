-- Добавляем таблицу для хранения еженедельных HTML-отчётов
CREATE TABLE weekly_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE,
    report_date DATE NOT NULL,
    html_content TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Уникальный индекс для предотвращения дубликатов
CREATE UNIQUE INDEX idx_weekly_reports_user_date ON weekly_reports(user_id, report_date);

-- Индексы для быстрого поиска
CREATE INDEX idx_weekly_reports_user_id ON weekly_reports(user_id);
CREATE INDEX idx_weekly_reports_report_date ON weekly_reports(report_date);