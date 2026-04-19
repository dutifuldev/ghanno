CREATE TABLE IF NOT EXISTS repository_access_grants (
    id BIGSERIAL PRIMARY KEY,
    github_repository_id BIGINT NOT NULL,
    github_user_id BIGINT NOT NULL,
    github_login TEXT NOT NULL,
    role TEXT NOT NULL,
    granted_by_github_user_id BIGINT NOT NULL,
    granted_by_github_login TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_repository_access_grants_repo_user
    ON repository_access_grants(github_repository_id, github_user_id);

CREATE INDEX IF NOT EXISTS idx_repository_access_grants_repo
    ON repository_access_grants(github_repository_id);

CREATE INDEX IF NOT EXISTS idx_repository_access_grants_user
    ON repository_access_grants(github_user_id);
