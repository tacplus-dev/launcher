# State: launcher

Overwritten to stay current (loopers and agents). Last updated: 2026-11-02 by looper (#4).

- **Version / line:** unreleased (no tags); a plain Go module, `go 1.27.0` with `toolchain go1.27.3`; no tacplus floor
- **Works:** forks `LAUNCH_CMD`, forwards SIGTERM/SIGINT, propagates the exit code, optional gRPC health polling; the land gates (`gates.sh`, `gates.yaml`) build, vet, tidy-check and govulncheck the source and the static linux amd64/arm64 binaries, all clean (2026-11-02); no binary is committed (#4)
- **Released / deployed:** not released from this repository; copied into the kernel, eventbus and stratum Java images (README)
- **Known gaps:** no tests; GitHub CI runs only the planning files check (govulncheck runs in the land gates); no release pipeline yet for `launcher-linux-<arch>` and `checksums.txt` (tacplus-dev/looper#17); fleet-seatrials pins `launcher_commit` c575e14, whose go.mod has the vulnerable grpc v1.81.1 / x/net v0.51.0 (PROPOSALS.md)
- **Blocked on:** nothing
