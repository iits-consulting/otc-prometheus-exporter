package otcdoc

// DocumentationPage holds the parsed contents of a single OTC metrics RST page.
type DocumentationPage struct {
	Namespace string
	Metrics   []MetricDocumentationEntry
}

// MetricDocumentationEntry is a single metric row extracted from an OTC documentation page.
type MetricDocumentationEntry struct {
	MetricId    string
	MetricName  string
	Unit        string
	Description string
}

// DocumentationSource describes where to fetch documentation for one OTC namespace.
type DocumentationSource struct {
	Namespace    string
	SubComponent string
	// GithubRawUrl is the raw.githubusercontent.com URL for the RST source.
	GithubRawUrl string
	// HuaweiFallbackUrl is a raw.githubusercontent.com URL to a Huawei Cloud metric catalog
	// Markdown file. When set, cmd/gen-metrics performs a second fetch after the RST pass and
	// adds any metrics present in the Huawei catalog but absent from the OTC RST docs.
	// Empty for sources where the Huawei catalog is incompatible or not useful (see inline comments).
	HuaweiFallbackUrl string
}

const (
	otcGithubRawBase  = "https://raw.githubusercontent.com/opentelekomcloud-docs/"
	huaweiCatalogBase = "https://raw.githubusercontent.com/huaweicloud/cloudeye-exporter/br_release_sdk_v3/cloudservice_metrics/"
)

// DocumentationSources is the list of all known OTC metric documentation sources.
var DocumentationSources = []DocumentationSource{
	{
		Namespace:         "ecs",
		GithubRawUrl:      otcGithubRawBase + "elastic-cloud-server/main/umn/source/cloud_eye_monitoring/basic_ecs_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.ECS.md",
	},
	{
		// bms: Huawei catalog uses namespace SERVICE.BMS; OTC uses SYS.BMS — prefix mismatch, no fallback.
		Namespace:    "bms",
		GithubRawUrl: otcGithubRawBase + "bare-metal-server/main/umn/source/server_monitoring/monitored_metrics_with_agent_installed.rst",
	},
	{
		Namespace:         "as",
		GithubRawUrl:      otcGithubRawBase + "auto-scaling/main/umn/source/as_management/as_group_and_instance_monitoring/monitoring_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.AS.md",
	},
	{
		Namespace:         "evs",
		GithubRawUrl:      otcGithubRawBase + "elastic-volume-service/main/umn/source/cloud_eye_monitoring/viewing_basic_evs_monitoring_data.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.EVS.md",
	},
	{
		Namespace:         "sfs",
		GithubRawUrl:      otcGithubRawBase + "scalable-file-service/main/umn/source/management/monitoring/sfs_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.SFS.md",
	},
	{
		Namespace:         "efs",
		GithubRawUrl:      otcGithubRawBase + "scalable-file-service/main/umn/source/management/monitoring/sfs_turbo_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.EFS.md",
	},
	{
		Namespace:         "cbr",
		GithubRawUrl:      otcGithubRawBase + "cloud-backup-recovery/main/umn/source/cloud_eye_monitoring/viewing_basic_cbr_monitoring_data.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.CBR.md",
	},
	{
		Namespace:         "vpc",
		GithubRawUrl:      otcGithubRawBase + "virtual-private-cloud/main/umn/source/monitoring/supported_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.VPC.md",
	},
	{
		Namespace:         "elb",
		GithubRawUrl:      otcGithubRawBase + "elastic-load-balancing/main/umn/source/monitoring/monitoring_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.ELB.md",
	},
	{
		Namespace:         "nat",
		GithubRawUrl:      otcGithubRawBase + "nat-gateway/main/umn/source/monitoring/supported_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.NAT.md",
	},
	{
		Namespace:         "waf",
		GithubRawUrl:      otcGithubRawBase + "web-application-firewall/main/umn/source/monitoring_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.WAF.md",
	},
	{
		Namespace:         "dms",
		GithubRawUrl:      otcGithubRawBase + "distributed-message-service/main/umn/source/monitoring_and_alarms/kafka_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.DMS.md",
	},
	{
		Namespace:         "dcs",
		GithubRawUrl:      otcGithubRawBase + "distributed-cache-service/main/umn/source/monitoring/dcs_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.DCS.md",
	},
	{
		// rds: Huawei combines all RDS engines in one file; cannot split per SubComponent, no fallback.
		Namespace:    "rds",
		SubComponent: "MySql",
		GithubRawUrl: otcGithubRawBase + "relational-database-service/main/umn/source/working_with_rds_for_mysql/metrics_and_alarms/configuring_displayed_metrics.rst",
	},
	{
		Namespace:    "rds",
		SubComponent: "Postgres",
		GithubRawUrl: otcGithubRawBase + "relational-database-service/main/umn/source/working_with_rds_for_postgresql/metrics_and_alarms/configuring_displayed_metrics.rst",
	},
	{
		Namespace:    "rds",
		SubComponent: "SqlServer",
		GithubRawUrl: otcGithubRawBase + "relational-database-service/main/umn/source/working_with_rds_for_sql_server/metrics_and_alarms/configuring_displayed_metrics.rst",
	},
	{
		Namespace:         "dds",
		GithubRawUrl:      otcGithubRawBase + "document-database-service/main/umn/source/monitoring_and_alarm_reporting/dds_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.DDS.md",
	},
	{
		// nosql: Huawei catalog covers multiple sub-engines (Cassandra, MongoDB, etc.),
		// producing 500+ panels with mostly "No data". Re-add in v2 when dashboards are split per engine.
		Namespace:    "nosql",
		GithubRawUrl: otcGithubRawBase + "gaussdb-nosql/main/umn/source/working_with_gaussdbfor_cassandra/monitoring_and_alarm_reporting/gaussdbfor_cassandra_monitoring_metrics.rst",
	},
	{
		// gaussdb: same multi-engine concern as nosql, no Huawei fallback.
		Namespace:    "gaussdb",
		GithubRawUrl: otcGithubRawBase + "gaussdb-mysql/main/umn/source/working_with_gaussdbfor_mysql/monitoring/configuring_displayed_metrics.rst",
	},
	{
		Namespace:         "gaussdbv5",
		GithubRawUrl:      otcGithubRawBase + "gaussdb-opengauss/main/umn/source/working_with_gaussdb/monitoring_and_alarming/monitoring_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.GAUSSDBV5.md",
	},
	{
		Namespace:         "dws",
		GithubRawUrl:      otcGithubRawBase + "data-warehouse-service/main/umn/source/gaussdbdws_cluster_o_and_m/viewing_gaussdbdws_cluster_monitoring_information_on_cloud_eye.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.DWS.md",
	},
	{
		// css: RST uses namespace SYS.ES (OTC publishes these under SYS.ES).
		Namespace:         "css",
		GithubRawUrl:      otcGithubRawBase + "cloud-search-service/main/umn/source/using_elasticsearch_for_data_search/elasticsearch_cluster_monitoring_and_log_management/elasticsearch_cluster_monitoring_metrics/monitoring_metrics_for_elasticsearch_clusters_in_cloud_eye.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.ES.md",
	},
	{
		Namespace:         "obs",
		GithubRawUrl:      otcGithubRawBase + "object-storage-service/main/umn/source/obs_console_operation_guide/monitoring/monitored_obs_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.OBS.md",
	},
	{
		Namespace:         "dcaas",
		GithubRawUrl:      otcGithubRawBase + "direct-connect/main/umn/source/monitoring/direct_connect_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.DCAAS.md",
	},
	{
		// vpn: Huawei catalog combines both VPN flavours in one file; cannot split per SubComponent, no fallback.
		Namespace:    "vpn",
		SubComponent: "Classic",
		GithubRawUrl: otcGithubRawBase + "virtual-private-network/main/umn/source/management/monitoring/metrics_s2c_classic_vpn.rst",
	},
	{
		Namespace:    "vpn",
		SubComponent: "Enterprise",
		GithubRawUrl: otcGithubRawBase + "virtual-private-network/main/umn/source/management/monitoring/metrics_s2c_enterprise_edition_vpn.rst",
	},
	{
		Namespace:         "apic",
		GithubRawUrl:      otcGithubRawBase + "api-gateway/main/umn/source/monitoring_and_analysis/api_monitoring/monitoring_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.APIC.md",
	},
	{
		// ddm RST has namespace SYS.DDMS; Huawei catalog also uses SYS.DDMS.
		Namespace:         "ddm",
		GithubRawUrl:      otcGithubRawBase + "distributed-database-middleware/main/umn/source/monitoring_management/supported_metrics/ddm_instance_metrics.rst",
		HuaweiFallbackUrl: huaweiCatalogBase + "SYS.DDMS.md",
	},
}
