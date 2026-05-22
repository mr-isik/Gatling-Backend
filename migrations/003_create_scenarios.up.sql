CREATE TABLE scenarios (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id    UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    tags          TEXT[],
    steps         JSONB NOT NULL DEFAULT '[]',
    is_ai_generated BOOLEAN NOT NULL DEFAULT FALSE,
    created_by    UUID REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_scenarios_project_id ON scenarios(project_id);
