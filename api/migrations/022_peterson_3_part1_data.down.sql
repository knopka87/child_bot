-- Удаление данных части 1 учебника Петерсон 3 класс
-- Удаление индексов (каскадно удалится при удалении задач, но для явности)
DELETE FROM textbook_task_index
WHERE textbook_id = (SELECT id FROM textbooks WHERE subject = 'math' AND grade = 3 AND authors = 'Петерсон Л.Г.' AND part = 1);

-- Удаление картинок (каскадно удалится при удалении задач, но для явности)
DELETE FROM textbook_task_images
WHERE task_id IN (
    SELECT t.id FROM textbook_tasks t
    INNER JOIN textbooks tb ON t.textbook_id = tb.id
    WHERE tb.subject = 'math' AND tb.grade = 3 AND tb.authors = 'Петерсон Л.Г.' AND tb.part = 1
);

-- Удаление задач
DELETE FROM textbook_tasks
WHERE textbook_id = (SELECT id FROM textbooks WHERE subject = 'math' AND grade = 3 AND authors = 'Петерсон Л.Г.' AND part = 1);

-- Удаление учебника
DELETE FROM textbooks WHERE subject = 'math' AND grade = 3 AND authors = 'Петерсон Л.Г.' AND part = 1;
