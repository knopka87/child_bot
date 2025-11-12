alter table hints_cache
    alter column level type text using level::text;
