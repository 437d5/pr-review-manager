CREATE TABLE pull_requests_reviewers (
    id BIGSETIALL PRIMARY KEY,
    pull_request_id VARCHAR(255) NOT NULL,
    reviewer_id VARCHAR(255) NOT NULL,
    asigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_reviewers_pull_request
        FOREIGN KEY (pull_request_id)
        REFERENCES pull_requests(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_reviewers_user
        FOREIGN KEY (reviewer_id)
        REFERENCES users(id)

    CONSTRAINT unique_reviewers UNIQUE (pull_request_id, reviewer_id)
);

CREATE INDEX idx_pr_reviewer_pull_request_id ON pull_requests_reviewers (pull_request_id);
CREATE INDEX idx_pr_reviewer_reviewer_id ON pull_requests_reviewers (reviewer_id);