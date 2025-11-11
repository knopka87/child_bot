-- Добавляем колонку session_id
ALTER TABLE hints_cache
    ADD COLUMN IF NOT EXISTS session_id text;

alter table hints_cache
    drop column IF exists image_hash;

TRUNCATE Table hints_cache;

alter table hints_cache
    add constraint hints_cache_sid_level_idx
        unique (session_id, level);