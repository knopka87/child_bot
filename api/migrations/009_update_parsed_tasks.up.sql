-- Добавляем колонку updated_at
ALTER TABLE public.parsed_tasks
    ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now();

-- Обновляем существующие записи
UPDATE public.parsed_tasks
SET updated_at = created_at
WHERE updated_at IS NULL;

-- Создаем функцию для триггера
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'UPDATE') THEN
        NEW.updated_at = now();
    ELSIF (TG_OP = 'INSERT') THEN
        NEW.updated_at = now();
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Создаем триггер
DROP TRIGGER IF EXISTS update_parsed_tasks_updated_at ON public.parsed_tasks;
CREATE TRIGGER update_parsed_tasks_updated_at
    BEFORE INSERT OR UPDATE ON public.parsed_tasks
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();