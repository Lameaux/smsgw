-- +goose Up
-- +goose StatementBegin
CREATE TABLE "delivery_notifications" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "message_type" char(1) NOT NULL,
    "message_id" uuid NOT NULL,
    "status" char(1) NOT NULL,
    "last_response" text,
    "next_attempt_at" timestamp NOT NULL,
    "attempt_counter" int NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp  NOT NULL
);

CREATE INDEX "delivery_outbound_next_attempt_at" ON "delivery_notifications"
( "message_type", "status", "next_attempt_at" )
WHERE "message_type" = 'o' AND "status" = 'n';
CREATE INDEX "delivery_inbound_next_attempt_at" ON "delivery_notifications"
( "message_type", "status", "next_attempt_at" )
WHERE "message_type" = 'i' AND "status" = 'n';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE delivery_notifications;
-- +goose StatementEnd
