#ª/bin/bash
mkdir -p dist
rm -Rf dist/*

go mod tidy
go build \
-o dist/rtsp-recorder \
-ldflags="-X github.com/Miguel-Dorta/rtsp-recorder/pkg.Version=$(git describe --tags)" \
./cmd/rtsp-recorder
