BEGIN;

CREATE TABLE "author" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "refresh_token" VARCHAR(255) NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

CREATE TABLE newsletter (
    "id" UUID PRIMARY KEY,
    "author_id" UUID NOT NULL REFERENCES "author" (id) ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "description" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

CREATE TABLE email (
    "id" UUID PRIMARY KEY,
    "newsletter_id" UUID NOT NULL REFERENCES "newsletter" (id) ON DELETE CASCADE,
    "subject" VARCHAR(255) NOT NULL,
    "html_content" VARCHAR(1024) NOT NULL,
    "created_at" TIMESTAMP NOT NULL
);

COMMIT;
