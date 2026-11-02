# Proposals: launcher

Written by loopers. Each proposal: what, why, links. Jody promotes accepted ones into ROADMAP.md (via an agent);
the looper then marks it `promoted` or `declined` here with the date.

## Open
- **CI build and govulncheck gate** (2026-10-28, tacplus-dev/launcher#3): CI runs only the planning files check. The land gates now run `gates.sh` (build, vet, tidy, govulncheck; #4), but a push outside the looper and the nightly are unchecked; a CI job running `./gates.sh` on push and nightly would close that.
- **Bump fleet-seatrials' `launcher_commit` pin** (2026-11-02, tacplus-dev/launcher#4): `images/pins.yaml` pins c575e14, whose go.mod has grpc v1.81.1 and x/net v0.51.0 (the 7 vulnerabilities of #4); its build image builds the launcher from that source. Pin main at or after #1/#3 (go.mod grpc v1.84.0, x/net v0.62.0). A change in tacplus-dev/fleet-seatrials.
- **Release pipeline for the launcher** (2026-11-02, tacplus-dev/launcher#4, tacplus-dev/looper#17): a tag build producing static `launcher-linux-amd64`, `launcher-linux-arm64` and `checksums.txt` for the GitHub release, so kernel, eventbus and stratum can take a pinned, checksummed release instead of schematic's image.

## Promoted / declined
