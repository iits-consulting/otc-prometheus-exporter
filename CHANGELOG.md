# Changelog

## [2.0.0](https://github.com/iits-consulting/otc-prometheus-exporter/compare/v1.3.0...v2.0.0) (2026-04-15)

> ⚠️ **Upgrading from v0.6.6? v1.x was never published — this is a direct jump.**
> Read the [migration guide](https://github.com/iits-consulting/otc-prometheus-exporter/blob/main/MIGRATION.md) before upgrading. Key breaks: metric names changed, Helm chart restructured, new IAM requirements for 7 namespaces, dashboards replaced.


### ⚠ BREAKING CHANGES

* extend providers and dashboards ([#76](https://github.com/iits-consulting/otc-prometheus-exporter/issues/76))

### Features

* Add resource name resolutions for ELB, DCS, VPC ([#39](https://github.com/iits-consulting/otc-prometheus-exporter/issues/39)) ([4beee5b](https://github.com/iits-consulting/otc-prometheus-exporter/commit/4beee5bc18a3759c986bb91005af278b020b36de))
* AK/SK Authentication ([#6](https://github.com/iits-consulting/otc-prometheus-exporter/issues/6)) ([e0452ff](https://github.com/iits-consulting/otc-prometheus-exporter/commit/e0452ffa02630192466e07ac44a947350553446b))
* **chart:** add support for default and additional Prometheus rules ([#69](https://github.com/iits-consulting/otc-prometheus-exporter/issues/69)) ([f6d4049](https://github.com/iits-consulting/otc-prometheus-exporter/commit/f6d4049fe63316995be58eaf602f541efa257062))
* expose version via --version flag ([#86](https://github.com/iits-consulting/otc-prometheus-exporter/issues/86)) ([be51adb](https://github.com/iits-consulting/otc-prometheus-exporter/commit/be51adba09f48905ef342bf94f39eceda2e4aca5))
* extend providers and dashboards ([#76](https://github.com/iits-consulting/otc-prometheus-exporter/issues/76)) ([0472c91](https://github.com/iits-consulting/otc-prometheus-exporter/commit/0472c9148ee9baa95b024b8abf28a272994629fc))
* **grafana:** generate dashboards for all namespaces based on the official documentation ([#47](https://github.com/iits-consulting/otc-prometheus-exporter/issues/47)) ([bf0b06d](https://github.com/iits-consulting/otc-prometheus-exporter/commit/bf0b06d3cfcc367af0c82176f72b8d8ac57913f7))
* **helm:** adding prometheus service monitor ([7f7f601](https://github.com/iits-consulting/otc-prometheus-exporter/commit/7f7f601d7772bc22c3daee715477c91695ecd0a7))
* **logging:** Add zap logging and fix some small issues ([#38](https://github.com/iits-consulting/otc-prometheus-exporter/issues/38)) ([e6c7380](https://github.com/iits-consulting/otc-prometheus-exporter/commit/e6c73802f9ed52b6b829c8eb81d68b62a1bbc86f))
* rewrite exporter with per-namespace on-demand scraping ([#75](https://github.com/iits-consulting/otc-prometheus-exporter/issues/75)) ([88223ee](https://github.com/iits-consulting/otc-prometheus-exporter/commit/88223ee40fdab0d2e0aa22c92005368354e18e0a))


### Bug Fixes

* **batchdl:** Don't panic if batch download request fails ([eb9195c](https://github.com/iits-consulting/otc-prometheus-exporter/commit/eb9195c0b4fe32134589d12cce0fad22f3a608fe))
* **chart:** update helm-docs chartInstallation section ([c939221](https://github.com/iits-consulting/otc-prometheus-exporter/commit/c939221d9c2bfa27d1885df7226dd4c1726719a1))
* **ci:** fix cli arguments ([5ef231b](https://github.com/iits-consulting/otc-prometheus-exporter/commit/5ef231b93f588a38f17215d79d63d0229dbbc583))
* **ci:** fix pipeline for helm documentation autogeneration  ([#46](https://github.com/iits-consulting/otc-prometheus-exporter/issues/46)) ([b234475](https://github.com/iits-consulting/otc-prometheus-exporter/commit/b234475ba13ccd08332ab35a702e562ef711b61b))
* **ci:** updated main path in goreleaser conv ([3337796](https://github.com/iits-consulting/otc-prometheus-exporter/commit/333779663bc5fe446dd9c197492412754cf91a2a))
* Comment out default env variable to allow override when used as helm dependency ([5c876b1](https://github.com/iits-consulting/otc-prometheus-exporter/commit/5c876b1fac9e37e52e08cb0aaed148fd8d587cbd))
* don't add domain name when authenticating with AK/SK ([bc35f1d](https://github.com/iits-consulting/otc-prometheus-exporter/commit/bc35f1d20869ef38e715734b82f793bd529ddf09))
* don't add domain name when authenticating with AK/SK ([986222f](https://github.com/iits-consulting/otc-prometheus-exporter/commit/986222fac1ad5ea5cdea40e134d0db89502e38ef))
* increase start time to fetch dcs metrics too ([#36](https://github.com/iits-consulting/otc-prometheus-exporter/issues/36)) ([e01ddb6](https://github.com/iits-consulting/otc-prometheus-exporter/commit/e01ddb66aa990f59da89e17d2bff8273170432ea))
* **lint:** fix linting issue ([273e9d8](https://github.com/iits-consulting/otc-prometheus-exporter/commit/273e9d8be27d728084dd20197427aa58d2762afb))
* Move name overrides in values file to the correct level ([e7fef61](https://github.com/iits-consulting/otc-prometheus-exporter/commit/e7fef6150123522a76ff27d8cbc15deb81ea2726))
* Update chart README.md to reflect movement of name override variable ([5f46279](https://github.com/iits-consulting/otc-prometheus-exporter/commit/5f46279328c88e1fa22203c70b3aa698da4ce439))
* Use correct property for  envFromSecret ([4a47776](https://github.com/iits-consulting/otc-prometheus-exporter/commit/4a47776e31cfd48f25690944616f6d5d88d0af73))

## v1.0.0 — Complete Rewrite

This is a full rewrite of the exporter. The architecture changed fundamentally — please read the breaking changes carefully before upgrading.

### Breaking Changes

#### Scraping model changed: poll&cache → on-demand-pull

The old exporter ran a background loop that fetched all namespaces on a fixed interval (`WAITDURATION`) and served cached metrics. The new exporter fetches metrics **on demand** — each Prometheus scrape triggers a real-time API call for exactly the requested namespace via `?namespace=SYS.XXX`.

#### Removed environment variables

| Removed                     | Replacement                                                                       |
|-----------------------------|-----------------------------------------------------------------------------------|
| `NAMESPACES`                | No longer needed — Prometheus scrape config controls which namespaces are fetched |
| `WAITDURATION`              | No longer needed — Prometheus scrape interval controls fetch frequency            |
| `FETCH_RESOURCE_ID_TO_NAME` | Replaced by `?enrich=false` query parameter (default: enrichment enabled)         |

#### New / renamed environment variables

| Variable               | Default         | Description                                                                                      |
|------------------------|-----------------|--------------------------------------------------------------------------------------------------|
| `REGION`               | `eu-de`         | **New.** OTC region (`eu-de` or `eu-nl`)                                                         |
| `OS_REGION_PROJECT_ID` | auto-discovered | **New.** Region-level project ID for global services (OBS). Auto-discovered with user/pass auth. |
| `REQUEST_TIMEOUT`      | `10`            | **New.** HTTP timeout in seconds for OTC API calls                                               |
| `IDLE_CONN_TIMEOUT`    | `90`            | **New.** Idle HTTP connection pool timeout in seconds                                            |
| `COLLECT_TIMEOUT`      | `55`            | **New.** Max collect time per scrape in seconds                                                  |
| `CES_BATCH_SIZE`       | `500`           | **New.** Max metrics per CES batch API request                                                   |
| `CES_LOOKBACK`         | `10`            | **New.** CES lookback window in minutes                                                          |
| `AOM_BATCH_SIZE`       | `20`            | **New.** Max metrics per AOM data API request                                                    |
| `AOM_CONCURRENCY`      | `5`             | **New.** Max concurrent AOM API calls per scrape                                                 |
| `LOG_LEVEL`            | `INFO`          | **New.** Log level (DEBUG, INFO, WARN, ERROR)                                                    |

#### Helm chart: ServiceMonitor → PodMonitor

The chart no longer creates a `Service` + `ServiceMonitor`. It now uses a **PodMonitor** with one `podMetricsEndpoint` per enabled namespace. If you have custom `ServiceMonitor` configurations, migrate them to the new `podMonitor` section in values.yaml.

#### Helm chart: values.yaml restructured

The values.yaml structure changed significantly:

| Old                           | New                                                                        |
|-------------------------------|----------------------------------------------------------------------------|
| `service`                     | Removed — no Service resource created                                      |
| `serviceMonitor`              | `podMonitor` (with `selfMonitoring` toggle)                                |
| `serviceMonitor.namespaces`   | `targets.namespaces` (unified config for scraping, dashboards, and alerts) |
| `prometheusRules.rules.rds`   | `targets.namespaces.SYS.RDS.rules: true` (per-namespace)                   |
| `prometheusRules.rules.elb`   | `targets.namespaces.SYS.ELB.rules: true` (per-namespace)                   |
| `prometheusRules.rules.obs`   | `targets.namespaces.SYS.OBS.rules: true` (per-namespace)                   |
| `prometheusRules.rules.alarm` | `targets.namespaces.SYS.ALARM.rules: true` (per-namespace)                 |

Similarly, the path for dashboards toggles changed:

| Old                                     | New                                              |
|-----------------------------------------|--------------------------------------------------|
| `dashboards.ecs.enabled`                | `targets.namespaces.SYS.ECS.dashboard: true`     |
| `dashboards.rds-mysql.enabled`          | `targets.namespaces.SYS.RDS.dashboard: true`     |
| `dashboards.elb.enabled`                | `targets.namespaces.SYS.ELB.dashboard: true`     |
| (all other `dashboards.<name>.enabled`) | `targets.namespaces.<NAMESPACE>.dashboard: true` |

#### Metric names changed

Metric names now follow a consistent `{namespace}_{metric}` pattern derived from the CES namespace and metric name. Old custom metric names are no longer produced. Check your dashboards and alert rules for references to old metric names.

### New Features

#### Per-namespace scraping

Each OTC namespace is scraped independently via `?namespace=SYS.XXX`. This means:
- Individual scrape intervals and timeouts per namespace
- One failing namespace doesn't block others
- Prometheus shows per-namespace scrape health (up/down)

#### Service-API enrichment

Providers that call service-specific APIs (ECS, RDS, ELB, DMS, NAT, DCS, DDS, VPC, CBR, AS) now enrich CES metrics with human-readable resource names and additional status/capacity metrics. Disable per-scrape with `?enrich=false`.

#### AOM host metrics (ECS)

When enrichment is enabled, ECS scrapes also fetch per-device AOM metrics (CPU, memory, disk I/O per device, network per NIC, filesystem per mountpoint) under the `ecs_aom_` prefix. Requires `AOM ReadOnlyAccess` and ICAgent installed in guest OS.

#### New services

Added providers for: VPN (`SYS.VPN`), CES Alarms (`SYS.ALARM`), Auto Scaling status metrics (`SYS.AS`), CBR backup status (`SYS.CBR`).

#### Grafana dashboards

Pre-generated dashboards for all 24 namespaces, deployed as ConfigMaps with `grafana_dashboard: "1"` label. Includes a self-monitoring dashboard for the exporter itself (HTTP trace metrics, scrape durations).

#### PrometheusRule alerts

Bundled alert rules for RDS, ELB, OBS, and Alarms, deployed as PrometheusRule CRDs. Enable per namespace via `targets.namespaces.SYS.XXX.rules: true`.
