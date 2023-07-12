FROM golang:1.19-alpine AS build

ARG CGO=1
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/iits-consulting/otc-prometheus-exporter
COPY . /go/src/github.com/iits-consulting/otc-prometheus-exporter

RUN go build -o otc-prometheus-exporter cmd/main.go && \
    mv otc-prometheus-exporter /usr/local/bin

FROM alpine:3.16
COPY --from=build /usr/local/bin/otc-prometheus-exporter /usr/bin/otc-prometheus-exporter
EXPOSE 39100
ENTRYPOINT [ "otc-prometheus-exporter" ]
