CREATE DATABASE dtms;

\connect dtms

CREATE SCHEMA tasks;
COMMENT ON SCHEMA tasks IS E'Хранит информацию о статуах задач';

-- Создание таблицы tasks.history
CREATE TABLE tasks.history (
    id SERIAL PRIMARY KEY,
    task_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    task_data JSONB NOT NULL,
    action VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,

    CONSTRAINT unique_task_action UNIQUE (task_id, timestamp, user_id)
);

-- Добавление комментариев к таблице и полям для лучшего понимания их назначения
COMMENT ON TABLE tasks.history IS 'Таблица для отслеживания изменений статуса задач в сервисе трекера задач.';
COMMENT ON COLUMN tasks.history.task_id IS 'Идентификатор задачи, связанной с действием.';
COMMENT ON COLUMN tasks.history.timestamp IS 'Временная метка выполнения действия.';
COMMENT ON COLUMN tasks.history.user_id IS 'Идентификатор пользователя, выполнившего действие.';
