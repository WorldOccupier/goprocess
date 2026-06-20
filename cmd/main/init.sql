CREATE TABLE IF NOT EXISTS t_url_term_count (
    term VARCHAR(1000) NOT NULL,
    termCount INT NOT NULL,
    url VARCHAR(2048) NOT NULL,
    UNIQUE(term, url)
);
