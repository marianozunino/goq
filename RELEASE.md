# Release Process

This document describes the simplified release process for GOQ.

## Quick Release (Recommended)

### Using Task (if you have Task installed)

```bash
# Create a patch release (v1.2.3 -> v1.2.4)
task release-patch

# Create a minor release (v1.2.3 -> v1.3.0)
task release-minor

# Create a major release (v1.2.3 -> v2.0.0)
task release-major

# Create a specific version release
VERSION=v1.5.0 task release
```


## What Happens During Release

1. **Pre-release Checks**: The script runs all quality checks (format, vet, test)
2. **Version Validation**: Ensures the version follows semantic versioning (v1.2.3)
3. **Clean Working Directory**: Verifies no uncommitted changes exist
4. **Tag Creation**: Creates an annotated git tag with the release version
5. **Tag Push**: Pushes the tag to the remote repository
6. **Automated Build**: GitHub Actions automatically builds and publishes the release

## Manual Release Process (if needed)

If you prefer to do it manually:

```bash
# 1. Ensure working directory is clean
git status

# 2. Run quality checks
task check

# 3. Create and push tag
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

## Release Automation

The release process is fully automated through:

- **GitHub Actions** (`.github/workflows/releaser.yml`): Triggers on tag push
- **GoReleaser** (`.goreleaser.yaml`): Builds binaries for multiple platforms
- **UPX Compression**: Reduces binary size
- **Multi-platform Support**: Linux, macOS, and Windows

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0): Breaking changes
- **MINOR** (v1.1.0): New features, backward compatible
- **PATCH** (v1.0.1): Bug fixes, backward compatible

## Troubleshooting

### Release fails with "working directory not clean"
```bash
git add .
git commit -m "Prepare for release"
# Then retry the release
```

### Release fails with test errors
```bash
# Fix failing tests first
task test
# Then retry the release
```

### Need to delete a release tag
```bash
git tag -d v1.2.3
git push origin --delete v1.2.3
```
