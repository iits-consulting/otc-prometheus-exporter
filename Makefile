BINARY := otc-prometheus-exporter
DOCKER ?= docker
HELM_UNITTEST_IMAGE ?= helmunittest/helm-unittest:3.18.6-1.0.3
HELM_DOCS_IMAGE ?= jnorwood/helm-docs:v1.14.2

.PHONY: build
build:
	go build -o $(BINARY) .

.PHONY: clean
clean:
	rm -f $(BINARY)

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test -v ./...

.PHONY: check
check: fmt vet lint test helm-lint helm-test

.PHONY: docker-build
docker-build:
	$(DOCKER) build -t $(BINARY) .

.PHONY: helm-lint
helm-lint:
	helm lint ./charts/otc-prometheus-exporter

.PHONY: helm-test
helm-test:
	$(DOCKER) run --rm -v "$(PWD):/apps" $(HELM_UNITTEST_IMAGE) ./charts/otc-prometheus-exporter

.PHONY: helm-docs
helm-docs:
	$(DOCKER) run --rm -v "$(PWD):/helm-docs" $(HELM_DOCS_IMAGE) -c charts/otc-prometheus-exporter -t charts/otc-prometheus-exporter/README.md.gotmpl

.PHONY: generate-dashboards-into-helm
generate-dashboards-into-helm:
	go run ./cmd/grafanadashboards/ --output-path ./charts/otc-prometheus-exporter/dashboards
