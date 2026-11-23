# History: launcher

Append-only: add entries at the end, never edit or remove one (CI enforces it). One line each:
`- <YYYY-MM-DD> <who: looper|agent|jody> <kind: issue|decision|release|deploy> <one line> (<links>)`

- 2026-10-28 looper issue deps: go.mod selects go1.27.0 (7 reachable stdlib vulnerabilities) and x/net v0.61.0 has 2 more (govulncheck) (#3, ee5a393)
- 2026-11-02 looper issue launcher: the committed binary is a darwin build with grpc v1.81.1 / x/net v0.51.0 (7 reachable vulnerabilities), not main's source (#4, 7a383c3; gates: full)
- 2026-11-15 looper issue gates fail on main: govulncheck finds GO-2026-6859 (net/http) in go.mod's toolchain go1.27.3, fixed in go1.27.4 (#8, ebc2151; gates: full)
- 2026-11-15 looper issue Build security floor: three HIGH gosec findings in main.go (G402, G702, G703) need a fix or an exception (#7)
- 2026-11-15 looper release v0.1.0 (https://github.com/tacplus-dev/launcher/releases/tag/v0.1.0)
- 2026-11-15 looper issue Pin the build: a committed build.lock naming this repository's build image by digest (#6)
- 2026-11-15 looper issue Release: v0.1.0 with static launcher-linux-{amd64,arm64} and checksums.txt, so kernel, eventbus and stratum can pin it (#5)
- 2026-11-15 looper issue STATE.md: Blocked on names tacplus-dev/looper#106, closed 2026-11-08; #5 now waits on tacplus-dev/looper#113 (#9)
- 2026-11-22 looper issue Pinned images behind their lines: relock (build.lock: ghcr.io/tacplus-dev/tacplus/build:2.6.12 → ghcr.io/tacplus-dev/tacplus/build:2.6.13) (#10)
