package provider

import "time"

// Config holds tunable parameters for metric collection.
// Set via CLI flags / environment variables.
var Config = struct {
	// CESBatchSize is the maximum number of metrics per CES batch API request.
	// The OTC SDK documents a limit of 10, but the API accepts up to 500 in practice.
	CESBatchSize int
	// AOMBatchSize is the maximum number of metrics per AOM data API request.
	// The documented limit is 20.
	AOMBatchSize int
	// AOMConcurrency is the maximum number of concurrent AOM API calls per scrape.
	AOMConcurrency int
	// CESLookback is how far back the CES batch query looks for datapoints.
	// CES updates metrics every ~5 minutes; 10 minutes ensures at least one
	// datapoint is captured even with slight delays.
	CESLookback time.Duration
	// CESObsStorageMetricLookback is how far back the query for OBS storage
	// metrics (capacity_*, object_num_*) looks for datapoints. It does not
	// apply to OBS request metrics (download_bytes etc.), which use
	// CESLookback. OBS storage metrics have a 30-minute monitoring
	// period, so values below 30 minutes lead to gaps in the data.
	CESObsStorageMetricLookback time.Duration
	// CollectTimeout is the maximum time a single Collect() call may run.
	// Set just below the Prometheus scrape timeout (typically 60s) so the
	// exporter can return partial results instead of timing out silently.
	CollectTimeout time.Duration
}{
	CESBatchSize:                500,
	AOMBatchSize:                20,
	AOMConcurrency:              5,
	CESLookback:                 10 * time.Minute,
	CESObsStorageMetricLookback: 60 * time.Minute,
	CollectTimeout:              55 * time.Second,
}
