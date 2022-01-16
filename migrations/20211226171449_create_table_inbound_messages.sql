-- +goose Up
-- +goose StatementBegin
CREATE TABLE "inbound_messages" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "shortcode" varchar(15) NOT NULL,
    "status" varchar(255) NOT NULL,
    "msisdn" varchar(15) NOT NULL,
    "body" text  NOT NULL,
    "provider_id" varchar(36) NOT NULL,
    "provider_message_id" varchar(255) NOT NULL,
    "next_attempt_at" timestamp NOT NULL,
    "attempt_counter" int NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp  NOT NULL
);

CREATE INDEX "inbound_shortcode_and_created_at" ON "inbound_messages" ( "shortcode", "created_at" );
CREATE INDEX "inbound_shortcode_and_msisdn" ON "inbound_messages" ( "shortcode", "msisdn" );
CREATE INDEX "inbound_shortcode_and_status" ON "inbound_messages" ( "shortcode", "status" );
CREATE INDEX "inbound_status_and_next_attempt_at" ON "inbound_messages" ( "status", "next_attempt_at" ) WHERE "status" = 'new';
CREATE UNIQUE INDEX "inbound_provider_message_id" ON "inbound_messages" ( "provider_id", "provider_message_id" );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE inbound_messages;
-- +goose StatementEnd
