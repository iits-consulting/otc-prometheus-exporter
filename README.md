# OTC Prometheus Exporter

This software gathers metrics from the [Open Telekom Cloud (OTC)](https://open-telekom-cloud.com/)
for [Prometheus](https://prometheus.io/). The metrics are then usable in any service which can use Prometheus as a
datasource. For example [Grafana](https://grafana.com/)

## How It Works

Each OTC service (namespace) is scraped independently via a `?namespace=` query parameter:

```
GET /metrics?namespace=SYS.ECS   -> ECS metrics
GET /metrics?namespace=SYS.RDS   -> RDS metrics
GET /metrics?namespace=SYS.ELB   -> ELB metrics
...
```

This allows Prometheus to scrape each namespace with its own interval and timeout, and if one service fails only that scrape shows as down.

Append `&enrich=false` to skip service-specific API calls and return only CES metrics (useful for debugging or reducing API load).

A `/healthz` endpoint returns 200 for liveness/readiness probes.

# OTC CES Metric Namespaces

All known Cloud Eye Service (CES) namespaces for Open Telekom Cloud.

Authoritative source:
[Services Interconnected with Cloud Eye](https://docs.otc.t-systems.com/cloud-eye/api-ref/appendix/services_interconnected_with_cloud_eye.html)

## Confirmed on OTC

These namespaces are listed on the official OTC "Services Interconnected with
Cloud Eye" page. Any of them can be scraped by the exporter — CES-only namespaces
require no code changes, just add them to `targets.namespaces` in `values.yaml`.

| Namespace           | Service                                  | Metrics Documentation                                                                                                                                                                                                                                                               |
|---------------------|------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `AGT.ECS`           | Elastic Cloud Server (OS Agent)          | [ECS Monitoring](https://docs.otc.t-systems.com/elastic-cloud-server/umn/monitoring/index.html)                                                                                                                                                                                     |
| `SERVICE.BMS`       | Bare Metal Server                        | [BMS Metrics](https://docs.otc.t-systems.com/bare-metal-server/umn/server_monitoring/monitored_metrics_with_agent_installed.html)                                                                                                                                                   |
| `SYS.ALARM`         | CES Alarm Rules (service-API-only)       | [CES Alarms](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/alarm_rules/index.html)                                                                                                                                                                               |
| `SYS.AS`            | Auto Scaling                             | [AS Metrics](https://docs.otc.t-systems.com/auto-scaling/umn/as_management/as_group_and_instance_monitoring/monitoring_metrics.html)                                                                                                                                                |
| `SYS.CBR`           | Cloud Backup and Recovery                | [CBR Metrics](https://docs.otc.t-systems.com/cloud-backup-recovery/umn/cloud_eye_monitoring/viewing_basic_cbr_monitoring_data.html)                                                                                                                                                 |
| `SYS.DCS`           | Distributed Cache Service (Redis)        | [DCS Metrics](https://docs.otc.t-systems.com/distributed-cache-service/umn/monitoring/dcs_metrics.html)                                                                                                                                                                             |
| `SYS.DDS`           | Document Database Service (MongoDB)      | [DDS Metrics](https://docs.otc.t-systems.com/document-database-service/umn/monitoring_and_alarm_reporting/dds_metrics.html)                                                                                                                                                         |
| `SYS.DMS`           | Distributed Message Service (Kafka)      | [DMS Metrics](https://docs.otc.t-systems.com/distributed-message-service/umn/monitoring/kafka_metrics.html)                                                                                                                                                                         |
| `SYS.DWS`           | Data Warehouse Service                   | [DWS Metrics](https://docs.otc.t-systems.com/data-warehouse-service/umn/monitoring_and_alarms/monitoring_clusters_using_cloud_eye.html)                                                                                                                                             |
| `SYS.ECS`           | Elastic Cloud Server                     | [ECS Metrics](https://docs.otc.t-systems.com/elastic-cloud-server/umn/monitoring/basic_ecs_metrics.html)                                                                                                                                                                            |
| `SYS.EFS`           | SFS Turbo                                | [EFS Metrics](https://docs.otc.t-systems.com/scalable-file-service/umn/management/monitoring/sfs_turbo_metrics.html)                                                                                                                                                                |
| `SYS.ELB`           | Elastic Load Balance                     | [ELB Metrics](https://docs.otc.t-systems.com/elastic-load-balancing/umn/monitoring/monitoring_metrics.html)                                                                                                                                                                         |
| `SYS.ES`            | Cloud Search Service (CSS/Elasticsearch) | [CSS Metrics](https://docs.otc.t-systems.com/cloud-search-service/umn/using_elasticsearch_for_data_search/elasticsearch_cluster_monitoring_and_log_management/elasticsearch_cluster_monitoring_metrics/monitoring_metrics_for_elasticsearch_clusters_in_cloud_eye.html#css-01-0042) |
| `SYS.EVS`           | Elastic Volume Service                   | [EVS Metrics](https://docs.otc.t-systems.com/elastic-volume-service/umn/cloud_eye_monitoring/viewing_basic_evs_monitoring_data.html)                                                                                                                                                |
| `SYS.GAUSSDB`       | GaussDB for MySQL (TaurusDB)             | [GaussDB Metrics](https://docs.otc.t-systems.com/gaussdb-mysql/umn/working_with_gaussdbfor_mysql/monitoring/configuring_displayed_metrics.html)                                                                                                                                     |
| `SYS.GAUSSDBV5`     | GaussDB for openGauss                    | [GaussDBV5 Metrics](https://docs.otc.t-systems.com/gaussdb-opengauss/umn/working_with_gaussdbopengauss/monitoring_and_alarming/metrics.html)                                                                                                                                        |
| `SYS.NAT`           | NAT Gateway                              | [NAT Metrics](https://docs.otc.t-systems.com/nat-gateway/umn/monitoring_management/supported_metrics.html)                                                                                                                                                                          |
| `SYS.NoSQL`         | GaussDB NoSQL (GeminiDB/Cassandra)       | [NoSQL Metrics](https://docs.otc.t-systems.com/gaussdb-nosql/umn/working_with_gaussdbfor_cassandra/monitoring_and_alarm_reporting/gaussdbfor_cassandra_monitoring_metrics.html)                                                                                                     |
| `SYS.OBS`           | Object Storage Service                   | [OBS Metrics](https://docs.otc.t-systems.com/object-storage-service/umn/obs_console_operation_guide/monitoring/obs_monitoring_metrics.html)                                                                                                                                         |
| `SYS.RDS`           | Relational Database Service              | [RDS Metrics](https://docs.otc.t-systems.com/relational-database-service/umn/working_with_rds_for_mysql/metrics_and_alarms/configuring_displayed_metrics.html)                                                                                                                      |
| `SYS.SFS`           | Scalable File Service                    | [SFS Metrics](https://docs.otc.t-systems.com/scalable-file-service/umn/management/monitoring/sfs_metrics.html)                                                                                                                                                                      |
| `SYS.VPC`           | Virtual Private Cloud (EIP/Bandwidth)    | [VPC Metrics](https://docs.otc.t-systems.com/virtual-private-cloud/umn/monitoring/supported_metrics.html)                                                                                                                                                                           |
| `SYS.VPN`           | Virtual Private Network                  | [VPN Docs](https://docs.otc.t-systems.com/virtual-private-network/umn/management/monitoring/metrics_s2c_classic_vpn.html)                                                                                                                                                           |
| `SYS.WAF`           | Web Application Firewall                 | [WAF Metrics](https://docs.otc.t-systems.com/web-application-firewall/umn/monitoring_metrics.html)                                                                                                                                                                                  |
| `SYS.APIG`          | API Gateway (Shared)                     | [APIG Monitoring](https://docs.otc.t-systems.com/api-gateway/umn/monitoring_and_analysis/api_monitoring/monitoring_metrics.html#apig-03-0032)                                                                                                                                       |
| `SYS.APIC`          | API Gateway (Dedicated)                  | [APIG Monitoring](https://docs.otc.t-systems.com/api-gateway/umn/monitoring_and_analysis/api_monitoring/monitoring_metrics.html#apig-03-0032)                                                                                                                                       |
| `SYS.CDM`           | Cloud Data Migration                     | [CDM Metrics](https://docs.otc.t-systems.com/data-arts-studio/umn/user_guide/dataarts_migration/managing_clusters/viewing_metrics/cdm_metrics.html)                                                                                                                                 |
| `SYS.DCAAS`         | Direct Connect                           | [DCAAS Docs](https://docs.otc.t-systems.com/direct-connect/umn/monitoring/direct_connect_metrics.html)                                                                                                                                                                              |
| `SYS.DDM`           | Distributed Database Middleware          | [DDM Docs](https://docs.otc.t-systems.com/distributed-database-middleware/umn/monitoring_management/supported_metrics/ddm_instance_metrics.html)                                                                                                                                    |
| `SYS.DLI`           | Data Lake Insight                        | [DLI Docs](https://docs.otc.t-systems.com/data-lake-insight/umn/monitoring_dli_using_cloud_eye.html)                                                                                                                                                                                |
| `SYS.DRS`           | Data Replication Service                 | [DRS Docs](https://docs.otc.t-systems.com/data-replication-service/umn/real-time_synchronization/interconnecting_with_cloud_eye/supported_metrics.html)                                                                                                                             |
| `SYS.ER`            | Enterprise Router                        | [ER Docs](https://docs.otc.t-systems.com/enterprise-router/umn/monitoring_and_audit/cloud_eye_monitoring/supported_metrics.html)                                                                                                                                                    |
| `SYS.FunctionGraph` | FunctionGraph                            | [FunctionGraph Metrics](https://docs.otc.t-systems.com/function-graph/umn/monitoring/metrics/functiongraph_metrics.html)                                                                                                                                                            |

### Note: ECS/BMS OS-Level Metrics Require ICAgent

Some ECS and BMS metrics (e.g., `disk_read_bytes_rate`, `disk_write_bytes_rate`, detailed memory metrics) require the **ICAgent (Telescope)** to be installed inside the guest OS. Without it, these metrics will show as "No data" in Cloud Eye and will not appear in the exporter output.

Basic hypervisor-level metrics (CPU utilization, network traffic) are available without the agent. See the [ICAgent installation guide](https://docs.otc.t-systems.com/cloud-eye/umn/server_monitoring/installing_and_configuring_the_agent_on_a_linux_ecs.html) for details.

### Special Case: Object Storage Service (OBS)

The Object Storage Service (OBS) metrics have some unique considerations:

1. **Global Project ID Requirement**: Unlike other services OBS is a global service. This means you need a global project ID to gather OBS metrics.
   The supported projects are `eu-de` and `eu-nl`.
2. **Limitation on Project Scoped Metrics**: OBS metrics cannot be collected with project scoped metrics since a global project ID is needed, which transcends individual project scopes.

## Requirements

- OTC credentials (username/password or AK/SK)
- IAM permissions (see below)

## IAM Permissions

The exporter needs read-only access to the Cloud Eye (CES) API and to each service's listing API.

### Simplest: Tenant Guest role

Assign the **Tenant Guest** role at both **global** and **regional project** scope. This grants `get*`, `list*`, `head*` on all services except IAM.

### Fine-grained: ReadOnlyAccess policies (least privilege)

**CES ReadOnlyAccess** is the core requirement -- it covers reading CES metrics for all namespaces. Additional per-service policies are only needed for services where the exporter calls the service-specific API:

| Policy                              | Needed for                                                                                                                      |
|-------------------------------------|---------------------------------------------------------------------------------------------------------------------------------|
| **CES ReadOnlyAccess**              | CES metrics for all namespaces (required)                                                                                       |
| **ECS ReadOnlyAccess**              | List servers (name mapping + instance status)                                                                                   |
| **AOM ReadOnlyAccess**              | AOM host metrics with `ecs_aom_` prefix (required for per-device/per-NIC/per-mountpoint metrics when ECS enrichment is enabled) |
| **RDS ReadOnlyAccess**              | List RDS instances (name mapping + status/nodes)                                                                                |
| **ELB ReadOnlyAccess**              | List load balancers (name mapping + operating status)                                                                           |
| **DMS ReadOnlyAccess**              | List Kafka instances (name mapping + storage/partitions)                                                                        |
| **NAT ReadOnlyAccess**              | List NAT gateways (name mapping + status)                                                                                       |
| **DCS ReadOnlyAccess**              | List Redis instances (name mapping + status/capacity)                                                                           |
| **DDS ReadOnlyAccess**              | List MongoDB instances (name mapping + node status)                                                                             |
| **VPC ReadOnlyAccess**              | List bandwidths (name mapping + size)                                                                                           |
| **CBR ReadOnlyAccess**              | List backups (status/size)                                                                                                      |
| **AutoScaling ReadOnlyAccess**      | List scaling groups (instance counts/status)                                                                                    |
| **BMS ReadOnlyAccess**              | List bare metal servers (name mapping + status)                                                                                 |
| **EVS ReadOnlyAccess**              | List volumes (name mapping + status/size)                                                                                       |
| **CSS ReadOnlyAccess**              | List Elasticsearch clusters (name mapping + status)                                                                             |
| **DWS ReadOnlyAccess**              | List data warehouse clusters (name mapping + status)                                                                            |
| **SFS ReadOnlyAccess**              | List file shares (name mapping + status)                                                                                        |
| **SFS Turbo ReadOnlyAccess**        | List SFS Turbo shares (name mapping + status)                                                                                   |
| **Direct Connect ReadOnlyAccess**   | List virtual gateways (name mapping + status)                                                                                   |

> **Note:** If the service API call fails (e.g. missing policy), enrichment degrades gracefully — the scrape still succeeds but resource names and `*_status` metrics will be absent for that namespace. A warning is logged.

CES-only namespaces (WAF, OBS, GaussDB, GaussDBv5, NoSQL, VPN) need **no additional policies** beyond CES ReadOnlyAccess.

All policies must be assigned at the **regional project** scope for the relevant region (`eu-de` or `eu-nl`).

## Usage & Configuration

### Environment Variables

| Environment variable   | Default | Description                                                                                  |
|------------------------|---------|----------------------------------------------------------------------------------------------|
| `OS_USERNAME`          | -       | **Required*** OTC username                                                                   |
| `OS_PASSWORD`          | -       | **Required*** OTC password                                                                   |
| `OS_ACCESS_KEY`        | -       | **Required*** AK (alternative to username/password)                                          |
| `OS_SECRET_KEY`        | -       | **Required*** SK (alternative to username/password)                                          |
| `OS_PROJECT_ID`        | -       | **Required** OTC project ID                                                                  |
| `OS_DOMAIN_NAME`       | -       | **Required** OTC domain name / tenant ID                                                     |
| `OS_REGION_PROJECT_ID` | -       | Region-level project ID for global services (OBS). Auto-discovered if not set.               |
| `REGION`               | `eu-de` | OTC region (`eu-de` or `eu-nl`)                                                              |
| `PORT`                 | `39100` | HTTP server port                                                                             |
| `LOG_LEVEL`            | `INFO`  | Log level (`DEBUG`, `INFO`, `WARN`, `ERROR`)                                                 |
| `REQUEST_TIMEOUT`      | `10`    | HTTP request timeout in seconds for OTC API calls                                            |
| `IDLE_CONN_TIMEOUT`    | `90`    | How long idle HTTP connections stay in the pool in seconds (set higher than scrape interval) |
| `COLLECT_TIMEOUT`      | `55`    | Maximum collect time in seconds per scrape (should be less than Prometheus scrape_timeout)   |
| `CES_BATCH_SIZE`       | `500`   | Max metrics per CES batch API request                                                        |
| `CES_LOOKBACK`         | `10`    | CES lookback window in minutes (how far back to query for datapoints)                        |
| `AOM_BATCH_SIZE`       | `20`    | Max metrics per AOM data API request                                                         |
| `AOM_CONCURRENCY`      | `5`     | Max concurrent AOM API calls per scrape                                                      |

*Either username+password or access-key+secret-key is required (mutually exclusive).

### Kubernetes (Helm)

```shell
helm repo add otc-prometheus-exporter https://iits-consulting.github.io/otc-prometheus-exporter/
helm search repo otc-prometheus-exporter
helm install otc-prometheus-exporter otc-prometheus-exporter/otc-prometheus-exporter -f your_values.yaml
```

The Helm chart uses a **PodMonitor** to scrape each namespace independently. All configuration is driven by `targets.namespaces` — the same block controls scraping, dashboards, and alert rules:

```yaml
targets:
  defaults:
    enabled: true
    enrich: true       # enrich CES metrics with names from service APIs
    dashboard: true    # deploy Grafana dashboard ConfigMap
    rules: false       # deploy PrometheusRule alerts (only if alerts/ file exists)
  namespaces:
    SYS.ECS: {}        # inherits all defaults
    SYS.RDS:
      rules: true      # enable RDS alerts
    SYS.WAF:
      enabled: false   # disabled — not scraped, no dashboard, no alerts

dashboards:
  enabled: true        # deploy dashboard ConfigMaps
  selfMonitoring: true # include the exporter self-monitoring dashboard
  serviceHealth: true  # include the Service Health Overview dashboard (aggregates all *_status metrics)

prometheusRules:
  enabled: true        # deploy PrometheusRule CRDs

podMonitor:
  enabled: true
  selfMonitoring: true # scrape exporter's own metrics (go_*, otc_scrape_*, otc_http_*)
```

See [`charts/otc-prometheus-exporter/values.yaml`](charts/otc-prometheus-exporter/values.yaml) for the full reference.

### Docker and Docker Compose

**Single container:**

```shell
cp .env.template .env
# Fill in your OTC credentials
docker run --env-file .env ghcr.io/iits-consulting/otc-prometheus-exporter:latest
```

**Docker Compose** (with Prometheus + Grafana):

```shell
cp .env.template .env
# Fill in your OTC credentials
docker-compose --env-file .env up --build otc-prometheus-exporter -d
```

This starts the exporter, Prometheus (localhost:9090), and Grafana (localhost:3000). Prometheus is pre-configured to scrape all namespaces. Edit `docker.local/prometheus/prometheus.yaml` to adjust which namespaces are scraped.

### Binary

1. Download and decompress the binary from the release page
2. `chmod +x otc-prometheus-exporter`
3. On macOS: `xattr -d com.apple.quarantine otc-prometheus-exporter`
4. `cp .env.template .env` and fill out the values
5. Run: `env $(cat .env) ./otc-prometheus-exporter`

To check the installed version: `./otc-prometheus-exporter --version`

## References

- [Open Telekom Cloud Docs](https://docs.otc.t-systems.com/)
- https://github.com/tiagoReichert/otc-cloudeye-prometheus-exporter
