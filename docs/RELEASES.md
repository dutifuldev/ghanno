# Releases

This document describes how `prtags` should ship the `prtags` CLI through GitHub Releases.

## Goal

Ship versioned `prtags` binaries through GitHub Releases using semantic version tags.

The first release path should stay simple:

- GitHub Releases only
- SemVer tags
- automated multi-platform binary builds
- generated checksums
- a one-command installer script

This is enough to make the CLI installable without committing to Homebrew, npm, or other package ecosystems yet.

## Versioning

Use semantic version tags:

- `v0.1.0`
- `v0.2.0`
- `v1.0.0`

The tag is the release trigger.

The release workflow should only run on pushed version tags that match the `vX.Y.Z` shape.

## Initial Release Targets

The initial release matrix should be:

- Linux `amd64`
- Linux `arm64`
- macOS `amd64`
- macOS `arm64`
- Windows `amd64`

That covers the main user platforms without making the first release setup too large.

## Release Artifacts

Each GitHub Release should include:

- compressed `prtags` binaries for each supported OS and architecture
- a checksums file
- generated release notes or a changelog section

The binary name should stay `prtags`.

## Tooling

Use GoReleaser.

GoReleaser should:

- build the `prtags` binary
- package the release archives
- generate checksums
- publish the release to GitHub

The repository should add:

- `.goreleaser.yaml`
- `.github/workflows/release.yml`

## Workflow

The release flow should be:

1. merge the desired changes to `main`
2. create a SemVer tag such as `v0.1.0`
3. push the tag
4. GitHub Actions runs the release workflow
5. GoReleaser publishes the release artifacts to GitHub Releases

## Documentation

The README should include a short install section that points users to:

- the latest GitHub Release
- the supported platforms
- the expected binary name
- the one-command installer script
