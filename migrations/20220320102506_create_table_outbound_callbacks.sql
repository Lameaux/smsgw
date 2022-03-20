-- +goose Up
-- +goose StatementBegin
CREATE TABLE "outbound_callbacks" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "merchant_id" uuid NOT NULL,
    "notification_url" text NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL
);

CREATE UNIQUE INDEX "outbound_callbacks_merchant_id" ON "outbound_callbacks" ( "merchant_id" );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outbound_callbacks;
-- +goose StatementEnd
