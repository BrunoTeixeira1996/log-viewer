#/usr/bin/sh

go build -o binaries/exporter ./cmd/journalctl-exporter/exporter.go
go build -o binaries/viewer ./cmd/journalctl-viewer/viewer.go
