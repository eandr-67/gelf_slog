package main

import (
	"fmt"
	"log"
	"time"

	"log/slog"

	"github.com/eandr-67/gelf"
	"github.com/eandr-67/gelf_slog"
)

func main() {
	// docker-compose up -d
	// or
	// ncat -l 12201 -u
	gelfWriter, err := gelf.NewUDPWriter("localhost:12201")
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}

	gelfWriter.CompressionType = gelf.CompressNone // for debugging only

	logger := slog.New(gelf_slog.Option{Level: slog.LevelDebug, Writer: gelfWriter}.NewGraylogHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now().AddDate(0, 0, -1)),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")

	time.Sleep(1 * time.Second)
}
