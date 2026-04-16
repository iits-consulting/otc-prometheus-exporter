# otc-prometheus-exporter

![Version: 2.0.0](https://img.shields.io/badge/Version-2.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 2.0.0](https://img.shields.io/badge/AppVersion-2.0.0-informational?style=flat-square)

Prometheus exporter for Open Telekom Cloud (OTC) metrics

## Installing the Chart

To install the chart with the release name otc-prometheus-exporter:

```shell
    helm repo add otc-prometheus-exporter https://iits-consulting.github.io/otc-prometheus-exporter
    helm search repo otc-prometheus-exporter
    helm install otc-prometheus-exporter otc-prometheus-exporter/otc-prometheus-exporter
```

## Additional Prometheus Rules

This chart can be used to create default `PrometheusRule` definitions for core services (e.g., RDS, OBS, ELB). 
If you need to monitor additional metrics, use the `additionalPrometheusRulesMap` field in your `values.yaml`:

```yaml
additionalPrometheusRulesMap:
  # Prometheus rules for RDS PostgreSQL
  rds-postgresql-alerts:
    groups:
      - name: postgresql-system
        rules:
          - alert: PostgreSQLHighCPUUtilization
            annotations:
              summary: '{{ $labels.instance }} CPU > 80%'
              description: 'CPU utilization for PostgreSQL instance {{ $labels.instance }} has been above 80%. Current value: {{ $value }}%'
            expr: >
              rds_rds001_cpu_util > 0.8
            labels:
              severity: warning
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| additionalPrometheusRulesMap | object | `{}` | Additional custom PrometheusRules (not tied to namespace files). Each key becomes a separate PrometheusRule resource. |
| dashboards | object | `{"enabled":true,"selfMonitoring":true,"serviceHealth":true}` | Grafana dashboard ConfigMap generation. Creates one ConfigMap per enabled namespace with grafana_dashboard: "1" label. |
| dashboards.enabled | bool | `true` | Deploy dashboard ConfigMaps |
| dashboards.selfMonitoring | bool | `true` | Deploy the exporter self-monitoring dashboard (HTTP traces, scrape durations) |
| dashboards.serviceHealth | bool | `true` | Deploy the Service Health Overview dashboard (aggregates all *_status metrics) |
| deployment | object | `{"affinity":{},"env":{"AOM_BATCH_SIZE":"20","AOM_CONCURRENCY":"5","CES_BATCH_SIZE":"500","CES_LOOKBACK":"10","IDLE_CONN_TIMEOUT":"90","REQUEST_TIMEOUT":"10"},"envFromSecret":"","health":{"liveness":{"path":"/healthz","port":39100},"readiness":{"path":"/healthz","port":39100},"startupProbe":{"path":"/healthz","port":39100}},"image":{"pullPolicy":"IfNotPresent","repository":"ghcr.io/iits-consulting/otc-prometheus-exporter","tag":""},"imagePullSecrets":[],"podAnnotations":{},"podSecurityContext":{"runAsGroup":65534,"runAsNonRoot":true,"runAsUser":65534,"seccompProfile":{"type":"RuntimeDefault"}},"ports":{"metrics":{"port":39100,"targetPort":39100}},"replicaCount":1,"resources":{"limits":{"memory":"128Mi"},"requests":{"cpu":"200m","memory":"128Mi"}},"securityContext":{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true},"volumeMounts":{},"volumes":{}}` | Deployment configuration |
| deployment.affinity | object | `{}` | Pod affinity/anti-affinity rules |
| deployment.env | object | `{"AOM_BATCH_SIZE":"20","AOM_CONCURRENCY":"5","CES_BATCH_SIZE":"500","CES_LOOKBACK":"10","IDLE_CONN_TIMEOUT":"90","REQUEST_TIMEOUT":"10"}` | Environment variables for authentication and tuning. Credentials should be set via envFromSecret in production. |
| deployment.env.REQUEST_TIMEOUT | string | `"10"` | Tuning (defaults are usually fine) |
| deployment.envFromSecret | string | `""` | Name of an existing Secret to source additional environment variables from. Use this for credentials instead of putting them in env directly. |
| deployment.health | object | `{"liveness":{"path":"/healthz","port":39100},"readiness":{"path":"/healthz","port":39100},"startupProbe":{"path":"/healthz","port":39100}}` | Health check configuration. All probes hit the /healthz endpoint. |
| deployment.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| deployment.image.repository | string | `"ghcr.io/iits-consulting/otc-prometheus-exporter"` | Container image repository |
| deployment.image.tag | string | `""` | Image tag (defaults to chart appVersion) |
| deployment.imagePullSecrets | list | `[]` | Image pull secrets for private registries |
| deployment.podAnnotations | object | `{}` | Additional pod annotations |
| deployment.podSecurityContext | object | `{"runAsGroup":65534,"runAsNonRoot":true,"runAsUser":65534,"seccompProfile":{"type":"RuntimeDefault"}}` | Pod-level security context |
| deployment.ports | object | `{"metrics":{"port":39100,"targetPort":39100}}` | Container port configuration |
| deployment.replicaCount | int | `1` | Number of replicas. Keep at 1 — multiple replicas produce duplicate metrics which break dashboards and alerting. |
| deployment.resources | object | `{"limits":{"memory":"128Mi"},"requests":{"cpu":"200m","memory":"128Mi"}}` | Container resource requests and limits |
| deployment.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true}` | Container-level security context |
| deployment.volumeMounts | object | `{}` | Additional volume mounts for the container |
| deployment.volumes | object | `{}` | Additional volumes to mount into the pod |
| fullnameOverride | string | `""` | Override the full resource name (release-name + chart-name) |
| nameOverride | string | `""` | Override the chart name used in resource names |
| namespaceOverride | string | `""` | Override the deployment namespace (useful for multi-namespace setups) |
| podMonitor | object | `{"annotations":{},"defaults":{"interval":"30s","scrapeTimeout":"15s"},"enabled":true,"labels":{},"selfMonitoring":true}` | PodMonitor configuration for Prometheus Operator scraping. Creates one scrape endpoint per enabled namespace, each passing ?namespace=SYS.XXX. |
| podMonitor.annotations | object | `{}` | Additional annotations on the PodMonitor |
| podMonitor.defaults | object | `{"interval":"30s","scrapeTimeout":"15s"}` | Default scrape settings (can be overridden per namespace) |
| podMonitor.enabled | bool | `true` | Deploy the PodMonitor resource |
| podMonitor.labels | object | `{}` | Additional labels on the PodMonitor (e.g. for Prometheus selector matching) |
| podMonitor.selfMonitoring | bool | `true` | Scrape the exporter's own metrics (go_*, process_*, otc_scrape_*, otc_http_*) |
| prometheusRules | object | `{"enabled":false}` | PrometheusRule generation. Creates one PrometheusRule per namespace that has rules enabled and a matching alerts/ file. |
| prometheusRules.enabled | bool | `false` | Deploy PrometheusRule resources |
| serviceAccount | object | `{"annotations":{},"create":true,"name":""}` | Service account configuration |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account (e.g. for IRSA/workload identity) |
| serviceAccount.create | bool | `true` | Create a dedicated service account |
| serviceAccount.name | string | `""` | Use an existing service account instead of creating one |
| targets | object | `{"defaults":{"dashboard":true,"enabled":true,"enrich":true,"rules":false},"namespaces":{"SERVICE.BMS":{"enabled":false},"SYS.ALARM":{"rules":true},"SYS.AS":{"enabled":false},"SYS.CBR":{"enabled":false},"SYS.DCS":{"enabled":false},"SYS.DDS":{"enabled":false},"SYS.DMS":{"enabled":false},"SYS.DWS":{"enabled":false},"SYS.ECS":null,"SYS.EFS":{"enabled":false},"SYS.ELB":{"rules":true},"SYS.ES":{"dashboard":"css","enabled":false},"SYS.EVS":{"enabled":false},"SYS.GAUSSDB":{"enabled":false},"SYS.GAUSSDBV5":{"enabled":false},"SYS.NAT":null,"SYS.NoSQL":{"enabled":false},"SYS.OBS":{"rules":true},"SYS.RDS":{"dashboard":true,"rules":true},"SYS.SFS":{"enabled":false},"SYS.VPC":null,"SYS.VPN":{"enabled":false},"SYS.WAF":{"enabled":false}}}` | Per-namespace target configuration. Controls which OTC namespaces are scraped, and whether dashboards/alerts are deployed. |
| targets.defaults | object | `{"dashboard":true,"enabled":true,"enrich":true,"rules":false}` | Defaults applied to every namespace (can be overridden per namespace) |
| targets.defaults.dashboard | bool | `true` | Deploy a Grafana dashboard ConfigMap. Set to false to skip. |
| targets.defaults.enabled | bool | `true` | Scrape this namespace |
| targets.defaults.enrich | bool | `true` | Enrich CES metrics with human-readable names from service APIs and gather various additional data points. |
| targets.defaults.rules | bool | `false` | Deploy PrometheusRule alerts. Set to true to enable, "name" to override filename. Only namespaces with a matching file in alerts/ will produce a resource. (Note: not all namespaces have pre-defines rules) |
| targets.namespaces | object | `{"SERVICE.BMS":{"enabled":false},"SYS.ALARM":{"rules":true},"SYS.AS":{"enabled":false},"SYS.CBR":{"enabled":false},"SYS.DCS":{"enabled":false},"SYS.DDS":{"enabled":false},"SYS.DMS":{"enabled":false},"SYS.DWS":{"enabled":false},"SYS.ECS":null,"SYS.EFS":{"enabled":false},"SYS.ELB":{"rules":true},"SYS.ES":{"dashboard":"css","enabled":false},"SYS.EVS":{"enabled":false},"SYS.GAUSSDB":{"enabled":false},"SYS.GAUSSDBV5":{"enabled":false},"SYS.NAT":null,"SYS.NoSQL":{"enabled":false},"SYS.OBS":{"rules":true},"SYS.RDS":{"dashboard":true,"rules":true},"SYS.SFS":{"enabled":false},"SYS.VPC":null,"SYS.VPN":{"enabled":false},"SYS.WAF":{"enabled":false}}` | Namespace-specific overrides. Keys are CES namespace names (e.g. SYS.ECS). Names for dashboards and alerts are auto-derived: SYS.ECS -> ecs, SERVICE.BMS -> bms. |
| targets.namespaces."SYS.ES".dashboard | string | `"css"` | CSS uses a non-standard dashboard name |

<img src="../../img/iits.svg" alt="iits consulting" id="logo" width="200" height="200">

<br>

*This chart is provided by [iits-consulting](https://iits-consulting.de/) - your Cloud-Native Innovation Teams as a Service!*

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)