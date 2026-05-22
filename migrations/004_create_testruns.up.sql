CREATE TABLE testruns (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scenario_id UUID NOT NULL REFERENCES scenarios(id),
    project_id  UUID NOT NULL REFERENCES projects(id),
    status      VARCHAR(50) NOT NULL DEFAULT 'pending',
    config      JSONB NOT NULL DEFAULT '{}',
    started_by  UUID REFERENCES users(id),
    started_at  TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    error_msg   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_testruns_scenario_id ON testruns(scenario_id);
CREATE INDEX idx_testruns_project_id ON testruns(project_id);
CREATE INDEX idx_testruns_status ON testruns(status);
