CREATE TABLE reports (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id      UUID NOT NULL REFERENCES testruns(id) ON DELETE CASCADE,
    summary     JSONB NOT NULL DEFAULT '{}',
    ai_summary  TEXT,
    ai_recommendations JSONB DEFAULT '[]',
    anomalies   JSONB DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_reports_run_id ON reports(run_id);
