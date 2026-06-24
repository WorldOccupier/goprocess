# goprocess

Processes website content so they can be ranked.

goprocess retrieves unprocessed HTML pages from a crawl database, parses the content to extract and count term frequencies per URL, and stores the results in a process database for downstream ranking.

## Prerequisites

- Go 1.26+
- Docker & Docker Compose (for containerized deployment)
- PostgreSQL 16 (crawl database)
- PostgreSQL 16 (process database)

## Build

### Local

```sh
make localbuild
```

### Docker

```sh
make up
# or
docker compose up -d --build
```

## Clean

Remove Docker volumes and containers:

```sh
make clean
```

Stop containers without removing volumes:

```sh
make down
```

## Usage

Configure via environment variables:

- `CRAWL_DB_URL` — crawl database connection string (default: `postgres://user:pass@host.docker.internal:5432/mydb`)
- `PROCESS_DB_URL` — process database connection string (default: `postgres://user:pass@postgres:5432/goprocess`)

The service processes pages in batches, parses HTML text (excluding script/style), counts term occurrences, and persists term counts to the process database.
