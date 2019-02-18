#!/bin/sh

export GO111MODULE=on
go get github.com/mitchellh/gox
gox -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmd/b2b