DROP TABLE IF EXISTS predictions;

DROP TABLE IF EXISTS matches;

DROP TABLE IF EXISTS telegram_users;

DROP TABLE IF EXISTS tournaments;

REASSIGN OWNED BY rushbanana_user TO postgres;
DROP OWNED BY rushbanana_user;
DROP ROLE IF EXISTS rushbanana_user;