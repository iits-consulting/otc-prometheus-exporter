# OTC Prometheus Exporter

This software gathers metrics from the [Open Telekom Cloud](https://open-telekom-cloud.com/) for [Prometheus](https://prometheus.io/).

## Available Metrics

Metrics for the following services are available

- Elastic Cloud Server
- Virtual Private Cloud
- NAT Gateway
- Elastic Load Balancing
- Distributed Message Service

## Required Credentials 
Since you're calling an API you'll to authenticate you'll need following information:
- domain name
- user name
- password
- project name/id (for most of the services)

## Usage & Configuration

In this section you will learn how to use and configure this software.
The configuration happens via environment variables and one configuration file.

1. Download the tool [otc-auth](https://github.com/iits-consulting/otc-auth)
2. Run `otc-auth` to obtain an unscoped and scoped token. The scoped token is required here, such that project based metrics can be fetched

4. Configure by setting environment variables:
    
    | environment variable | default value | allowed values        | description                                                        |
    | -------------------- | ------------- | --------------------- | ------------------------------------------------------------------ |
    | `PROJECT_NAME`       | none          |                       | Project name on CloudEye where your instances are running          |
    | `NAMESPACES`         | none          | ECS,DMS,VPC,NAT,ELB   | specific namespaces for instances you want to get the metrics from |
    | `PORT`               | `8000`        | any valid unused port | port on which metrics are served                                   |
    | `WAITDURATION`       | `60`          | any positive integer  | time in seconds between two API call fetches                       |

5. Run with `go run cmd/main.go`



## References

- [Open Telekom Cloud Docs](https://docs.otc.t-systems.com/)
-  <!--- [other sources](https://github.com/tiagoReichert/otc-cloudeye-prometheus-exporter) --->