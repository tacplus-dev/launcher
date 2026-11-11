# launcher

Process supervisor for distroless Java containers. It forks the command in
`LAUNCH_CMD`, forwards SIGTERM/SIGINT to the child, propagates the child's exit
code, and optionally polls the child's gRPC health endpoint so a deadlocked JVM
(alive, not serving) is killed and the container restarted.

Built static from source and copied into every Java service image as that
service's entrypoint binary (see the `launcher-linux-<arch>` COPY in kernel,
eventbus and stratum Dockerfiles). The fleet's build image carries it at
`/usr/local/bin/launcher`. No binary is committed here: a release ships
`launcher-linux-amd64`, `launcher-linux-arm64` and `checksums.txt`, built as

```sh
for arch in amd64 arm64; do
  CGO_ENABLED=0 GOOS=linux GOARCH=$arch go build -trimpath -ldflags="-s -w" -o launcher-linux-$arch .
done
```

`release.sh DIR` builds them so, checks each with `govulncheck -mode=binary` and
writes `checksums.txt` (sha256), the assets of a `vX.Y.Z` GitHub release:
`https://github.com/tacplus-dev/launcher/releases/download/vX.Y.Z/<asset>`.
Pin a version and its sha256 from `checksums.txt`. (The tag workflow that runs it
and attaches them is tacplus-dev/launcher#5, waiting on tacplus-dev/looper#106.)

## Gates

`./gates.sh` (`gates.yaml`): `go build`, `go vet`, `go mod tidy -diff`,
`govulncheck` of the source, the release assets built and checked
(`tests/release.test.sh`: static, stripped, `-trimpath`, `checksums.txt` verifies,
`govulncheck` of both binaries), and a checkout the build leaves clean.
They run, with no Docker, in the image `build.lock` pins by digest (tacplus build 2.6.8, the fleet's build
image), from which tacplus-dev/build plans every build of a commit (looper#71); to move it, change its digest:
a new digest is a change like any other, landed on the full gates.
