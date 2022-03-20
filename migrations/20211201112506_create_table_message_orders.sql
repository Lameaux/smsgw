-- +goose Up
-- +goose StatementBegin
CREATE TABLE "message_orders" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "merchant_id" uuid NOT NULL,
    "sender" varchar(15),
    "body" text  NOT NULL,
    "client_transaction_id" varchar(36),
    "notification_url" text,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL
);

CREATE INDEX "message_orders_merchant_id_and_created_at" ON "message_orders" ( "merchant_id", "created_at" );
CREATE UNIQUE INDEX "message_orders_client_transaction_id" ON "message_orders" ( "merchant_id", "client_transaction_id" )
WHERE "client_transaction_id" IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE message_orders;
-- +goose StatementEnd
