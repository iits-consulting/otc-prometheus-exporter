FROM golang:1.20-alpine AS build

ARG CGO=0
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/iits-consulting/otc-prometheus-exporter
COPY . /go/src/github.com/iits-consulting/otc-prometheus-exporter

RUN go build -o otc-prometheus-exporter main.go && \
    mv otc-prometheus-exporter /usr/local/bin

FROM alpine:3.17
COPY --from=build /usr/local/bin/otc-prometheus-exporter /usr/bin/otc-prometheus-exporter
EXPOSE 39100
ENTRYPOINT [ "otc-prometheus-exporter" ]
