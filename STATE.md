# State: launcher

Overwritten to stay current (loopers and agents). Last updated: 2026-10-28 by looper (#3).

- **Version / line:** unreleased (no tags); a plain Go module, `go 1.27.0` with `toolchain go1.27.3`; no tacplus floor
- **Works:** forks `LAUNCH_CMD`, forwards SIGTERM/SIGINT, propagates the exit code, optional gRPC health polling; `go build`, `go vet`, `go mod tidy -diff` clean and `govulncheck ./...` reports no vulnerabilities (2026-10-28)
- **Released / deployed:** not released from this repository; copied into the kernel, eventbus and stratum Java images (README)
- **Known gaps:** no tests; CI runs only the planning files check (no build or govulncheck gate); a stray macOS binary `launcher` is committed (PROPOSALS.md)
- **Blocked on:** nothing
