package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
	"github.com/iits-consulting/otc-prometheus-exporter/provider"
)

// metricsHandler returns an HTTP handler that serves Prometheus metrics.
// Without a namespace param it returns the exporter's own metrics (Go runtime +
// scrape instrumentation). With ?namespace=X it scrapes the given OTC namespace.
func metricsHandler(registry *provider.Registry, client *otcclient.Client, logger internal.ILogger, exporterReg *prometheus.Registry) http.HandlerFunc {
	exporterHandler := promhttp.HandlerFor(exporterReg, promhttp.HandlerOpts{})

	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			exporterHandler.ServeHTTP(w, r)
			return
		}

		prov := registry.GetOrFallback(namespace)

		scrapeLog := logger.WithFields("trace", newTraceID(), "namespace", namespace)

		enrich := r.URL.Query().Get("enrich") != "false"
		ctx := provider.WithEnrich(r.Context(), enrich)
		ctx, cancel := context.WithTimeout(ctx, provider.Config.CollectTimeout)
		defer cancel()

		scrapeClient := client
		if client != nil {
			scrapeClient = client.WithContext(ctx)
			scrapeClient.Logger = scrapeLog
		}

		start := time.Now()
		families, err := prov.Collect(ctx, scrapeClient)
		duration := time.Since(start)

		if err != nil {
			scrapeDuration.WithLabelValues(namespace, "false").Observe(duration.Seconds())
			var notFound *provider.ErrNamespaceNotFound
			if errors.As(err, &notFound) {
				scrapeLog.Warn("unknown namespace",
					"duration", duration.String(),
					"error", err.Error())
				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprintf(w, "unknown namespace %q\n", namespace)
				return
			}
			scrapeLog.Error("collect failed",
				"duration", duration.String(),
				"error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "collect error: %v\n", err)
			return
		}

		provider.EnrichWithHelp(families)

		metricCount := 0
		for _, fam := range families {
			metricCount += len(fam.GetMetric())
		}

		scrapeDuration.WithLabelValues(namespace, "true").Observe(duration.Seconds())
		scrapeMetrics.WithLabelValues(namespace).Set(float64(metricCount))
		scrapeFamilies.WithLabelValues(namespace).Set(float64(len(families)))

		// Check if Prometheus already disconnected (scrape timeout) while
		// we were collecting. If so, log a warning — the data we collected
		// will be discarded since the connection is gone.
		if r.Context().Err() != nil {
			scrapeLog.Warn("client disconnected during collect, response discarded",
				"duration", duration.String(),
				"families", len(families),
				"metrics", metricCount)
			return
		}

		scrapeLog.Debug("collect completed",
			"duration", duration.String(),
			"families", len(families),
			"metrics", metricCount,
			"enrich", enrich)

		contentType := expfmt.Negotiate(r.Header)
		w.Header().Set("Content-Type", string(contentType))

		enc := expfmt.NewEncoder(w, contentType)
		for _, fam := range families {
			if len(fam.GetMetric()) == 0 {
				continue
			}
			if err := enc.Encode(fam); err != nil {
				scrapeLog.Warn("client disconnected during response",
					"duration", duration.String())
				return
			}
		}
	}
}

// newTraceID generates a UUID v4 for trace correlation.
func newTraceID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// healthzHandler returns a simple 200 "ok" response for liveness probes.
func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "ok")
}
