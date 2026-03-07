-- Индексная таблица для быстрого поиска задач
CREATE TABLE IF NOT EXISTS textbook_task_index (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL UNIQUE REFERENCES textbook_tasks(id) ON DELETE CASCADE,
    textbook_id BIGINT NOT NULL REFERENCES textbooks(id) ON DELETE CASCADE,
    grade INT NOT NULL,

    -- Нормализованный хэш условия для точного поиска
    normalized_hash VARCHAR(64) NOT NULL,

    -- Числа из условия (для быстрой фильтрации)
    numbers_signature VARCHAR(255) DEFAULT NULL,

    -- Ключевые слова (JSON массив)
    keywords JSONB DEFAULT NULL,

    -- Нормализованный текст (для fulltext поиска)
    normalized_text TEXT,

    -- Вектор для полнотекстового поиска
    search_vector TSVECTOR,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_textbook_task_index_hash ON textbook_task_index(normalized_hash);
CREATE INDEX IF NOT EXISTS idx_textbook_task_index_grade ON textbook_task_index(grade);
CREATE INDEX IF NOT EXISTS idx_textbook_task_index_numbers ON textbook_task_index(numbers_signature);
CREATE INDEX IF NOT EXISTS idx_textbook_task_index_textbook_grade ON textbook_task_index(textbook_id, grade);
CREATE INDEX IF NOT EXISTS idx_textbook_task_index_search ON textbook_task_index USING GIN(search_vector);

-- Триггер для автоматического обновления search_vector
CREATE OR REPLACE FUNCTION update_textbook_task_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('russian', COALESCE(NEW.normalized_text, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_textbook_task_index_search ON textbook_task_index;
CREATE TRIGGER trg_textbook_task_index_search
    BEFORE INSERT OR UPDATE ON textbook_task_index
    FOR EACH ROW
    EXECUTE FUNCTION update_textbook_task_search_vector();
