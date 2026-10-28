# Proposals: launcher

Written by loopers. Each proposal: what, why, links. Jody promotes accepted ones into ROADMAP.md (via an agent);
the looper then marks it `promoted` or `declined` here with the date.

## Open
- **Remove the committed macOS binary `launcher`** (2026-10-28, tacplus-dev/launcher#3): the repository tracks `launcher`, a Mach-O arm64 executable, likely an accidental local build from the copy out of schematic. It is not the Linux static binary the images use, bloats history and is unverifiable provenance; delete it and add it to `.gitignore`.
- **CI build and govulncheck gate** (2026-10-28, tacplus-dev/launcher#3): CI runs only the planning files check. A workflow running `go build ./...`, `go vet ./...`, `go mod tidy -diff` and `govulncheck ./...` on push and nightly would have caught #3 at the commit that raised the floor without a patched `toolchain` line.

## Promoted / declined
