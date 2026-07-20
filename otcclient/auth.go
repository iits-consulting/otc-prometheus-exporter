package otcclient

import (
	"errors"
	"fmt"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"
)

// Config holds the authentication and endpoint configuration for an OTC client.
type Config struct {
	Username        string
	Password        string
	AccessKey       string
	SecretKey       string
	ProjectID       string
	DomainName      string
	Region          string
	RequestTimeout  time.Duration // Overall HTTP request timeout (default: 10s)
	IdleConnTimeout time.Duration // How long idle connections stay in the pool (default: 90s)
}

func validateConfig(cfg Config) error {
	// Region must be eu-de, eu-nl, or eu-ch2 (Swiss OTC, experimental).
	validRegions := map[string]bool{"eu-de": true, "eu-nl": true, "eu-ch2": true}
	if !validRegions[cfg.Region] {
		return fmt.Errorf("otcclient: invalid region %q, must be \"eu-de\", \"eu-nl\", or \"eu-ch2\"", cfg.Region)
	}

	// ProjectID is always required.
	if cfg.ProjectID == "" {
		return fmt.Errorf("otcclient: ProjectID is required")
	}

	hasUserPass := cfg.Username != "" && cfg.Password != ""
	hasAKSK := cfg.AccessKey != "" && cfg.SecretKey != ""

	hasPartialUserPass := (cfg.Username != "") != (cfg.Password != "")
	hasPartialAKSK := (cfg.AccessKey != "") != (cfg.SecretKey != "")
	if hasPartialUserPass {
		return errors.New("both username and password are required")
	}
	if hasPartialAKSK {
		return errors.New("both access-key and secret-key are required")
	}

	if hasUserPass && hasAKSK {
		return fmt.Errorf("otcclient: provide either username/password or access-key/secret-key, not both")
	}
	if !hasUserPass && !hasAKSK {
		return fmt.Errorf("otcclient: must provide username/password or access-key/secret-key")
	}

	if hasUserPass {
		if cfg.DomainName == "" {
			return fmt.Errorf("otcclient: DomainName is required for user/password auth")
		}
	}

	return nil
}

func iamEndpoint(region string) string {
	// The Swiss OTC uses an irregular public IAM hostname.
	if region == "eu-ch2" {
		return "https://iam-pub.eu-ch2.sc.otc.t-systems.com/v3"
	}
	return fmt.Sprintf("https://iam.%s.otc.t-systems.com/v3", region)
}

func authOptions(cfg Config) golangsdk.AuthOptionsProvider {
	if cfg.AccessKey != "" {
		return golangsdk.AKSKAuthOptions{
			IdentityEndpoint: iamEndpoint(cfg.Region),
			AccessKey:        cfg.AccessKey,
			SecretKey:        cfg.SecretKey,
			ProjectId:        cfg.ProjectID,
		}
	}
	return golangsdk.AuthOptions{
		IdentityEndpoint: iamEndpoint(cfg.Region),
		Username:         cfg.Username,
		Password:         cfg.Password,
		DomainName:       cfg.DomainName,
		TenantID:         cfg.ProjectID,
		AllowReauth:      true,
	}
}

// SetRegionProjectID manually sets the region-level project ID and
// authenticates a region-scoped provider. Use this when auto-discovery
// is not possible (e.g. with AK/SK auth that can't list projects).
func (c *Client) SetRegionProjectID(id string) error {
	c.authenticateRegionProject(c.Region, id)
	if c.regionProvider == nil {
		return fmt.Errorf("failed to authenticate with region project ID %s", id)
	}
	return nil
}

// DiscoverRegionProjectID queries the IAM API to find the region-level
// project for the configured region. Global services like OBS require
// this project scope instead of a specific subproject.
func (c *Client) DiscoverRegionProjectID() error {
	identityClient, err := openstack.NewIdentityV3(c.provider, golangsdk.EndpointOpts{})
	if err != nil {
		return fmt.Errorf("failed to create identity client for region project discovery: %w", err)
	}

	pages, err := projects.List(identityClient, projects.ListOpts{}).AllPages()
	if err != nil {
		return fmt.Errorf("failed to list projects for region project discovery: %w", err)
	}

	allProjects, err := projects.ExtractProjects(pages)
	if err != nil {
		return fmt.Errorf("failed to extract projects for region project discovery: %w", err)
	}

	projectNames := make([]string, len(allProjects))
	for i, p := range allProjects {
		projectNames[i] = p.Name + "(" + p.ID + ")"
	}
	c.Logger.Debug("discovered projects", "projects", projectNames)

	// The region-level project is typically named after the region (e.g. "eu-de").
	for _, p := range allProjects {
		if p.Name == c.Region {
			c.authenticateRegionProject(p.Name, p.ID)
			return nil
		}
	}

	return fmt.Errorf("region-level project not found, global services like OBS may not work")
}

func (c *Client) authenticateRegionProject(name, id string) {
	regionCfg := c.cfg
	regionCfg.ProjectID = id
	regionOpts := authOptions(regionCfg)
	regionProvider, err := openstack.AuthenticatedClient(regionOpts)
	if err != nil {
		c.Logger.Warn("failed to authenticate with region project",
			"project_name", name, "project_id", id,
			"error", err.Error())
		return
	}

	// Reuse the tuned HTTP transport from the primary provider so the
	// region provider also benefits from increased connection pool limits.
	regionProvider.HTTPClient = c.provider.HTTPClient

	c.regionProvider = regionProvider
	c.RegionProjectID = id
	c.Logger.Info("authenticated with region-level project",
		"name", name, "id", id)
}
