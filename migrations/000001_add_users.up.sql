CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "login_at" timestamptz NOT NULL DEFAULT now(),
    "email" varchar UNIQUE NOT NULL,
    "first_name" varchar,
    "last_name" varchar,
    "active" bool NOT NULL DEFAULT true,
    "role" varchar(20) NOT NULL
);