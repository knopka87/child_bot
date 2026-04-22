-- Attempts table - unified table for help and check attempts
CREATE TABLE IF NOT EXISTS attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_profile_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE,

    -- Тип попытки
    attempt_type VARCHAR(20) NOT NULL CHECK (attempt_type IN ('help', 'check')),

    -- Статус
    status VARCHAR(20) NOT NULL DEFAULT 'created'
        CHECK (status IN ('created', 'processing', 'completed', 'failed')),

    -- Изображения (храним URL или путь к S3)
    task_image_url TEXT,
    answer_image_url TEXT,

    -- Результаты LLM (JSON)
    detect_result JSONB, -- DetectResponse
    parse_result JSONB,  -- ParseResponse
    hints_result JSONB,  -- HintResponse (для help)
    check_result JSONB,  -- CheckResponse (для check)

    -- Текущая подсказка (для help)
    current_hint_index INTEGER DEFAULT 0,

    -- Статистика
    hints_used INTEGER DEFAULT 0,
    time_spent_seconds INTEGER, -- время решения в секундах

    -- Результат (для статистики)
    is_correct BOOLEAN, -- для check
    has_errors BOOLEAN, -- для check

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_attempts_child_profile
    ON attempts (child_profile_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_attempts_status
    ON attempts (status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_attempts_type
    ON attempts (attempt_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_attempts_unfinished
    ON attempts (child_profile_id, status)
    WHERE status IN ('created', 'processing');

-- Триггер для обновления updated_at
CREATE TRIGGER attempts_updated_at
    BEFORE UPDATE ON attempts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); -- используем существующую функцию

-- Комментарии
COMMENT ON TABLE attempts IS 'Попытки решения задач (help и check)';
COMMENT ON COLUMN attempts.attempt_type IS 'Тип: help (подсказки) или check (проверка решения)';
COMMENT ON COLUMN attempts.status IS 'Статус: created, processing, completed, failed';
COMMENT ON COLUMN attempts.detect_result IS 'JSON результат Detect API';
COMMENT ON COLUMN attempts.parse_result IS 'JSON результат Parse API';
COMMENT ON COLUMN attempts.hints_result IS 'JSON результат Hint API';
COMMENT ON COLUMN attempts.check_result IS 'JSON результат CheckSolution API';
