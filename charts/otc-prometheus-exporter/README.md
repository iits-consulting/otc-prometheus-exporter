# otc-prometheus-exporter

![Version: 1.3.0](https://img.shields.io/badge/Version-1.3.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.6.6](https://img.shields.io/badge/AppVersion-0.6.6-informational?style=flat-square)

A Helm chart for Kubernetes

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
| dashboards.as.enabled | bool | `true` |  |
| dashboards.bms.enabled | bool | `true` |  |
| dashboards.cbr.enabled | bool | `true` |  |
| dashboards.css.enabled | bool | `true` |  |
| dashboards.dcs.enabled | bool | `true` |  |
| dashboards.dds.enabled | bool | `true` |  |
| dashboards.dms.enabled | bool | `true` |  |
| dashboards.dws.enabled | bool | `true` |  |
| dashboards.ecs.enabled | bool | `true` |  |
| dashboards.efs.enabled | bool | `true` |  |
| dashboards.elb.enabled | bool | `true` |  |
| dashboards.evs.enabled | bool | `true` |  |
| dashboards.gaussdb.enabled | bool | `true` |  |
| dashboards.gaussdbv5.enabled | bool | `true` |  |
| dashboards.nat.enabled | bool | `true` |  |
| dashboards.nosql.enabled | bool | `true` |  |
| dashboards.rds-mysql.enabled | bool | `true` |  |
| dashboards.rds-postgres.enabled | bool | `true` |  |
| dashboards.rds-sqlserver.enabled | bool | `true` |  |
| dashboards.sfs.enabled | bool | `true` |  |
| dashboards.vpc.enabled | bool | `true` |  |
| dashboards.waf.enabled | bool | `true` |  |
| deployment.affinity | object | `{}` |  |
| deployment.env | object | `{}` |  |
| deployment.envFromSecret | string | `""` |  |
| deployment.health.liveness.path | string | `"/metrics"` |  |
| deployment.health.liveness.port | int | `39100` |  |
| deployment.health.readiness.path | string | `"/metrics"` |  |
| deployment.health.readiness.port | int | `39100` |  |
| deployment.health.startupProbe.path | string | `"/metrics"` |  |
| deployment.health.startupProbe.port | int | `39100` |  |
| deployment.image.pullPolicy | string | `"IfNotPresent"` |  |
| deployment.image.repository | string | `"ghcr.io/iits-consulting/otc-prometheus-exporter"` |  |
| deployment.image.tag | string | `""` |  |
| deployment.imagePullSecrets | list | `[]` |  |
| deployment.podAnnotations | object | `{}` |  |
| deployment.podSecurityContext | string | `nil` |  |
| deployment.ports.metrics.port | int | `39100` |  |
| deployment.ports.metrics.targetPort | int | `39100` |  |
| deployment.replicaCount | int | `1` |  |
| deployment.resources.limits.memory | string | `"128Mi"` |  |
| deployment.resources.requests.cpu | string | `"100m"` |  |
| deployment.resources.requests.memory | string | `"128Mi"` |  |
| deployment.securityContext | string | `nil` |  |
| deployment.volumeMounts | object | `{}` |  |
| deployment.volumes | object | `{}` |  |
| fullnameOverride | string | `""` |  |
| nameOverride | string | `""` |  |
| service.ports.metrics.port | int | `39100` |  |
| service.ports.metrics.targetPort | int | `39100` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| serviceMonitor.enabled | bool | `true` |  |

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

<img src="../../img/iits.svg" alt="iits consulting" id="logo" width="200" height="200">

<br>

*This chart is provided by [iits-consulting](https://iits-consulting.de/) - your Cloud-Native Innovation Teams as a Service!*

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)