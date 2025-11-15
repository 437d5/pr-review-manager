CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL REFERENCES users(id),
    author_id VARCHAR(255) NOT NULL,
    status  pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    merged_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT check_merged_at_valid CHECK (
        (status = 'MERGED' AND merged_at IS NOT NULL) OR 
        (status = 'OPEN' AND merged_at IS NULL)
    )
);

CREATE INDEX idx_pull_requests_author_id ON pull_requests(author_id);
CREATE INDEX idx_pull_requests_status ON pull_requests(status);