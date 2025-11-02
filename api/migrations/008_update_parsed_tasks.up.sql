-- Переименования (по одному действию на команду)
ALTER TABLE public.parsed_tasks
    RENAME COLUMN raw_text TO raw_task_text;

ALTER TABLE public.parsed_tasks
    RENAME COLUMN subject_hint TO subject;

ALTER TABLE public.parsed_tasks
    RENAME COLUMN confirmation_needed TO needs_user_confirmation;

-- Добавление новых колонок (идемпотентно)
ALTER TABLE public.parsed_tasks
    ADD COLUMN IF NOT EXISTS task_type text,
    ADD COLUMN IF NOT EXISTS combined_subpoints boolean;

-- Бэкофилл значений перед установкой NOT NULL/DEFAULT
UPDATE public.parsed_tasks
SET combined_subpoints = true
WHERE combined_subpoints IS NULL;

-- Жёсткие ограничения и дефолты
ALTER TABLE public.parsed_tasks
    ALTER COLUMN raw_task_text SET NOT NULL,
    ALTER COLUMN needs_user_confirmation SET DEFAULT true,
    ALTER COLUMN needs_user_confirmation SET NOT NULL,
    ALTER COLUMN combined_subpoints SET DEFAULT true,
    ALTER COLUMN combined_subpoints SET NOT NULL;

-- CHECK для subject
ALTER TABLE public.parsed_tasks
    ADD CONSTRAINT parsed_tasks_subject_enum CHECK (
        subject IN ('math','russian','generic')
    );
