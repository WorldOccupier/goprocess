package metadata

import (
	"context"
	"goprocess/logger"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ctx                = context.Background()
	defaultDatabaseUrl = "postgres://user:pass@postgres:5432/goprocess"
	UpdateInterval     = 1 * time.Minute
)

func getDBUrl() string {
	databaseUrl := os.Getenv("PROCESS_DB_URL")
	if databaseUrl == "" {
		databaseUrl = defaultDatabaseUrl
	}

	return databaseUrl
}

func startUpdater() {
	dbUrl := getDBUrl()
	connection, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		logger.Log.Error("Metadata updater: failed to connect", "error", err)
		return
	}
	defer connection.Close()

	ticker := time.NewTicker(UpdateInterval)
	defer ticker.Stop()

	for {
		_, err := connection.Exec(ctx, "UPDATE t_metadata SET value = (SELECT COUNT(DISTINCT url) FROM t_url_term_count) WHERE key = 'total_urls'")
		if err != nil {
			logger.Log.Error("Metadata updater: update failed", "err", err)
		} else {
			logger.Log.Info("Metadata updater: total_urls refreshed")
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func StartUpdater(ctx context.Context, wg *sync.WaitGroup) {
	wg.Go(startUpdater)
}
