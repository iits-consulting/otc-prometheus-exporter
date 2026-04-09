package otcclient

import (
	"strings"
	"sync"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
)

// clientCache caches service clients so endpoint discovery only happens once per service.
type clientCache struct {
	mu    sync.Mutex
	items map[string]clientCacheEntry
}

type clientCacheEntry struct {
	client *golangsdk.ServiceClient
}

func (cc *clientCache) getOrCreate(name string, factory func() (*golangsdk.ServiceClient, error)) (*golangsdk.ServiceClient, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if r, ok := cc.items[name]; ok {
		return r.client, nil
	}
	client, err := factory()
	if err != nil {
		return nil, err
	}

	cc.items[name] = clientCacheEntry{client: client}
	return client, nil
}

// ---------------------------------------------------------------------------
// Service client factories — each caches its result via clientCache.
// ---------------------------------------------------------------------------

func (c *Client) CESClient() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("ces", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewCESClient(c.provider, c.endpointOpts())
	})
}

// RegionCESClient returns a CES client scoped to the region-level project.
// This is needed for global services like OBS whose CES metrics are only
// available under the region project (e.g. eu-de), not under a specific project.
// Falls back to the regular CES client if no region project was discovered.
func (c *Client) RegionCESClient() (*golangsdk.ServiceClient, error) {
	if c.regionProvider != nil {
		return openstack.NewCESClient(c.regionProvider, c.endpointOpts())
	}
	c.Logger.Warn("no region project available, falling back to project-scoped CES; global services like OBS may not return metrics")
	return c.CESClient()
}

func (c *Client) ComputeV2() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("compute", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewComputeV2(c.provider, c.endpointOpts())
	})
}

func (c *Client) RDSV3() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("rds", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewRDSV3(c.provider, c.endpointOpts())
	})
}

// DMSV2 returns a DMS v2 service client. Note: the previous implementation used v1;
// v2 provides topic, consumer group, and broker-level APIs needed for rich metrics.
func (c *Client) DMSV2() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("dms", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewDMSServiceV2(c.provider, c.endpointOpts())
	})
}

func (c *Client) ELBV3() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("elb", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewELBV3(c.provider, c.endpointOpts())
	})
}

func (c *Client) NatV2() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("nat", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewNatV2(c.provider, c.endpointOpts())
	})
}

func (c *Client) DDSV3() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("dds", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewDDSServiceV3(c.provider, c.endpointOpts())
	})
}

func (c *Client) DCSV1() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("dcs", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewDCSServiceV1(c.provider, c.endpointOpts())
	})
}

func (c *Client) NetworkV1() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("network", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewNetworkV1(c.provider, c.endpointOpts())
	})
}

func (c *Client) BlockStorageV3() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("block", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewBlockStorageV3(c.provider, c.endpointOpts())
	})
}

func (c *Client) CBRV3() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("cbr", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewCBRService(c.provider, c.endpointOpts())
	})
}

func (c *Client) AutoScalingV1() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("as", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewAutoScalingV1(c.provider, c.endpointOpts())
	})
}

func (c *Client) WAFV1() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("waf", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewWAFV1(c.provider, c.endpointOpts())
	})
}

func (c *Client) DWSV1() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("dws", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewDWSV1(c.provider, c.endpointOpts())
	})
}

func (c *Client) CSSV1() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("css", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewCSSService(c.provider, c.endpointOpts())
	})
}

func (c *Client) SFSV2() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("sfs", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewSharedFileSystemV2(c.provider, c.endpointOpts())
	})
}

func (c *Client) SFSTurboV1() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("sfsTurbo", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewSharedFileSystemTurboV1(c.provider, c.endpointOpts())
	})
}

func (c *Client) DCaaSV2() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("dcaas", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewDCaaSV2(c.provider, c.endpointOpts())
	})
}

func (c *Client) OBS() (*golangsdk.ServiceClient, error) {
	return c.cache.getOrCreate("obs", func() (*golangsdk.ServiceClient, error) {
		return openstack.NewOBSService(c.provider, c.endpointOpts())
	})
}

// AOMV1 returns an AOM v1 service client. NOT cached because it mutates the
// endpoint derived from CES — caching would share the mutated ServiceClient.
func (c *Client) AOMV1() (*golangsdk.ServiceClient, error) {
	sc, err := openstack.NewCESClient(c.provider, c.endpointOpts())
	if err != nil {
		return nil, err
	}
	// CES endpoint: https://ces.{region}.otc.t-systems.com/V1.0/{project_id}/
	// AOM endpoint: https://aom.{region}.otc.t-systems.com/v1/{project_id}/
	sc.Endpoint = strings.Replace(sc.Endpoint, "ces.", "aom.", 1)
	sc.Endpoint = strings.Replace(sc.Endpoint, "/V1.0/", "/v1/", 1)
	sc.ResourceBase = sc.Endpoint
	return sc, nil
}
