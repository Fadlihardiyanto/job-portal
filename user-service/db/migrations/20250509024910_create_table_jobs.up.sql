CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    jobs_title VARCHAR(255),
    company_id BIGINT REFERENCES company(id) ON DELETE CASCADE,
    location INT(11),
    workspace_type VARCHAR(255),
    min_salary VARCHAR(255),
    max_salary VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
