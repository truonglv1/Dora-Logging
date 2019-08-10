#!/usr/bin/env bash
echo "remove application..."
rm bin/dora-logging
echo "build application..."
go build -o bin/dora-logging cmd/main.go