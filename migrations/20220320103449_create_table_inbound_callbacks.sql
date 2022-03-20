-- +goose Up
-- +goose StatementBegin
CREATE TABLE "inbound_callbacks" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "shortcode" varchar(15) NOT NULL,
    "notification_url" text NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp  NOT NULL
);

CREATE UNIQUE INDEX "inbound_callbacks_shortcode" ON "inbound_callbacks" ( "shortcode" );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE inbound_callbacks;
-- +goose StatementEnd
