# Homebrew Distribution for `edx`

**Date:** 2026-06-15
**Status:** Approved design, pending implementation
**Author:** Fatih Yildiz (with Claude Code)

## Goal

Distribute the Edge Delta `edx` CLI via Homebrew for both macOS and Linux,
modeled on [`datadog-labs/homebrew-pack`](https://github.com/datadog-labs/homebrew-pack).
Users should be able to run:

```bash
brew install edgedelta/pack/edx
```

on any supported OS/arch and get a working `edx` binary.

## Background

`edx` is a Go/Cobra CLI (`github.com/edgedelta/edx`). It builds with the version
injected at link time:

```
-ldflags "-X github.com/edgedelta/edx/internal/cli.Version=<version>"
```

`internal/cli/version.go` declares `var Version = "dev"` and `edx version`
prints `edx <Version>`.

The repo currently has **no GitHub releases, no CI, and no GoReleaser config**.
The Datadog tap installs **pre-built GoReleaser tarballs** from GitHub Releases
(e.g. `pup_1.1.0_Darwin_arm64.tar.gz`), referenced per OS/arch with a sha256.
Faithfully reproducing that model therefore requires building the release
pipeline as well as the tap.

## Decisions (from brainstorming)

| Decision | Choice |
|----------|--------|
| Scope | Both the release pipeline (in `edx`) **and** the tap repo |
| Formula updates | **Automatic** via GoReleaser's `brews:` block |
| Tap license | Apache-2.0 |
| Execution | Create + push the tap repo; commit pipeline changes to a **branch** in `edx` (no direct push to `main`) |
| Tap location | Sibling dir: `/Users/yildiz/workspace/gocode/src/github.com/edgedelta/homebrew-pack` |
| Tap brew-audit CI | Included |

## Part A — Release pipeline in the `edx` repo

Delivered on a feature branch (e.g. `homebrew-release-pipeline`), not pushed to `main`.

### A1. `.goreleaser.yaml`

- **`builds`**: one build of `edx`
  - `main: .`
  - `binary: edx`
  - `env: [CGO_ENABLED=0]`
  - `ldflags: -s -w -X github.com/edgedelta/edx/internal/cli.Version={{.Version}}`
  - `goos: [darwin, linux]`
  - `goarch: [amd64, arm64]`
- **`archives`**: format `tar.gz`, name template producing the Datadog/Homebrew
  convention — Title-case OS and `amd64`→`x86_64`:
  - `edx_<version>_Darwin_arm64.tar.gz`
  - `edx_<version>_Darwin_x86_64.tar.gz`
  - `edx_<version>_Linux_arm64.tar.gz`
  - `edx_<version>_Linux_x86_64.tar.gz`
- **`checksum`**: `checksums.txt` (sha256).
- **`changelog`**: enabled, grouping conventional-commit prefixes.
- **`brews:`** block — the auto-update mechanism:
  - `repository: { owner: edgedelta, name: homebrew-pack, branch: main }`
  - `directory: Formula`
  - `homepage: https://github.com/edgedelta/edx`
  - `description: "Edge Delta command line interface"`
  - `license: "Apache-2.0"`
  - `install`: `bin.install "edx"`
  - `test`: `assert_match version.to_s, shell_output("#{bin}/edx version")`
  - Authenticated with `HOMEBREW_TAP_GITHUB_TOKEN` (env, set in CI).
  - On each release, GoReleaser regenerates `Formula/edx.rb` (version, per-arch
    URLs, sha256) and commits it to the tap repo.

### A2. `.github/workflows/release.yml`

- **Trigger**: `push` on tags matching `v*`.
- **Permissions**: `contents: write`.
- **Steps**:
  1. `actions/checkout` with `fetch-depth: 0` (GoReleaser needs full history/tags).
  2. `actions/setup-go` (Go 1.23, matching `go.mod`).
  3. `goreleaser/goreleaser-action` running `release --clean`.
  4. Env: `GITHUB_TOKEN` (default, creates the release in `edx`) and
     `HOMEBREW_TAP_GITHUB_TOKEN` (secret, cross-repo push to the tap).

### A3. `LICENSE` (edx repo)

Add an Apache-2.0 `LICENSE` to `edx` so the source license matches the
`license "Apache-2.0"` the formula declares. Copyright Edge Delta.

### A4. `README.md` (edx repo)

Add a "Install via Homebrew" subsection to the existing Install section:

```bash
brew install edgedelta/pack/edx
# or
brew tap edgedelta/pack
brew install edx
```

## Part B — `edgedelta/homebrew-pack` tap repo

Created locally at the sibling path, then `gh repo create` + pushed.

### B1. `README.md`

Install instructions (direct, tap-then-install, and Brewfile), mirroring the
Datadog tap README. Title: "Edge Delta Homebrew Tap".

### B2. `LICENSE`

Apache-2.0, copyright Edge Delta.

### B3. `Formula/edx.rb`

A **seed** formula matching the exact shape GoReleaser will emit, with:
- `version "0.0.0"` placeholder and a header comment stating the file is
  auto-managed by GoReleaser and overwritten on each release.
- `on_macos` / `on_linux` blocks with `Hardware::CPU.arm?` branches and the
  four release-tarball URLs (placeholder sha256 values).
- `desc`, `homepage`, `license "Apache-2.0"`.
- `install` → `bin.install "edx"`.
- `test` → asserts `edx version` output.

This makes the repo structure reviewable immediately; the first real release
(`v0.1.0`) overwrites it with valid URLs and checksums.

### B4. `.github/workflows/tests.yml`

Minimal formula-quality CI on PRs and pushes:
- Runs on `ubuntu-latest` with Homebrew available.
- `brew style Formula/edx.rb` and `brew audit --tap edgedelta/pack edx`
  (audit may be `|| true` until the first real release exists, since audit
  checks live URLs).

## Manual setup steps (require the user's GitHub permissions)

Documented in the tap README and/or a setup note; not automatable here:

1. **Create a cross-repo token**: a fine-grained PAT (or classic PAT) with
   `contents: write` on `edgedelta/homebrew-pack`. Add it to the **`edx`** repo
   as the Actions secret `HOMEBREW_TAP_GITHUB_TOKEN`. The default `GITHUB_TOKEN`
   cannot push to another repository, so this is required for auto-publishing.
2. **Cut a release**:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```
   The release workflow then builds the tarballs, creates the GitHub Release,
   and updates `Formula/edx.rb` in the tap.

## End-to-end flow

```
git push origin v0.1.0
        │
        ▼
edx Actions: release.yml → goreleaser release --clean
        │
        ├─ builds 4 tarballs + checksums.txt
        ├─ creates GitHub Release on edgedelta/edx
        └─ regenerates Formula/edx.rb on edgedelta/homebrew-pack
                │
                ▼
        brew install edgedelta/pack/edx   (macOS + Linux, amd64 + arm64)
```

## Implementation note: Formula vs Cask

GoReleaser has **deprecated** Formula generation (`brews`) in favor of
`homebrew_casks`. However, **Homebrew Casks are macOS-only** — Linux has no
cask support. Because `edx` must install on both Linux and macOS, a **Formula**
is the correct and only cross-platform choice, so we keep `brews` despite the
deprecation warning. `goreleaser check` exits non-zero on the deprecation, but
`goreleaser release` runs fine and emits a working formula. Revisit if/when
Homebrew adds cask support on Linux. This rationale is also documented inline
in `.goreleaser.yaml`.

## Out of scope

- Submitting `edx` to `homebrew-core` (the official Homebrew repo).
- Windows / Scoop / other package managers.
- Signing / notarization of macOS binaries.
- Existing release history (this starts at `v0.1.0`).

## Verification

- `goreleaser check` validates `.goreleaser.yaml` locally.
- `goreleaser release --snapshot --clean` builds tarballs locally without
  publishing (sanity check of build matrix + archive names).
- `brew style` / `brew audit` validate the seed formula syntax.
- True end-to-end (`brew install`) is verifiable only after the manual token
  setup and first tag push.
