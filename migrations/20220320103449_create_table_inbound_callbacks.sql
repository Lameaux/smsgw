-- +goose Up
-- +goose StatementBegin
CREATE TABLE "inbound_callbacks" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "merchant_id" uuid NOT NULL,
    "shortcode" varchar(15) NOT NULL,
    "url" text NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp  NOT NULL
);

CREATE UNIQUE INDEX "inbound_callbacks_merchant_id_shortcode" ON "inbound_callbacks" ( "merchant_id", "shortcode" );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE inbound_callbacks;
-- +goose StatementEnd
