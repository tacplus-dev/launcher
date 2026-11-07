#!/bin/sh
# The release assets of tacplus-dev/launcher (#5), written to DIR: the static
# linux binaries launcher-linux-amd64 and launcher-linux-arm64, built as the
# README says, each checked with govulncheck -mode=binary, and checksums.txt,
# their sha256. The release workflow attaches them to a vX.Y.Z tag's release;
# the gates build them too (tests/release.test.sh).
set -eu
[ $# -eq 1 ] || { echo "usage: release.sh DIR" >&2; exit 2; }
out=$(mkdir -p "$1" && cd "$1" && pwd)
cd "$(dirname "$0")"
govulncheck="go run golang.org/x/vuln/cmd/govulncheck@v1.9.0"
for arch in amd64 arm64; do
  CGO_ENABLED=0 GOOS=linux GOARCH=$arch go build -trimpath -ldflags="-s -w" -o "$out/launcher-linux-$arch" .
  $govulncheck -mode=binary "$out/launcher-linux-$arch"
done
(cd "$out" && sha256sum launcher-linux-amd64 launcher-linux-arm64 > checksums.txt)
