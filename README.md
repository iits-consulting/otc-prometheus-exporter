# OTC Prometheus Exporter

This software gathers metrics from the [Open Telekom Cloud (OTC)](https://open-telekom-cloud.com/)
for [Prometheus](https://prometheus.io/). The metrics are then usable in any service which can use Prometheus as a
datasource. For example [Grafana](https://grafana.com/)

## Available Metrics

Metrics for the following services are available

- Elastic Cloud Server (ECS)
- Virtual Private Cloud (VPC)
- NAT Gateway (NAT)
- Elastic Load Balancing (ELB)
- Distributed Message Service (DMS)

## Requirements

- OTC-Credentials

## Usage & Configuration

In this section you will learn how to use and configure this software.
The configuration happens via environment variables and one configuration file.

1. Obtain the necessary credentials for the OTC. You need a username, password, project id and domain name.

2. Set the desired namespaces as a list of comma seperated values in the environment variable `NAMESPACES`.

3. The other environment variables are not required. The following table covers all environment variables.

| environment variable        | default value | allowed values                   | description                                                                                                                                                                                                                                                                |
| --------------------------- | ------------- | -------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `OS_USERNAME`               | none          | valid username                   | **REQUIRED** User in the OTC with access to the API                                                                                                                                                                                                                        |
| `OS_PASSWORD`               | none          | valid password                   | **REQUIRED** Password for the user                                                                                                                                                                                                                                         |
| `OS_ACCESS_KEY`             | none          | valid access key                 | **REQUIRED** You can instead of username/password also provide the users AK and SK                                                                                                                                                                                         |
| `OS_SECRET_KEY`             | none          | valid secret key                 | **REQUIRED** You can instead of username/password also provide the users AK and SK                                                                                                                                                                                         |
| `OS_PROJECT_ID`             | none          | valid project id                 | **REQUIRED** Project from which the metrics should be gathered. Obtainable in the OTC console IAM -> Projects -> View your Project -> You can now see the ProjectID                                                                                                        |
| `OS_DOMAIN_NAME`            | none          | valid domain name                | **REQUIRED** Domainname/Tenant ID. The value in the OTC console on the top right `OTC-EU-DE-{somenumberhere}`.                                                                                                                                                             |
| `NAMESPACES`                | none          | e.g. SYS.ECS,SERVICE.BMS,ECS,BMS | **REQUIRED** Specific namespaces for instances you want to get the metrics from. See list of all namespaces in the CloudEye docs [CloudEyeDoc]. If namespace is available in CloudEye docs you can also use namespace without prefix (SYS.ECS -> ECS, SERVICE.BMS -> BMS). |
| `REGION`                    | `eu-de`       | `eu-de`, `eu-nl`                 | Region where your project is located                                                                                                                                                                                                                                       |
| `PORT`                      | `39100`       | any valid unused port            | Port on which metrics are served                                                                                                                                                                                                                                           |
| `WAITDURATION`              | `60`          | any positive integer             | Time in seconds between two API call fetches                                                                                                                                                                                                                               |
| `FETCH_RESOURCE_ID_TO_NAME` | false         | boolean                          | Turns the mapping of resource id to resource name on or off                                                                                                                                                                                                                |


Below is a comprehensive list of metrics that this software can gather. 
Note: The "Translatable" column indicates whether the resource ID can be mapped to its corresponding resource name:


| Category       | Service                                               | Namespace     | Translatable | Reference                                              |
| -------------- | ----------------------------------------------------- | ------------- | ------------ | ------------------------------------------------------ |
| Compute        | Elastic Cloud Server                                  | SYS.ECS       | Yes          | [Basic ECS metrics](https://docs.otc.t-systems.com/elastic-cloud-server/umn/monitoring/basic_ecs_metrics.html) |
|                | Bare Metal Server                                     | SERVICE.BMS   | No           | [BMS Metrics Under OS Monitoring](https://docs.otc.t-systems.com/bare-metal-server/umn/server_monitoring/monitored_metrics_with_agent_installed.html) |
|                | Auto Scaling                                          | SYS.AS        | No           | [AS metrics](https://docs.otc.t-systems.com/auto-scaling/umn/as_management/as_group_and_instance_monitoring/monitoring_metrics.html) |
| Storage        | Elastic Volume Service (attached to an ECS or BMS)    | SYS.EVS       | Yes          | [EVS metrics](https://docs.otc.t-systems.com/elastic-volume-service/umn/viewing_evs_monitoring_data.html) |
|                | Scalable File Service                                 | SYS.SFS       | No           | [SFS metrics](https://docs.otc.t-systems.com/scalable-file-service/umn/management/monitoring/sfs_metrics.html) |
|                | SFS Turbo                                             | SYS.EFS       | Yes          | [SFS Turbo metrics](https://docs.otc.t-systems.com/scalable-file-service/umn/management/monitoring/sfs_turbo_metrics.html) |
|                | Cloud Backup and Recovery                             | SYS.CBR       | Yes          | [CBR metrics](https://docs.otc.t-systems.com/cloud-backup-recovery/umn/monitoring/cbr_metrics.html) |
| Network        | Elastic IP and bandwidth                              | SYS.VPC       | Yes          | [VPC metrics](https://docs.otc.t-systems.com/virtual-private-cloud/umn/operation_guide_new_console_edition/monitoring/supported_metrics.html) |
|                | Elastic Load Balance                                  | SYS.ELB       | Yes          | [ELB metrics](https://docs.otc.t-systems.com/elastic-load-balancing/umn/monitoring/monitoring_metrics.html) |
|                | NAT Gateway                                           | SYS.NAT       | Yes          | [NAT Gateway metrics](https://docs.otc.t-systems.com/nat-gateway/umn/monitoring_management/supported_metrics.html) |
| Security       | Web Application Firewall                              | SYS.WAF       | Yes          | [WAF metrics](https://docs.otc.t-systems.com/web-application-firewall/umn/monitoring_metrics.html) |
| Application    | Distributed Message Service                           | SYS.DMS       | Yes          | [DMS metrics](https://docs.otc.t-systems.com/distributed-message-service/umn/monitoring/kafka_metrics.html) |
|                | Distributed Cache Service                             | SYS.DCS       | Yes          | [DCS metrics](https://docs.otc.t-systems.com/distributed-cache-service/umn/monitoring/dcs_metrics.html) |
| Database       | Relational Database Service                           | SYS.RDS       | Yes          | [RDS for MySQL metrics](https://docs.otc.t-systems.com/relational-database-service/umn/working_with_rds_for_mysql/metrics_and_alarms/configuring_displayed_metrics.html)<br>[RDS for PostgreSQL metrics](https://docs.otc.t-systems.com/relational-database-service/umn/working_with_rds_for_postgresql/metrics_and_alarms/configuring_displayed_metrics.html)<br>[RDS for SQL Server metrics](https://docs.otc.t-systems.com/relational-database-service/umn/working_with_rds_for_sql_server/metrics_and_alarms/configuring_displayed_metrics.html) |
|                | Document Database Service                             | SYS.DDS       | No           | [DDS metrics](https://docs.otc.t-systems.com/document-database-service/umn/monitoring/interconnected_with_cloud_eye/dds_metrics.html) |
|                | GaussDB NoSQL                                         | SYS.NoSQL     | Yes          | [GaussDB(for Cassandra) metrics](https://docs.otc.t-systems.com/gaussdb-nosql/umn/working_with_gaussdbfor_cassandra/monitoring_and_alarm_reporting/gaussdbfor_cassandra_monitoring_metrics.html) |
|                | GaussDB(for MySQL)                                    | SYS.GAUSSDB   | No           | [GaussDB(for MySQL) metrics](https://docs.otc.t-systems.com/gaussdb-mysql/umn/working_with_gaussdbfor_mysql/monitoring/configuring_displayed_metrics.html) |
|                | GaussDB(for openGauss)                                | SYS.GAUSSDBV5 | Yes          | [GaussDB(for openGauss) metrics](https://docs.otc.t-systems.com/gaussdb-opengauss/umn/working_with_gaussdbopengauss/monitoring_and_alarming/metrics.html) |
| Data analysis  | Data Warehouse Service                                | SYS.DWS       | Yes          | [DWS metrics](https://docs.otc.t-systems.com/data-warehouse-service/umn/monitoring_and_alarms/monitoring_clusters_using_cloud_eye.html) |
|                | Cloud Search Service                                  | SYS.ES        | No           | [CSS metrics](https://docs.otc.t-systems.com/cloud-search-service/umn/monitoring_a_cluster/supported_metrics.html) |




[CloudEyeDoc]:https://docs.otc.t-systems.com/cloud-eye/api-ref/appendix/services_interconnected_with_cloud_eye.html#ces-03-0059

### Binary

If you want to run the application directly as a binary then you can do it by following these steps.

1. Download and decompress the binary from the release page
2. `chmod +x otc-prometheus-exporter` to make it executable.
3. On macOS, it might be necessary to remove the Apple quarantine property from it too.
   This can be done with: `xattr -d com.apple.quarantine otc-prometheues-exporter`
4. `cp .env.template .env`
5. Fill out the values in the `.env` file
6. Run the programm: `env $(cat .env) ./otc-prometheues-exporter`

### Docker

If you want a single docker container with the application running then you can do it by following these steps.

1. Make sure you have docker installed and running
2. Copy the `.env.template` to `.env` and fill it out. This makes the docker command much shorter this way and your
   secrets are not listed in your shell history.
3. Run the following:

```shell
docker pull ghcr.io/iits-consulting/otc-prometheus-exporter:latest
docker run --env-file .env ghcr.io/iits-consulting/otc-prometheus-exporter:latest
```

### Docker Compose

If you want to start the application, a Prometheus and Grafana server all at once then you can do it by following these steps.
This is suitable for a quick test or local development because the entire tool chain is running.

1. Make sure you have docker and docker-compose installed and running
2. Copy the `.env.template` to `.env` and fill it out. This makes the docker command much shorter this way and your
   secrets are not listed in your shell history.
3. Run the following: `docker compose --env-file .env up`

### Kubernetes (Helm)

```shell
helm repo add otc-prometheus-exporter https://iits-consulting.github.io/otc-prometheus-exporter/
helm search repo otc-prometheus-exporter
helm install otc-prometheus-exporter otc-prometheus-exporter/otc-prometheus-exporter --set your_values.yaml
```

## References

- [Open Telekom Cloud Docs](https://docs.otc.t-systems.com/)
- https://github.com/tiagoReichert/otc-cloudeye-prometheus-exporter
