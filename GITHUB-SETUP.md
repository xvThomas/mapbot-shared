# GitHub Repository Setup Guide

This guide explains how to create and configure the `mapbot-shared` repository on GitHub.

## 📝 Prerequisites

- GitHub account with organization `pixime` (or personal account)
- Git installed locally
- Repository already initialized locally

## 🚀 Step 1: Create Repository on GitHub

### Option A: Using GitHub CLI (Recommended)

```bash
cd /Users/xavierthomas/Documents/dev/LLM/mapbot/mapbot-shared

# Create public repository
gh repo create xvThomas/mapbot-shared --public --source=. --remote=origin --push

# Or create private repository
gh repo create xvThomas/mapbot-shared --private --source=. --remote=origin --push
```

### Option B: Using GitHub Web Interface

1. Go to https://github.com/new (or https://github.com/organizations/pixime/repositories/new)
2. Fill in:
   - **Repository name:** `mapbot-shared`
   - **Description:** `Shared Go utilities for MapBot ecosystem - database, config, logging, testing`
   - **Visibility:** Public (recommended for Go modules) or Private
   - **DO NOT** initialize with README, .gitignore, or license (already done locally)
3. Click **Create repository**

## 🔗 Step 2: Link Local Repository to GitHub

If you created the repo via web interface:

```bash
cd /Users/xavierthomas/Documents/dev/LLM/mapbot/mapbot-shared

# Add remote
git remote add origin git@github.com:xvThomas/mapbot-shared.git

# Or if using HTTPS:
# git remote add origin https://github.com/xvThomas/mapbot-shared.git

# Set main branch
git branch -M main

# Push to GitHub
git push -u origin main
```

## 🏷️ Step 3: Create Initial Release

```bash
cd /Users/xavierthomas/Documents/dev/LLM/mapbot/mapbot-shared

# Create and push v0.1.0 tag
git tag -a v0.1.0 -m "Initial release

Features:
- Database manager with PostgreSQL/PostGIS support
- Configuration utilities
- Structured logging
- Test utilities with testcontainers
- CI/CD pipelines"

git push origin v0.1.0
```

This will trigger the `.github/workflows/release.yml` workflow and create a GitHub release automatically.

## ⚙️ Step 4: Configure Repository Settings

### Branch Protection (Recommended)

1. Go to repository settings → Branches
2. Add branch protection rule for `main`:
   - ✅ Require a pull request before merging
   - ✅ Require status checks to pass before merging
     - Select: `golangci-lint`, `go fmt`, `go vet`, `Test`
   - ✅ Require branches to be up to date before merging
   - ✅ Do not allow bypassing the above settings

### Secrets (If using private features)

1. Go to repository settings → Secrets and variables → Actions
2. Add secrets if needed (e.g., `CODECOV_TOKEN`)

### Topics (For discoverability)

1. Go to repository main page
2. Click gear icon next to "About"
3. Add topics:
   - `go`
   - `golang`
   - `postgresql`
   - `database`
   - `testing`
   - `utilities`
   - `mapbot`

## 📊 Step 5: Verify CI/CD

After pushing, verify that GitHub Actions workflows run successfully:

1. Go to repository → Actions tab
2. Check that all workflows pass:
   - ✅ Lint
   - ✅ Test
   - ✅ Release (after creating tag)

## 🔐 Step 6: Make Module Publicly Available

For public Go modules:

```bash
# Trigger Go proxy to index the module
curl -X POST "https://proxy.golang.org/github.com/xvThomas/mapbot-shared/@v/v0.1.0.info"

# Verify it's available
curl "https://proxy.golang.org/github.com/xvThomas/mapbot-shared/@latest"
```

## ✅ Verification Checklist

- [ ] Repository created on GitHub
- [ ] Local repository pushed to GitHub
- [ ] Initial tag v0.1.0 created and pushed
- [ ] GitHub Actions workflows passing
- [ ] Branch protection rules configured (optional but recommended)
- [ ] Module indexed by Go proxy
- [ ] README displays correctly on GitHub
- [ ] License file present

## 📝 Next Steps

After repository is set up:

1. **Update dependent projects:**

   ```bash
   cd ../mapbot-ai
   go get github.com/xvThomas/mapbot-shared@v0.1.0
   ```

2. **Create Go Workspace** (see main WORKSPACE-SETUP.md)

3. **Migrate code** from mapbot-ai and french-admin-etl to use mapbot-shared

## 🆘 Troubleshooting

### Error: remote origin already exists

```bash
git remote remove origin
git remote add origin git@github.com:xvThomas/mapbot-shared.git
```

### Error: failed to push some refs

```bash
# If remote has commits not in local (shouldn't happen for new repo)
git pull origin main --rebase
git push -u origin main
```

### Go module not found

Wait a few minutes for Go proxy to index, or force update:

```bash
GOPROXY=direct go get github.com/xvThomas/mapbot-shared@v0.1.0
```

## 📚 Additional Resources

- [Go Modules Documentation](https://go.dev/doc/modules)
- [GitHub Actions Documentation](https://docs.github.com/actions)
- [Semantic Versioning](https://semver.org/)
