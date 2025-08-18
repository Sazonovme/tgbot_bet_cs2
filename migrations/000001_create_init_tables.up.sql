/*
migrate create -ext sql -dir migrations -seq create_users_table
migrate -path ../../migrations -database "postgres://postgres@localhost:5432/rushbanana_db?sslmode=disable" up
*/

-- Таблица пользователей телеграма
CREATE TABLE IF NOT EXISTS telegram_users (
   chat_id BIGINT PRIMARY KEY,
   user_id BIGINT UNIQUE NOT NULL,
   username TEXT UNIQUE,
   is_active BOOLEAN DEFAULT true,
   created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица турниров
CREATE TABLE IF NOT EXISTS tournaments(
   id SERIAL PRIMARY KEY,
   name VARCHAR(300) NOT NULL,
   is_active BOOLEAN
);

-- Таблица матчей
CREATE TABLE IF NOT EXISTS matches (
   id SERIAL PRIMARY KEY,
   tournament_id INT REFERENCES tournaments(id),
   name VARCHAR(300) NOT NULL,
   team_1 VARCHAR(100) NOT NULL,
   team_2 VARCHAR(100) NOT NULL,
   date TIMESTAMPTZ NOT NULL,
   result VARCHAR(5) -- например: '2-1'
);

-- Таблица предсказаний
CREATE TABLE IF NOT EXISTS predictions (
   chat_id BIGINT NOT NULL REFERENCES telegram_users(chat_id),
   match_id INT REFERENCES matches(id),
   prediction VARCHAR(5) NOT NULL, -- например: '2-1' или "1"
   UNIQUE (chat_id, match_id)
);

-- Пользователь базы
CREATE ROLE rushbanana_user WITH LOGIN PASSWORD 'secret_password';

-- Доступ к таблицам
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO rushbanana_user;

-- Доступ к sequence (нужен для SERIAL / IDENTITY)
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO rushbanana_user;

-- Дефолтные права для новых таблиц
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO rushbanana_user;

-- Дефолтные права для новых sequence
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO rushbanana_user;
