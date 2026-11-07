#!/bin/sh
# The gates of tacplus-dev/launcher (gates.yaml): build, vet and tidy the module,
# check it with govulncheck, build the static linux binaries the Java images
# take (CGO_ENABLED=0, amd64 and arm64) and their checksums.txt as a release does
# (release.sh, tests/release.test.sh), each binary checked with govulncheck too, and
# require that a build leaves the checkout as it found it (tacplus-dev/launcher#4).
set -eu
cd "$(dirname "$0")"
govulncheck="go run golang.org/x/vuln/cmd/govulncheck@v1.9.0"
before=$(git status --porcelain)
go build ./...
go vet ./...
go mod tidy -diff
$govulncheck ./...
tests/release.test.sh
after=$(git status --porcelain)
if [ "$before" != "$after" ]; then
  echo "gates: the build changed the checkout:" >&2
  git status --porcelain >&2
  exit 1
fi
echo "gates: ok"
