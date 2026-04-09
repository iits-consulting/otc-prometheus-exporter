package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
	"github.com/iits-consulting/otc-prometheus-exporter/provider"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// nopLogger is a no-op logger for tests.
type nopLogger struct{}

func (n *nopLogger) Info(_ string, _ ...interface{})              {}
func (n *nopLogger) Debug(_ string, _ ...interface{})             {}
func (n *nopLogger) Warn(_ string, _ ...interface{})              {}
func (n *nopLogger) Error(_ string, _ ...interface{})             {}
func (n *nopLogger) Panic(_ string, _ ...interface{})             {}
func (n *nopLogger) Sync() error                                  { return nil }
func (n *nopLogger) WithFields(_ ...interface{}) internal.ILogger { return n }

var _ internal.ILogger = (*nopLogger)(nil)

// mockProvider implements provider.MetricProvider for testing.
type mockProvider struct {
	namespace string
	families  []*dto.MetricFamily
	err       error
}

func (m *mockProvider) Namespace() string { return m.namespace }
func (m *mockProvider) Collect(_ context.Context, _ *otcclient.Client) ([]*dto.MetricFamily, error) {
	return m.families, m.err
}

func TestMetricsHandlerNoNamespaceReturnsExporterMetrics(t *testing.T) {
	registry := provider.NewRegistry()
	handler := metricsHandler(registry, nil, &nopLogger{}, newExporterRegistry())

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "go_goroutines") {
		t.Fatalf("expected Go runtime metrics, got: %s", body[:200])
	}
}

func TestMetricsHandlerUnknownNamespaceReturnsOK(t *testing.T) {
	registry := provider.NewRegistry()
	handler := metricsHandler(registry, nil, &nopLogger{}, prometheus.NewRegistry())

	req := httptest.NewRequest(http.MethodGet, "/metrics?namespace=UNKNOWN", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestMetricsHandlerSuccess(t *testing.T) {
	registry := provider.NewRegistry()

	metricName := "sys_ecs_cpu_util"
	gaugeValue := 42.5
	families := []*dto.MetricFamily{
		{
			Name: &metricName,
			Type: dto.MetricType_GAUGE.Enum(),
			Metric: []*dto.Metric{
				{
					Gauge: &dto.Gauge{Value: &gaugeValue},
				},
			},
		},
	}

	mock := &mockProvider{
		namespace: "SYS.ECS",
		families:  families,
	}
	registry.Register(mock)

	handler := metricsHandler(registry, nil, &nopLogger{}, prometheus.NewRegistry())

	req := httptest.NewRequest(http.MethodGet, "/metrics?namespace=SYS.ECS", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), metricName) {
		t.Fatalf("expected body to contain %q, got: %s", metricName, rec.Body.String())
	}
}

func TestMetricsHandlerProviderError(t *testing.T) {
	registry := provider.NewRegistry()

	mock := &mockProvider{
		namespace: "SYS.FAIL",
		err:       errors.New("boom"),
	}
	registry.Register(mock)

	handler := metricsHandler(registry, nil, &nopLogger{}, prometheus.NewRegistry())

	req := httptest.NewRequest(http.MethodGet, "/metrics?namespace=SYS.FAIL", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("expected body 'ok', got: %s", rec.Body.String())
	}
}
