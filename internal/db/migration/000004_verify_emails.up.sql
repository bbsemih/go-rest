CREATE TABLE "verify_emails" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "email" VARCHAR NOT NULL,
    "secret_code" VARCHAR NOT NULL,
    "is_used" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" timestampz NOT NULL DEFAULT now(),
    "expired_at" timestampz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE;

ALTER TABLE "users" ADD COLUMN "is_email_verified" BOOLEAN NOT NULL DEFAULT FALSE;