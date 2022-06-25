-- +goose Up
-- +goose StatementBegin
CREATE TABLE "outbound_messages" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "merchant_id" uuid NOT NULL,
    "message_group_id" uuid NOT NULL,
    "status" char(1) NOT NULL,
    "msisdn" bigint NOT NULL,
    "provider_id" varchar(36),
    "provider_message_id" varchar(255),
    "provider_response" text,
    "next_attempt_at" timestamp NOT NULL,
    "attempt_counter" int NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp  NOT NULL
);

CREATE INDEX "outbound_merchant_id_and_message_group_id" ON "outbound_messages" ( "merchant_id", "message_group_id" );
CREATE INDEX "outbound_merchant_id_and_created_at" ON "outbound_messages" ( "merchant_id", "created_at");
CREATE INDEX "outbound_merchant_id_and_status" ON "outbound_messages" ( "merchant_id", "status" );
CREATE INDEX "outbound_merchant_id_and_msisdn" ON "outbound_messages" ( "merchant_id", "msisdn" );
CREATE INDEX "outbound_status_and_next_attempt_at" ON "outbound_messages" ( "status", "next_attempt_at" )
    WHERE "status" = 'n';
CREATE UNIQUE INDEX "outbound_provider_message_id" ON "outbound_messages" ( "provider_id", "provider_message_id" )
    WHERE "provider_message_id" IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outbound_messages;
-- +goose StatementEnd
