-- +goose Up
CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "username" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "client_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT FALSE,
    "expires_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

-- +goose Down
DROP TABLE IF EXISTS "sessions";
