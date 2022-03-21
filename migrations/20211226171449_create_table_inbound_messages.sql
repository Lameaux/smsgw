-- +goose Up
-- +goose StatementBegin
CREATE TABLE "inbound_messages" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "merchant_id" uuid NOT NULL,
    "shortcode" varchar(15) NOT NULL,
    "status" char(1) NOT NULL,
    "msisdn" bigint NOT NULL,
    "body" text  NOT NULL,
    "provider_id" varchar(36) NOT NULL,
    "provider_message_id" varchar(255) NOT NULL,
    "next_attempt_at" timestamp NOT NULL,
    "attempt_counter" int NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp  NOT NULL
);

CREATE INDEX "inbound_merchant_id_and_created_at" ON "inbound_messages" ( "merchant_id", "created_at" );
CREATE INDEX "inbound_merchant_id_and_msisdn" ON "inbound_messages" ( "merchant_id", "msisdn" );
CREATE INDEX "inbound_merchant_id_and_status" ON "inbound_messages" ( "merchant_id", "status" );
CREATE INDEX "inbound_merchant_id_and_shortcode" ON "inbound_messages" ( "merchant_id", "shortcode" );
CREATE INDEX "inbound_status_and_next_attempt_at" ON "inbound_messages" ( "status", "next_attempt_at" ) 
WHERE "status" = 'n';
CREATE UNIQUE INDEX "inbound_provider_message_id" ON "inbound_messages" ( "provider_id", "provider_message_id" );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE inbound_messages;
-- +goose StatementEnd
