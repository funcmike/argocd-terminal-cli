#!/bin/sh
set -eo pipefail

mkdir -p build

# osx
export GOOS=darwin GOARCH=amd64
go build -o build/atc-macos-$GOARCH cmd/atc/main.go
export GOOS=darwin GOARCH=arm64
go build -o build/atc-macos-$GOARCH cmd/atc/main.go

# linux
export GOOS=linux GOARCH=amd64
go build -o build/atc-$GOOS-$GOARCH cmd/atc/main.go
export GOOS=linux GOARCH=arm64
go build -o build/atc-$GOOS-$GOARCH cmd/atc/main.go

# windows
export GOOS=windows GOARCH=amd64
go build -o build/atc-$GOOS-$GOARCH.exe cmd/atc/main.go
export GOOS=windows GOARCH=arm64
go build -o build/atc-$GOOS-$GOARCH.exe cmd/atc/main.go

unset GOOS GOARCH



