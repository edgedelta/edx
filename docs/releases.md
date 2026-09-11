# Releasing edx

Pushing a `v*` tag runs Windows tests and smoke checks, then GoReleaser builds
macOS, Linux, and Windows binaries for amd64 and arm64. Windows assets are ZIPs;
macOS/Linux assets remain tar.gz archives. All archives are covered by the
release's `checksums.txt`.

## Submit WinGet using your existing gh login

No new token is required for this route. The release workflow generates the
manifests and saves them as the `winget-manifests` Actions artifact. Once a
stable release containing the Windows ZIPs has finished:

1. Download its manifests using `gh run download <run-id> --repo edgedelta/edx
   --name winget-manifests --dir /tmp/edx-winget`. Use the release workflow run,
   not a snapshot: the manifest URLs must point to published release assets.
2. Create a fork with `gh repo fork microsoft/winget-pkgs --org edgedelta
   --clone=false --remote=false` (only needed once).
3. Clone the fork into a separate directory with `gh repo clone
   edgedelta/winget-pkgs -- --filter=blob:none --sparse`. Create a branch such
   as `edx-0.21.0`, enable sparse checkout for `manifests/e/EdgeDelta/edx`, and
   copy the downloaded `manifests/` tree into the checkout. Validate the version
   directory with `winget validate --manifest <directory>` on Windows.
4. Commit and push that branch to the fork, then submit it using:

   ```sh
   gh pr create --repo microsoft/winget-pkgs --base master \
     --head edgedelta:edx-0.21.0 \
     --title 'New package: EdgeDelta.edx version 0.21.0' \
     --body 'Adds the Edge Delta CLI Windows x64 and ARM64 portable packages.'
   ```

Substitute the actual release version. Follow the PR's validation/review results
through acceptance. The install command becomes available after acceptance and
indexing in the community source. Repeat submission for subsequent releases.

## Optional unattended WinGet submissions

This is only needed if GitHub Actions should submit PRs without a local `gh`
session. Homebrew already uses the existing `HOMEBREW_TAP_GITHUB_TOKEN` secret.

1. Fork `microsoft/winget-pkgs` into `edgedelta/winget-pkgs`. The configured
   upstream base branch is `master`.
2. Add the repository Actions secret `WINGET_GITHUB_TOKEN`. Use a release bot
   token that can push branches to that fork and open pull requests against
   `microsoft/winget-pkgs` (for example, a classic PAT with `public_repo` for
   these public repositories). Authorize it for the organization if required.
   The default workflow `GITHUB_TOKEN` cannot perform this cross-repository work.
3. Publish a stable release. GoReleaser generates the `EdgeDelta.edx` manifests,
   pushes branch `edx-<version>` to the fork, and opens an upstream PR. Follow the
   PR's validation/review results through acceptance. The install command only
   becomes available after acceptance and indexing in the community source.

Without the WinGet secret, ZIP/Homebrew releases still publish and WinGet manifests
are generated locally in `dist/winget`. The release workflow preserves those
files as the `winget-manifests` Actions artifact. Prereleases also skip WinGet
publication. Do not advertise WinGet availability before the first acceptance.

## Validation

```sh
goreleaser check
goreleaser release --snapshot --clean --skip=publish,announce,validate
```

Inspect the generated Windows ZIPs and `dist/winget` manifests. On Windows,
validate the directory containing the generated manifests using
`winget validate --manifest <directory>`. Before the first public release,
exercise `edx auth login` and skills installation on a Windows machine; CI
covers tests and non-interactive startup, not browser sign-in.

GoReleaser's [WinGet integration](https://goreleaser.com/customization/publish/winget/)
documents the fork/token configuration. Check the upstream PR on every release:
an unsuccessful PR creation may be logged without failing the release workflow.
