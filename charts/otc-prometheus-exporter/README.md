# otc-prometheus-exporter

![Version: 1.2.0](https://img.shields.io/badge/Version-1.2.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.1](https://img.shields.io/badge/AppVersion-0.5.1-informational?style=flat-square)

A Helm chart for Kubernetes

## Installing the Chart

To install the chart with the release name otc-prometheus-exporter:

```shell
    helm repo add otc-prometheus-exporter https://iits-consulting.github.io/otc-prometheus-exporter
    helm search repo otc-prometheus-exporter
    helm install otc-prometheus-exporter otc-prometheus-exporter/otc-prometheus-exporter
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| fullnameOverride | string | `""` |  |
| nameOverride | string | `""` |  |
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
| deployment.env.FETCH_RESOURCE_ID_TO_NAME | bool | `false` |  |
| deployment.env.OS_DOMAIN_NAME | string | `"OTC-EU-DE-00000000001000058635"` |  |
| deployment.env.OS_PASSWORD | string | `""` |  |
| deployment.env.OS_PROJECT_ID | string | `"f0f45389d6a947d88c8658fb8e1a1053"` |  |
| deployment.env.OS_USERNAME | string | `""` |  |
| deployment.env.PORT | int | `39100` |  |
| deployment.env.REGION | string | `"eu-de"` |  |
| deployment.env.WAITDURATION | int | `60` |  |
| deployment.envFromSecret | object | `{}` |  |
| deployment.health.liveness.path | string | `"/metrics"` |  |
| deployment.health.liveness.periodSeconds | int | `180` |  |
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
| deployment.securityContext | string | `nil` |  |
| deployment.volumeMounts | object | `{}` |  |
| deployment.volumes | object | `{}` |  |
| service.ports.metrics.port | int | `39100` |  |
| service.ports.metrics.targetPort | int | `39100` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| serviceMonitor.enabled | bool | `true` |  |

<img src="../../img/iits.svg" alt="iits consulting" id="logo" width="200" height="200">

<br>

*This chart is provided by [iits-consulting](https://iits-consulting.de/) - your Cloud-Native Innovation Teams as a Service!*

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.13.0](https://github.com/norwoodj/helm-docs/releases/v1.13.0)
