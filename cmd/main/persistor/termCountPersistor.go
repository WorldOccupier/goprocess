package persistor

import (
	"context"
	"goprocess/logger"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ctx                = context.Background()
	defaultDatabaseUrl = "postgres://user:pass@postgres:5432/goprocess"
)

type UrlTermCount struct {
	Url       string
	TermCount map[string]int
}

func getDBUrl() string {
	databaseUrl := os.Getenv("PROCESS_DB_URL")
	if databaseUrl == "" {
		databaseUrl = defaultDatabaseUrl
	}

	return databaseUrl
}

func SaveUrlTermCounts(urlTermCounts []UrlTermCount) {
	dbUrl := getDBUrl()
	connection, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		logger.Log.Error("Unable to get DB connection", "error", err)
		return
	}
	defer connection.Close()
	batch := &pgx.Batch{}
	for _, utc := range urlTermCounts {
		for term, count := range utc.TermCount {
			batch.Queue(
				`INSERT INTO t_url_term_count (term, termCount, url)
                 VALUES ($1, $2, $3)
                 ON CONFLICT (term, url) DO UPDATE SET termCount = $2`,
				term, count, utc.Url,
			)
		}
	}

	br := connection.SendBatch(ctx, batch)
	br.Close()
}
