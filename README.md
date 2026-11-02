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

## Gates

`./gates.sh` (`gates.yaml`): `go build`, `go vet`, `go mod tidy -diff`,
`govulncheck` of the source and of both static linux binaries, and a checkout
the build leaves clean.
