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
|-----------------------------|---------------|----------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `OS_USERNAME`               | none          | valid username                   | **REQUIRED** User in the OTC with access to the API                                                                                                                                                                                                                        |
| `OS_PASSWORD`               | none          | valid password                   | **REQUIRED** Password for the user                                                                                                                                                                                                                                         |
| `OS_ACCESS_KEY`             | none          | valid access key                 | **REQUIRED** You can instead of username/password also provide the users AK and SK                                                                                                                                                                                         |
| `OS_SECRET_KEY`             | none          | valid secret key                 | **REQUIRED** You can instead of username/password also provide the users AK and SK                                                                                                                                                                                         |
| `OS_PROJECT_ID`             | none          | valid project id                 | **REQUIRED** Project from which the metrics should be gathered. Obtainable in the OTC console IAM -> Projects -> View your Project -> You can now see the ProjectID                                                                                                        |
| `OS_DOMAIN_NAME`            | none          | valid domain name                | **REQUIRED** Domainname/Tenant ID. The value in the OTC console on the top right `OTC-EU-DE-{somenumberhere}`.                                                                                                                                                             |
| `NAMESPACES`                | none          | e.g. SYS.ECS,SERVICE.BMS,ECS,BMS | **REQUIRED** Specific namespaces for instances you want to get the metrics from. See list of all namespaces in the CloudEye docs [CloudEyeDoc]. If namespace is available in CloudEye docs you can also use namespace without prefix (SYS.ECS -> ECS, SERVICE.BMS -> BMS). |
| `REGION`                    | `eu-de`       | `eu-de`, `eu-nl`                 | Region where your project is located                                                                                                                                                                                                                                       |
| `PORT`                      | `8000`        | any valid unused port            | Port on which metrics are served                                                                                                                                                                                                                                           |
| `WAITDURATION`              | `60`          | any positive integer             | Time in seconds between two API call fetches                                                                                                                                                                                                                               |
| `FETCH_RESOURCE_ID_TO_NAME` | false         | boolean                          | Turns the mapping of resource id to resource name on or off                                                                                                                                                                                                                |

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
