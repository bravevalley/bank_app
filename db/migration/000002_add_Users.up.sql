CREATE TABLE IF NOT EXISTS "users" (
    "username" VARCHAR(50) PRIMARY KEY,
    "hashed_password" VARCHAR NOT NULL,
    "full_name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(100) UNIQUE NOT NULL,
    "password_last_changed" TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00Z',    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);


ALTER TABLE "account" ADD FOREIGN KEY ("name") REFERENCES "users" ("username");

ALTER TABLE "account" ADD CONSTRAINT "name_currency_key" UNIQUE("name", "currency");
