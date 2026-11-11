# State: launcher

Overwritten to stay current (loopers and agents). Last updated: 2026-11-11 by looper (#5).

- **Version / line:** unreleased (no tags; v0.1.0 is #5); a plain Go module, `go 1.27.0` with `toolchain go1.27.3`; no tacplus floor
- **Works:** forks `LAUNCH_CMD`, forwards SIGTERM/SIGINT, propagates the exit code, optional gRPC health polling; the land gates (`gates.sh`, `gates.yaml`) build, vet, tidy-check and govulncheck the source and the static linux amd64/arm64 binaries, all clean (2026-11-02); no binary is committed (#4); `release.sh` builds the release assets (both binaries, govulncheck-checked, and `checksums.txt`), tested in the gates by `tests/release.test.sh` (#5); a `vX.Y.Z` tag on main runs `.github/workflows/release.yml`, which builds them with `release.sh` and attaches them to the tag's GitHub release (#5); pinned for the build pipeline (#6): `build.lock` names tacplus build 2.6.1 by digest, in which the gates pass with no Docker
- **Released / deployed:** not released yet: v0.1.0 (#5) for kernel, eventbus and stratum to pin (tacplus-dev/kernel#22); copied into the kernel, eventbus and stratum Java images (README)
- **Known gaps:** no tests of `main.go` (only `tests/release.test.sh`, of the release build); GitHub CI runs only the planning files check (govulncheck runs in the land gates); the release binaries are built by GitHub Actions on the tag, not by the deterministic build pipeline (tacplus-dev/looper#17); fleet-seatrials pins `launcher_commit` c575e14, whose go.mod has the vulnerable grpc v1.81.1 / x/net v0.51.0 (PROPOSALS.md)
- **Blocked on:** nothing
