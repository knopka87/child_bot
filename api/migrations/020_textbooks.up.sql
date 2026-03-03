-- Таблица учебников
CREATE TABLE IF NOT EXISTS textbooks (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    subject VARCHAR(50) NOT NULL,
    grade INT NOT NULL,
    authors VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    part INT DEFAULT NULL,
    year INT DEFAULT NULL,
    publisher VARCHAR(255) DEFAULT NULL,
    source_url VARCHAR(500) DEFAULT NULL,
    UNIQUE (subject, grade, authors, part)
);

-- Таблица задач из учебников
CREATE TABLE IF NOT EXISTS textbook_tasks (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    textbook_id BIGINT NOT NULL REFERENCES textbooks(id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    task_number VARCHAR(20) NOT NULL,
    task_order INT NOT NULL DEFAULT 0,
    condition_text TEXT,
    condition_html TEXT,
    solution_text TEXT,
    solution_html TEXT,
    hints_text TEXT,
    hints_html TEXT,
    has_sub_items BOOLEAN NOT NULL DEFAULT FALSE,
    sub_items_json JSONB DEFAULT NULL,
    source_url VARCHAR(500) DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_textbook_tasks_textbook_page ON textbook_tasks(textbook_id, page_number);
CREATE INDEX IF NOT EXISTS idx_textbook_tasks_task_number ON textbook_tasks(textbook_id, task_number);

-- Таблица картинок к задачам
CREATE TABLE IF NOT EXISTS textbook_task_images (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    task_id BIGINT NOT NULL REFERENCES textbook_tasks(id) ON DELETE CASCADE,
    image_type VARCHAR(20) NOT NULL,
    image_order INT NOT NULL DEFAULT 0,
    sub_item_letter VARCHAR(5) DEFAULT NULL,
    original_url VARCHAR(500) NOT NULL,
    local_path VARCHAR(500) DEFAULT NULL,
    alt_text VARCHAR(500) DEFAULT NULL,
    width INT DEFAULT NULL,
    height INT DEFAULT NULL,
    file_size INT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_textbook_task_images_task ON textbook_task_images(task_id, image_type, image_order);
