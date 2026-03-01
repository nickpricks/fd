#!/bin/bash
set -e

echo "=> Fetching dependencies..."
go mod tidy

echo "=> Building FeatherTrailMD..."
make build

echo "=> Installing to GOPATH..."
make install

echo "=> Installation complete! Run 'ft help' to get started."
