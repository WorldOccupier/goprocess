CREATE TABLE IF NOT EXISTS t_url_term_count (
    term VARCHAR(1000) NOT NULL,
    termCount INT NOT NULL,
    url VARCHAR(2048) NOT NULL,
    UNIQUE(term, url)
);

CREATE INDEX IF NOT EXISTS idx_url_term_count_lookup ON t_url_term_count (term, url, termCount);

CREATE TABLE IF NOT EXISTS t_metadata (
    key VARCHAR(100) PRIMARY KEY,
    value INT NOT NULL DEFAULT 0
);
INSERT INTO t_metadata (key, value)
VALUES ('total_urls', (SELECT COUNT(DISTINCT url) FROM t_url_term_count))
ON CONFLICT (key) DO NOTHING;
