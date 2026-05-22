package influx

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type metricRepository struct {
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
	bucket   string
}

func NewMetricRepository(writeAPI api.WriteAPIBlocking, queryAPI api.QueryAPI, bucket string) repository.MetricRepository {
	return &metricRepository{
		writeAPI: writeAPI,
		queryAPI: queryAPI,
		bucket:   bucket,
	}
}

func (r *metricRepository) Write(ctx context.Context, metrics []domain.Metric) error {
	var points []*write.Point
	for _, m := range metrics {
		p := influxdb2.NewPointWithMeasurement("http_request").
			AddTag("run_id", m.RunID).
			AddTag("step_name", m.StepName).
			AddTag("status_code", fmt.Sprintf("%d", m.StatusCode)).
			AddField("latency_ms", m.Latency).
			AddField("bytes_sent", m.BytesSent).
			AddField("bytes_received", m.BytesRecv).
			AddField("success", m.Success).
			SetTime(m.Timestamp)
		points = append(points, p)
	}
	return r.writeAPI.WritePoint(ctx, points...)
}

func (r *metricRepository) QueryTimeSeries(ctx context.Context, runID string, from, to time.Time) ([]domain.AggregatedMetric, error) {
	// Simple stub for now. Real Flux query will be added later.
	return nil, nil
}

func (r *metricRepository) QueryLatest(ctx context.Context, runID string) (*domain.AggregatedMetric, error) {
	// Simple stub for now
	return nil, nil
}
