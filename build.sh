#!/usr/bin/env bash
set -euo pipefail

OS=${1:-darwin}
ARCH=${2:-amd64}

BUILD_TIME=$(date +"%Y%m%d.%H%M%S")
CommitHash=$(git rev-parse HEAD)
GoVersion=$(go version | cut -c 14- | cut -d' ' -f1)
GitTag=$(git describe --tags)

TRG_PKG="github.com/awwithro/makemea/cmd"
FLAG="-X $TRG_PKG.BuildTime=$BUILD_TIME"
FLAG="$FLAG -X $TRG_PKG.CommitHash=$CommitHash"
FLAG="$FLAG -X $TRG_PKG.GoVersion=$GoVersion"
FLAG="$FLAG -X $TRG_PKG.GitTag=$GitTag"

GOOS=${OS} GOARCH=${ARCH} go build -v -ldflags "$FLAG" -o makemea
