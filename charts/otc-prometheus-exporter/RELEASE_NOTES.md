Helm chart release for {{ .Name }} `{{ .Version }}` (app version `{{ .AppVersion }}`).

{{ .Description }}

## Installation

```sh
helm repo add otc-prometheus-exporter https://iits-consulting.github.io/otc-prometheus-exporter
helm install otc-prometheus-exporter otc-prometheus-exporter/otc-prometheus-exporter
```

## Docker Image

```
ghcr.io/iits-consulting/otc-prometheus-exporter:{{ .Version }}
```

## Changelog

See the [full release notes](https://github.com/iits-consulting/otc-prometheus-exporter/releases/tag/v{{ .AppVersion }}) for this version.
