#!/bin/sh
# The gates of tacplus-dev/launcher (gates.yaml): build, vet and tidy the module,
# check it with govulncheck, build the static linux binaries the Java images
# take (CGO_ENABLED=0, amd64 and arm64) and check each with govulncheck too, and
# require that a build leaves the checkout as it found it (tacplus-dev/launcher#4).
set -eu
cd "$(dirname "$0")"
govulncheck="go run golang.org/x/vuln/cmd/govulncheck@v1.9.0"
before=$(git status --porcelain)
go build ./...
go vet ./...
go mod tidy -diff
$govulncheck ./...
out=$(mktemp -d)
trap 'rm -rf "$out"' EXIT
for arch in amd64 arm64; do
  CGO_ENABLED=0 GOOS=linux GOARCH=$arch go build -trimpath -ldflags="-s -w" -o "$out/launcher-linux-$arch" .
  $govulncheck -mode=binary "$out/launcher-linux-$arch"
done
after=$(git status --porcelain)
if [ "$before" != "$after" ]; then
  echo "gates: the build changed the checkout:" >&2
  git status --porcelain >&2
  exit 1
fi
echo "gates: ok"
