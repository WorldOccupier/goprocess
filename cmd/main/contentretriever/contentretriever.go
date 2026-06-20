package contentretriever

import (
	"context"
	"goprocess/logger"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ctx = context.Background()
	defaultDatabaseUrl = "postgres://user:pass@host.docker.internal:5432/mydb"
)

type UrlContent struct {
	Url     string
	Content string
}

func getDBUrl() string {
	databaseUrl := os.Getenv("CRAWL_DB_URL")
	if databaseUrl == "" {
		databaseUrl = defaultDatabaseUrl
	}

	return databaseUrl
}

func rowToUrlContent(row pgx.CollectableRow) (UrlContent, error) {
	var u UrlContent
	err := row.Scan(&u.Url, &u.Content)
	return u, err
}

func Get(count int) []UrlContent {
	dbUrl := getDBUrl()
	connection, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		logger.Log.Error("Could not get connection pool", "err", err)
		return nil
	}
	defer connection.Close()

	rows, err := connection.Query(ctx, "SELECT url, content FROM t_web_page_details LIMIT $1", count)
	if err != nil {
		logger.Log.Error("Error retrieving url content", "err", err)
		return nil
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, rowToUrlContent)
	if err != nil {
		logger.Log.Error("Error collecting rows", "err", err)
		return nil
	}
	logger.Log.Info("Rows count: " + strconv.Itoa(len(results)))

	return results
}
