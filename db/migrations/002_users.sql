-- +goose Up
CREATE TABLE "users" (
  "username" varchar NOT NULL PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamp NOT NULL DEFAULT ('0001-01-01 00:00:00'),
  "created_at" timestamp NOT NULL DEFAULT (now())
);

-- CREATE UNIQUE INDEX ON "account" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- +goose Down
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

DROP TABLE IF EXISTS "users";
