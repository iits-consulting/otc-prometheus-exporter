# Migration Guide

## v0.6.6 → v2.0.0

> **Note on versioning:** v1.x was never published. The version was bumped directly to v2.0.0 to align the Go binary, Docker image, and Helm chart under a single version number going forward.

This release is a complete rewrite. The architecture, scraping model, metric names, Helm chart structure, and IAM requirements have all changed. Read through each section before upgrading.

---

### 1. Scraping model changed: poll & cache → on-demand pull

The old exporter ran a background loop that fetched all configured namespaces on a fixed interval and cached the result. The new exporter fetches metrics **on demand** — each Prometheus scrape triggers a live API call for exactly the requested namespace via `?namespace=SYS.XXX`.

This means Prometheus (not the exporter) controls fetch frequency, and each namespace is scraped independently with its own interval and timeout.

---

### 2. Removed environment variables

| Removed variable        | What to do instead                                                        |
|-------------------------|---------------------------------------------------------------------------|
| `NAMESPACES`            | Remove it. Namespaces are now controlled by your Prometheus scrape config (or `targets.namespaces` in values.yaml for Helm). |
| `WAITDURATION`          | Remove it. Use Prometheus `scrape_interval` per job instead.              |
| `FETCH_RESOURCE_ID_TO_NAME` | Remove it. Enrichment is now on by default. Append `?enrich=false` to a scrape URL to disable it. |

---

### 3. New environment variables

| Variable               | Default         | Description                                                        |
|------------------------|-----------------|--------------------------------------------------------------------|
| `REGION`               | `eu-de`         | OTC region (`eu-de` or `eu-nl`)                                    |
| `OS_REGION_PROJECT_ID` | auto-discovered | Region-level project ID for OBS. Auto-discovered with user/pass auth. |
| `REQUEST_TIMEOUT`      | `10`            | HTTP timeout in seconds for OTC API calls                          |
| `IDLE_CONN_TIMEOUT`    | `90`            | Idle HTTP connection pool timeout in seconds                       |
| `COLLECT_TIMEOUT`      | `55`            | Max collect time per scrape in seconds                             |
| `CES_BATCH_SIZE`       | `500`           | Max metrics per CES batch API request                              |
| `CES_LOOKBACK`         | `10`            | CES lookback window in minutes                                     |
| `AOM_BATCH_SIZE`       | `20`            | Max metrics per AOM data API request                               |
| `AOM_CONCURRENCY`      | `5`             | Max concurrent AOM API calls per scrape                            |
| `LOG_LEVEL`            | `INFO`          | Log level (`DEBUG`, `INFO`, `WARN`, `ERROR`)                       |

---

### 4. Metric names changed

Metric names now follow a consistent `{namespace}_{metric}` pattern derived from the CES namespace and metric name (e.g. `sys_ecs_cpu_util`). The old custom metric names are no longer produced.

**You must update all dashboards and alert rules that reference old metric names.**

---

### 5. Helm chart restructured

#### ServiceMonitor replaced by PodMonitor

The chart no longer creates a `Service` or `ServiceMonitor`. It creates a **PodMonitor** with one scrape endpoint per enabled namespace. Migrate your scrape configuration to the new `podMonitor` section.

#### values.yaml structure changed

| Old path                            | New path                                                         |
|-------------------------------------|------------------------------------------------------------------|
| `service`                           | Removed — no Service resource is created                         |
| `serviceMonitor`                    | `podMonitor`                                                     |
| `serviceMonitor.namespaces`         | `targets.namespaces`                                             |
| `prometheusRules.rules.rds`         | `targets.namespaces.SYS.RDS.rules: true`                         |
| `prometheusRules.rules.elb`         | `targets.namespaces.SYS.ELB.rules: true`                         |
| `prometheusRules.rules.obs`         | `targets.namespaces.SYS.OBS.rules: true`                         |
| `prometheusRules.rules.alarm`       | `targets.namespaces.SYS.ALARM.rules: true`                       |
| `dashboards.ecs.enabled`            | `targets.namespaces.SYS.ECS.dashboard: true`                     |
| `dashboards.rds-mysql.enabled`      | `targets.namespaces.SYS.RDS.dashboard: true`                     |
| `dashboards.elb.enabled`            | `targets.namespaces.SYS.ELB.dashboard: true`                     |
| (all other `dashboards.<x>.enabled`) | `targets.namespaces.<NAMESPACE>.dashboard: true`                |
| `deployment.env.NAMESPACES`         | Remove — see section 2                                           |
| `deployment.env.WAITDURATION`       | Remove — see section 2                                           |
| `deployment.env.FETCH_RESOURCE_ID_TO_NAME` | Remove — see section 2                                  |

---

### 6. Grafana dashboards replaced

All bundled Grafana dashboards have been regenerated from the official OTC documentation. Dashboard UIDs and panel IDs have changed.

**What breaks:**
- Dashboards you customized by editing the bundled ConfigMap JSONs.
- Links or bookmarks using a Grafana dashboard UID.
- Alert rules or annotations that reference a specific panel ID.

**What to do:**
Export any customizations from Grafana before upgrading. After upgrading, re-apply them to the new dashboards.

A new **Service Health Overview** dashboard (`dashboards.serviceHealth: true`) is deployed by default. It aggregates all `*_status` metrics across namespaces into a single overview. Disable with `dashboards.serviceHealth: false`.

---

### 7. New IAM requirements

The exporter now calls service-specific APIs for more namespaces. The following were CES-only in v0.6.6 and now require additional policies:

| Namespace     | New API call                 | Policy to add                 |
|---------------|------------------------------|-------------------------------|
| `SERVICE.BMS` | List bare metal servers      | BMS ReadOnlyAccess            |
| `SYS.EVS`     | List block storage volumes   | EVS ReadOnlyAccess            |
| `SYS.ES`      | List Elasticsearch clusters  | CSS ReadOnlyAccess            |
| `SYS.DWS`     | List data warehouse clusters | DWS ReadOnlyAccess            |
| `SYS.SFS`     | List file shares             | SFS ReadOnlyAccess            |
| `SYS.EFS`     | List SFS Turbo shares        | SFS Turbo ReadOnlyAccess      |
| `SYS.DCAAS`   | List Direct Connect gateways | Direct Connect ReadOnlyAccess |

These already had service API calls in v0.6.6 and still require their respective policies:
ECS, RDS, ELB, DMS, NAT, DCS, DDS, VPC, CBR, AutoScaling.

**If a required policy is missing:** the scrape does **not** fail. Enrichment degrades gracefully — resource names will be absent from labels and `*_status` metrics will not be emitted. A `WARN` log entry is written. To suppress the API call entirely, set `enrich: false` for that namespace in values.yaml or scrape with `?enrich=false`.

The simplest approach remains assigning the **Tenant Guest** role at both global and regional project scope.
