# Proposals: launcher

Written by loopers. Each proposal: what, why, links. Jody promotes accepted ones into ROADMAP.md (via an agent);
the looper then marks it `promoted` or `declined` here with the date.

## Open
- **CI build and govulncheck gate** (2026-10-28, tacplus-dev/launcher#3): CI runs only the planning files check. The land gates now run `gates.sh` (build, vet, tidy, govulncheck; #4), but a push outside the looper and the nightly are unchecked; a CI job running `./gates.sh` on push and nightly would close that.
- **Bump fleet-seatrials' `launcher_commit` pin** (2026-11-02, tacplus-dev/launcher#4): `images/pins.yaml` pins c575e14, whose go.mod has grpc v1.81.1 and x/net v0.51.0 (the 7 vulnerabilities of #4); its build image builds the launcher from that source. Pin main at or after #1/#3 (go.mod grpc v1.84.0, x/net v0.62.0). A change in tacplus-dev/fleet-seatrials.
- **Release assets from the deterministic build** (2026-11-07, tacplus-dev/launcher#5, tacplus-dev/looper#17): the release workflow (#5) builds the assets on a GitHub runner with `-trimpath` and the go.mod toolchain (the main module's version comes from the tag, so a reproduction needs the tag fetched); once looper#17's pipeline exists, build (or reproduce and compare) them there and attach its attestation.

- **Point eventbus' and stratum's health probe at the plaintext port before they take #7** (2026-11-09, tacplus-dev/launcher#7): their Dockerfiles set `HEALTH_PORT=8443`, the chassis TLS port. Since #7 the probe never skips verification: without `HEALTH_CA` it is plaintext, fails on a TLS port and kills the child after `HEALTH_TIMEOUT`. Set `HEALTH_PORT=8444` (the chassis `health_probe_port`) in tacplus-dev/eventbus and tacplus-dev/stratum before, or with, the launcher pin that takes #7 (kernel sets no `HEALTH_PORT`: it probes 5000, check what it serves there). Changes in those repositories.
- **Verify the health probe's peer by SPIFFE ID** (2026-11-09, tacplus-dev/launcher#7): `HEALTH_CA` verifies the chain and a DNS name (`HEALTH_SERVER_NAME`); a chassis certificate carries a SPIFFE URI SAN, which crypto/tls's hostname check does not match. Verifying the URI SAN (or reading the SPIFFE bundle the service carries) would let the TLS probe use the service's own certificate as issued.
- **Scope the G702/G703 exceptions with `where`** (2026-11-23, tacplus-dev/launcher#11): the security report of 54e9674 warns that both exceptions in `security-exceptions.yaml` have no `where`, so they accept G702/G703 in every file, not only the findings they were decided for (main.go:81, :217, :240). Scoping them to `main.go` (or those call sites) keeps a new G702/G703 elsewhere failing.

## Promoted / declined
- **Release pipeline for the launcher** (2026-11-02, tacplus-dev/launcher#4): promoted 2026-11-06 (#5, approved by Jody); `release.sh` and the tag workflow `.github/workflows/release.yml` landed (#5); v0.1.0 released 2026-11-15 from 0e139ca, its assets reproduced bit for bit by a local `release.sh` of the tag.
