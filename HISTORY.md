# History: launcher

Append-only: add entries at the end, never edit or remove one (CI enforces it). One line each:
`- <YYYY-MM-DD> <who: looper|agent|jody> <kind: issue|decision|release|deploy> <one line> (<links>)`

- 2026-10-28 looper issue deps: go.mod selects go1.27.0 (7 reachable stdlib vulnerabilities) and x/net v0.61.0 has 2 more (govulncheck) (#3, ee5a393)
- 2026-11-02 looper issue launcher: the committed binary is a darwin build with grpc v1.81.1 / x/net v0.51.0 (7 reachable vulnerabilities), not main's source (#4, 7a383c3; gates: full)
