-- +goose Up
-- +goose StatementBegin
CREATE TABLE "merchants" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "email" text NOT NULL,
    "apikey1" char(36) NOT NULL,
    "apikey2" char(36) NOT NULL,
    "callback_url" text,
    "shortcode" varchar(15),
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    "deleted_at" timestamp
);

CREATE UNIQUE INDEX "merchants_email" ON "merchants" ( "email" )
    WHERE "deleted_at" IS NULL;

CREATE UNIQUE INDEX "merchants_shortcode" ON "merchants" ( "shortcode" )
    WHERE "shortcode" IS NOT NULL AND "deleted_at" IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE merchants;
-- +goose StatementEnd
