package utils

import (
	"context"
	"fmt"
	"time"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/logger"
)

var tables = []string{ //nolint:gochecknoglobals
	"message_orders",
	"outbound_messages",
	"inbound_messages",
	"inbound_notifications",
}

const execTimeout = 2 * time.Second

func CleanupDatabase(db db.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s", table))
		if err != nil {
			logger.Fatal(err)
		}
	}
}
