# Versioning Guide

This document explains how to version and release `mapbot-shared`.

## 📏 Semantic Versioning

We follow [Semantic Versioning 2.0.0](https://semver.org/):

```
v<MAJOR>.<MINOR>.<PATCH>[-<PRERELEASE>]

Example: v1.2.3, v2.0.0-beta.1
```

### Version Components

- **MAJOR** (v**X**.0.0): Breaking changes that require code updates
- **MINOR** (vX.**Y**.0): New features, backward-compatible
- **PATCH** (vX.Y.**Z**): Bug fixes, backward-compatible
- **PRERELEASE**: alpha, beta, rc (e.g., v1.0.0-beta.1)

## 🚀 Release Process

### 1. Prepare Release

```bash
cd /Users/xavierthomas/Documents/dev/LLM/mapbot/mapbot-shared

# Make sure everything is committed
git status

# Pull latest changes
git pull origin main

# Run tests locally
go test ./...

# Run linter
golangci-lint run
```

### 2. Update Version

Determine the version number based on changes:

**Breaking changes (MAJOR):**
- Changed function signatures
- Removed public functions/types
- Changed behavior that breaks existing code

**New features (MINOR):**
- Added new functions/packages
- New optional parameters
- Enhanced existing features (backward-compatible)

**Bug fixes (PATCH):**
- Fixed bugs
- Performance improvements
- Documentation updates

### 3. Create and Push Tag

```bash
# For a new minor release (e.g., v0.2.0)
git tag -a v0.2.0 -m "Release v0.2.0

New features:
- Added connection pool metrics
- Improved error messages
- Enhanced test utilities

Bug fixes:
- Fixed connection leak in database manager
"

# Push the tag
git push origin v0.2.0
```

### 4. Verify Release

1. Check GitHub Actions: https://github.com/pixime/mapbot-shared/actions
2. Verify release created: https://github.com/pixime/mapbot-shared/releases
3. Test module availability:
   ```bash
   go list -m github.com/pixime/mapbot-shared@v0.2.0
   ```

## 📋 Version History

### v0.1.0 (Initial Release)
- Database manager with PostgreSQL/PostGIS support
- Configuration utilities
- Structured logging
- Test utilities with testcontainers
- CI/CD pipelines

### Future Versions

See [GitHub Releases](https://github.com/pixime/mapbot-shared/releases) for complete version history.

## 🔄 Updating Dependent Projects

After releasing a new version:

```bash
# In mapbot-ai
cd ../mapbot-ai
go get github.com/pixime/mapbot-shared@v0.2.0
go mod tidy
git add go.mod go.sum
git commit -m "Update mapbot-shared to v0.2.0"
git push

# In french-admin-etl
cd ../french-admin-etl
go get github.com/pixime/mapbot-shared@v0.2.0
go mod tidy
git add go.mod go.sum
git commit -m "Update mapbot-shared to v0.2.0"
git push
```

## 🏷️ Version Conventions

### Stable Releases

```
v1.0.0, v1.1.0, v1.2.3
```

Use for production-ready code.

### Pre-releases

```
v1.0.0-alpha.1    # Early testing
v1.0.0-beta.1     # Feature-complete, testing
v1.0.0-rc.1       # Release candidate, final testing
```

Use during development or testing phases.

### Development Versions

```bash
# Work in progress (use pseudo-versions)
go get github.com/pixime/mapbot-shared@main

# Specific commit
go get github.com/pixime/mapbot-shared@abcd1234
```

## 🔧 Version Management Tips

### Check Current Version

```bash
# Latest stable
go list -m github.com/pixime/mapbot-shared@latest

# All versions
go list -m -versions github.com/pixime/mapbot-shared
```

### Rollback to Previous Version

```bash
go get github.com/pixime/mapbot-shared@v0.1.0
go mod tidy
```

### Force Update

```bash
# Clear cache and update
go clean -modcache
go get -u github.com/pixime/mapbot-shared@latest
```

## 📝 Release Checklist

Before creating a new tag:

- [ ] All tests pass locally (`go test ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Documentation updated (README, CHANGELOG)
- [ ] Breaking changes documented
- [ ] Version number follows semantic versioning
- [ ] Commit message is descriptive
- [ ] Tag annotation describes changes

After creating tag:

- [ ] GitHub Actions workflows pass
- [ ] GitHub Release created automatically
- [ ] Module indexed by Go proxy
- [ ] Dependent projects updated
- [ ] Team notified of release

## 🆘 Common Issues

### Tag already exists

```bash
# Delete local tag
git tag -d v0.2.0

# Delete remote tag
git push origin --delete v0.2.0

# Recreate tag
git tag -a v0.2.0 -m "..."
git push origin v0.2.0
```

### Wrong version published

Cannot delete published versions from Go proxy. Options:

1. **Retract the version** (Go 1.16+):
   ```go
   // In go.mod
   retract v0.2.0 // Accidentally published wrong version
   ```

2. **Publish a new patch version** (recommended)

3. **Document in release notes** that version should not be used

## 📞 Questions?

See [CONTRIBUTING.md](CONTRIBUTING.md) or open an issue.
