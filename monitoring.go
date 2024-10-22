package main

import (
	"context"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredres "google.golang.org/genproto/googleapis/api/monitoredres"
)

const (
	metricType        = "custom.googleapis.com/conversation"
	revisionName      = "myherodotus-00004-jcg" // TODO(telpirion): Get this dynamically
	configurationName = "test"                  // TODO(telpirion): Figure out what this means?
)

// writeTimeSeriesValue writes a value for the custom metric created
func writeTimeSeriesValue(projectID, label string) {
	configurationName := os.Getenv("CONFIGURATION_NAME")

	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		LogError(err.Error())
		return
	}
	defer c.Close()

	now := &timestamp.Timestamp{
		Seconds: time.Now().Unix(),
	}
	req := &monitoringpb.CreateTimeSeriesRequest{
		Name: "projects/" + projectID,
		TimeSeries: []*monitoringpb.TimeSeries{{
			Metric: &metricpb.Metric{
				Type: metricType,
				Labels: map[string]string{
					"conversation": label,
				},
			},
			Resource: &monitoredres.MonitoredResource{
				Type: "cloud_run_revision",
				Labels: map[string]string{
					"project_id":         projectID,
					"service_name":       "myherodotus",
					"revision_name":      revisionName,
					"location":           "us-west1",
					"configuration_name": configurationName,
				},
			},
			Points: []*monitoringpb.Point{{
				Interval: &monitoringpb.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_Int64Value{
						Int64Value: 1,
					},
				},
			}},
		}},
	}

	err = c.CreateTimeSeries(ctx, req)
	if err != nil {
		LogError(err.Error())
	}
}
