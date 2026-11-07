#!/bin/sh
# The test of release.sh (tacplus-dev/launcher#5): it writes launcher-linux-amd64,
# launcher-linux-arm64 and checksums.txt, nothing else; the binaries are static
# linux builds made as the README says (CGO_ENABLED=0, -trimpath, -s -w), and
# checksums.txt is their sha256 in sha256sum's format, by asset name.
set -eu
cd "$(dirname "$0")/.."
out=$(mktemp -d)
trap 'rm -rf "$out"' EXIT
fail() { echo "release.test: $*" >&2; exit 1; }
./release.sh "$out"
[ "$(ls "$out" | tr '\n' ' ')" = "checksums.txt launcher-linux-amd64 launcher-linux-arm64 " ] ||
  fail "assets: $(ls "$out" | tr '\n' ' ')"
(cd "$out" && sha256sum -c --strict checksums.txt) || fail "checksums.txt does not verify"
[ "$(cut -d' ' -f3 "$out/checksums.txt" | tr '\n' ' ')" = "launcher-linux-amd64 launcher-linux-arm64 " ] ||
  fail "checksums.txt names: $(cat "$out/checksums.txt")"
for arch in amd64 arm64; do
  info=$(go version -m "$out/launcher-linux-$arch")
  for want in "CGO_ENABLED=0" "GOOS=linux" "GOARCH=$arch" "-trimpath=true"; do
    echo "$info" | grep -qF -- "build	$want" || fail "launcher-linux-$arch: no build $want in: $info"
  done
  # -trimpath leaves -ldflags out of the build info: -s shows as no symbol table.
  if go tool nm "$out/launcher-linux-$arch" >/dev/null 2>&1; then fail "launcher-linux-$arch: not stripped (-s -w)"; fi
done
echo "release.test: ok"
