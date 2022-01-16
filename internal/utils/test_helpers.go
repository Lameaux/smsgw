package utils

import (
	"context"
	"fmt"
	"time"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/logger"
)

var tables = []string{
	"message_orders",
	"outbound_messages",
	"inbound_messages",
	"inbound_notifications",
}

func CleanupDatabase(db db.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s", table))
		if err != nil {
			logger.Fatal(err)
		}
	}
}
