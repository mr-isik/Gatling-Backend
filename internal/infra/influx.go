package infra

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func NewInflux(cfg InfluxConfig) (influxdb2.Client, api.WriteAPIBlocking, api.QueryAPI, error) {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)

	ok, err := client.Health(context.Background())
	if err != nil || ok.Status != "pass" {
		return nil, nil, nil, fmt.Errorf("error checking influxdb health: %v", err)
	}

	writeAPI := client.WriteAPIBlocking(cfg.Org, cfg.Bucket)
	queryAPI := client.QueryAPI(cfg.Org)

	return client, writeAPI, queryAPI, nil
}
