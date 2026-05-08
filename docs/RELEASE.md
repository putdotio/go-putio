# Release

## Delivery Model

Every merge to `master` should already be releasable.

GitHub Actions owns releases for this repo. The workflow verifies the Go module first, then runs semantic-release on `master`.

The release lane:

- reads Conventional Commits since the latest `v*` tag
- calculates the next semantic version
- creates the git tag
- creates the GitHub Release and generated notes

Go modules are published by git tag, so this repo does not build or upload package artifacts. It does not use GoReleaser because there are no binaries to publish.

## Required Tokens

The release job uses the built-in GitHub Actions bot token:

- `GITHUB_TOKEN`

No custom secret is required. The job grants `contents: write` only so semantic-release can create tags and GitHub Releases in this repository.

## Versioning

Conventional Commits drive automated version selection:

- `fix:` creates a patch release
- `feat:` creates a minor release
- `feat!:` or `BREAKING CHANGE:` creates a major release
- `ci:`, `docs:`, `test:`, and `chore:` do not create a release by default

The current release branch is `master`.

## Local Checks

Before changing release wiring, validate the repo-local guardrails:

```bash
go test -v ./...
go test -race -count=1 ./...
golangci-lint run ./...
actionlint
```
