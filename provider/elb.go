package provider

import (
	"context"

	elbListeners "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/listeners"
	elbLB "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// ELBProvider collects CES metrics for the OTC Elastic Load Balancer service,
// enriches them with load balancer and listener names, and reports LB status.
type ELBProvider struct{}

func (p *ELBProvider) Namespace() string { return "SYS.ELB" }

func (p *ELBProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.ELB", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		elbClient, err := client.ELBV3()
		if err != nil {
			return nil, err
		}

		lbPages, err := elbLB.List(elbClient, elbLB.ListOpts{}).AllPages()
		if err != nil {
			return nil, err
		}
		lbs, err := elbLB.ExtractLoadbalancers(lbPages)
		if err != nil {
			return nil, err
		}

		lbNameMap := make(map[string]string, len(lbs))
		for _, lb := range lbs {
			lbNameMap[lb.ID] = lb.Name
		}

		// Listeners are optional — continue without them on failure.
		listenerNameMap := make(map[string]string)
		listenerPages, err := elbListeners.List(elbClient, elbListeners.ListOpts{}).AllPages()
		if err != nil {
			client.Logger.Warn("listener enrichment failed", "namespace", "SYS.ELB", "error", err.Error())
		} else {
			listeners, err := elbListeners.ExtractListeners(listenerPages)
			if err != nil {
				client.Logger.Warn("listener extract failed", "namespace", "SYS.ELB", "error", err.Error())
			} else {
				for _, l := range listeners {
					listenerNameMap[l.ID] = l.Name
				}
			}
		}

		client.Logger.Debug("enrichment completed", "namespace", "SYS.ELB",
			"loadbalancers", len(lbs), "listeners", len(listenerNameMap))

		return &EnrichResult{
			// Merged map for cache: on fallback, EnrichWithNames uses resource_id
			// which matches LB IDs. Listener composite names are lost on fallback,
			// but a LB name is better than no name for this rare error path.
			NameMap:       buildELBMergedNameMap(lbNameMap, listenerNameMap),
			ExtraFamilies: convertELBToMetrics(lbs),
			Enrich: func(families []*dto.MetricFamily) {
				enrichELBNames(families, lbNameMap, listenerNameMap)
			},
		}, nil
	})
}

// enrichELBNames sets resource_name based on LB and listener name maps.
// For metrics with a lbaas_listener_id label, the name is "lb/listener".
// For LB-level metrics (no listener ID), just the LB name.
func enrichELBNames(families []*dto.MetricFamily, lbNames, listenerNames map[string]string) {
	for _, fam := range families {
		for _, m := range fam.GetMetric() {
			var lbID, listenerID string
			var nameLabel *dto.LabelPair
			for _, lp := range m.GetLabel() {
				switch lp.GetName() {
				case "lbaas_instance_id":
					lbID = lp.GetValue()
				case "lbaas_listener_id":
					listenerID = lp.GetValue()
				case "resource_name":
					nameLabel = lp
				}
			}
			if nameLabel == nil {
				continue
			}

			lbName := lbNames[lbID]
			if lbName == "" {
				continue
			}

			if listenerID != "" {
				listenerName := listenerNames[listenerID]
				if listenerName == "" {
					if len(listenerID) > 8 {
						listenerName = listenerID[:8] // short ID as fallback
					} else {
						listenerName = listenerID
					}
				}
				composite := lbName + "/" + listenerName
				nameLabel.Value = &composite
			} else {
				nameLabel.Value = &lbName
			}
		}
	}
}

// buildELBMergedNameMap creates a single map containing both LB and listener ID-to-name mappings.
// Used for caching: a single cache entry covers both enrichment sources.
func buildELBMergedNameMap(lbNames, listenerNames map[string]string) map[string]string {
	merged := make(map[string]string, len(lbNames)+len(listenerNames))
	for id, name := range lbNames {
		merged[id] = name
	}
	for id, name := range listenerNames {
		merged[id] = name
	}
	return merged
}

// convertELBToMetrics creates a MetricFamily "elb_loadbalancer_status" with a
// gauge metric per load balancer. The value is 0.0 for ONLINE LBs, 1.0 otherwise
// (OTC convention: 0=normal, 1=abnormal).
func convertELBToMetrics(lbs []elbLB.LoadBalancer) []*dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, len(lbs))
	for _, lb := range lbs {
		value := 1.0
		if lb.OperatingStatus == "ONLINE" {
			value = 0.0
		}
		metrics = append(metrics, NewGaugeMetric(value, map[string]string{
			"resource_id":         lb.ID,
			"resource_name":       lb.Name,
			"provisioning_status": lb.ProvisioningStatus,
			"operating_status":    lb.OperatingStatus,
		}))
	}
	return []*dto.MetricFamily{NewGaugeMetricFamily("elb_loadbalancer_status", metrics)}
}
