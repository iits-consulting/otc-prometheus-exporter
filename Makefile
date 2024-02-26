.PHONY: test
test:
	go test -v ./...

.PHONY: generate-dashboards-into-helm
generate-dashboards-into-helm: 
	go run cmd/grafanadashboards/main.go ./... --output-path ./charts/otc-prometheus-exporter/dashboards
			
