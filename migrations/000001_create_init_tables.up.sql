/*
migrate create -ext sql -dir migrations -seq create_users_table
migrate -path ../../migrations -database "postgres://postgres@localhost:5432/rushbanana_db?sslmode=disable" up
*/

CREATE TABLE IF NOT EXISTS telegram_users (
   chat_id BIGINT PRIMARY KEY,
   username TEXT UNIQUE,
   first_name TEXT,
   last_name TEXT,
   is_active BOOLEAN DEFAULT true,
   created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tournaments (
   id serial PRIMARY KEY,
   name VARCHAR (300) NOT NULL,
   is_active BOOLEAN
);

CREATE TABLE IF NOT EXISTS matches (
   id SERIAL PRIMARY KEY,
   tournament_id INT REFERENCES tournaments(id),
   name VARCHAR(300) NOT NULL,
   team_1 VARCHAR(100) NOT NULL,
   team_2 VARCHAR(100) NOT NULL,
   date TIMESTAMPTZ NOT NULL,
   result VARCHAR(5) -- например: '2-1'
);

CREATE TABLE IF NOT EXISTS predictions (
   username TEXT NOT NULL REFERENCES telegram_users(username),
   match_id INT REFERENCES matches(id),
   prediction VARCHAR (5) NOT NULL -- например: '2-1' или "1"
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_prediction 
ON predictions (username, match_id);

CREATE ROLE rushbanana_user WITH LOGIN PASSWORD 'secret_password';

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO rushbanana_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO rushbanana_user;