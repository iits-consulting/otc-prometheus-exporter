# OTC Prometheus Exporter

This software gathers metrics from the [Open Telekom Cloud](https://open-telekom-cloud.com/)
for [Prometheus](https://prometheus.io/).

## Available Metrics

Metrics for the following services are available

- Elastic Cloud Server (ECS)
- Virtual Private Cloud (VPC)
- NAT Gateway (NAT)
- Elastic Load Balancing (ELB)
- Distributed Message Service (DMS)

## Requirements

- OTC-Credentials
- Download the tool [otc-auth](https://github.com/iits-consulting/otc-auth)

## Usage & Configuration

In this section you will learn how to use and configure this software.
The configuration happens via environment variables and one configuration file.

1. To authenticate please use `otc-auth`, to do so please follow the instructions in the readme of the project.

2. After authenticating you obtained an unscoped token and scoped token. Set environment variable `PROJECT_NAME` with
   your target project name the way it's written in the config file (`~/.otc-auth-config`).

3. Set the desired namespaces as a list of comma seperated values in the environment variable `NAMESPACES`.

4. The other environment variables are not required. The following table covers all environment variables.

| environment variable   | default value        | allowed values        | description                                                                     |
|------------------------|----------------------|-----------------------|---------------------------------------------------------------------------------|
| `PROJECT_NAME`         | none                 |                       | **REQUIRED** Project name on CloudEye where your instances are running          |
| `NAMESPACES`           | none                 | ECS,DMS,VPC,NAT,ELB   | **REQUIRED** Specific namespaces for instances you want to get the metrics from |
| `PORT`                 | `8000`               | any valid unused port | Port on which metrics are served                                                |
| `WAITDURATION`         | `60`                 | any positive integer  | Time in seconds between two API call fetches                                    |
| `OTC_AUTH_CONFIG_PATH` | `~/.otc-auth-config` | any valid path        | Path to the `.otc-auth-config`                                                  |


### Binary

1. Download and decompress the binary from the release page
2. `chmod +x otc-prometheus-exporter` to make it executable.
3. On MacOs it might be necessary to remove the Apple quarantine property from it too. This can be done with: `xattr -d com.apple.quarantine otc-prometheues-exporter`
4. Export the required environment variables and run the programm.

```shell
export PROJECT_NAME="eu-de_iits-cool-project"
export NAMESPACES="ECS,VPC,RDS"
./otc-prometheues-exporter
```

### Docker

```shell
docker pull ghcr.io/iits-consulting/otc-prometheus-exporter:latest
docker run --platform 'linux/amd64' -e "PROJECT_NAME=eu-de_iits-central" -e "OTC_AUTH_CONFIG_PATH=/data/otc-auth-config" -e "NAMESPACES=VPC,ECS" --mount type=bind,source="/Users/zeljkobekcic/.otc-auth-config",target="/data/otc-auth-config"  ghcr.io/iits-consulting/otc-prometheus-exporter:latest
```

### Kubernetes (Helm)

```shell
helm repo add otc-prometheus-exporter https://iits-consulting.github.io/otc-prometheus-exporter/
helm search repo otc-prometheus-exporter
helm install otc-prometheus-exporter otc-prometheus-exporterE/otc-prometheus-exporter --set your_values.yaml
```

## References

- [Open Telekom Cloud Docs](https://docs.otc.t-systems.com/)
- https://github.com/tiagoReichert/otc-cloudeye-prometheus-exporter
