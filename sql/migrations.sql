CREATE DATABASE example_data
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LOCALE_PROVIDER = 'libc'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

COMMENT ON DATABASE example_data
    IS 'Main database for example application';

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product TEXT NOT NULL,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    user_uuid UUID NOT NULL,
    CONSTRAINT fk_orders_user
    FOREIGN KEY (user_uuid)
    REFERENCES users (uuid)
    ON DELETE RESTRICT
    ON UPDATE CASCADE
);
