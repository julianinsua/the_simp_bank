-- +goose Up
ALTER TABLE "sessions" ADD CONSTRAINT "unique_user_sessions" UNIQUE ("username");

-- +goose Down
ALTER TABLE "sessions" DROP CONSTRAINT IF EXISTS "unique_user_sessions";
