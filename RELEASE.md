# 🚀 Fuselage Release Guide

This document outlines the process for releasing new versions of Fuselage.

## 📋 Version Types

### Semantic Versioning (SemVer)
- **Patch (v1.0.1)** - Bug fixes, documentation updates
- **Minor (v1.1.0)** - New features, backward compatible
- **Major (v2.0.0)** - Breaking changes, API changes

## 🔄 Release Process

### 1. Pre-Release Checklist
```bash
# Run all tests
go test -v ./...

# Check code formatting
go fmt ./...

# Build verification
go build ./...

# Lint check (if available)
golangci-lint run
```

### 2. Update Documentation

#### Update CHANGELOG.md
```markdown
## [vX.Y.Z] - YYYY-MM-DD

### Added
- 🆕 New feature description

### Changed
- 🔄 Modified functionality

### Fixed
- 🐛 Bug fixes

### Removed
- 🗑️ Deprecated features
```

#### Update README.md
```bash
# Update version references
sed -i 's/v1.0.0/vX.Y.Z/g' README.md
```

#### Update Version History
Add new version to README.md versioning section:
```markdown
### Version History
- **vX.Y.Z** - Brief description of changes
- **v1.0.0** - Initial stable release
```

### 3. Commit Changes
```bash
git add .
git commit -m "Prepare release vX.Y.Z"
git push origin main
```

### 4. Create and Push Tag
```bash
# Create annotated tag
git tag -a vX.Y.Z -m "Release vX.Y.Z - Brief description"

# Push tag (triggers automatic release)
git push origin vX.Y.Z
```

### 5. Verify Release
- Check GitHub Actions workflow completion
- Verify GitHub Releases page
- Test installation: `go get github.com/k-tsurumaki/fuselage@vX.Y.Z`

## 📝 Release Templates

### Patch Release (Bug Fixes)
```bash
# Example: v1.0.1
git tag -a v1.0.1 -m "Release v1.0.1 - Bug fixes and improvements"
git push origin v1.0.1
```

### Minor Release (New Features)
```bash
# Example: v1.1.0
git tag -a v1.1.0 -m "Release v1.1.0 - New middleware and features"
git push origin v1.1.0
```

### Major Release (Breaking Changes)
```bash
# Example: v2.0.0
git tag -a v2.0.0 -m "Release v2.0.0 - Major API changes"
git push origin v2.0.0
```

## 🎯 Quick Release Commands

### For Patch Release
```bash
# Update CHANGELOG.md and README.md first, then:
git add . && git commit -m "Release v1.0.1"
git tag v1.0.1 && git push origin main && git push origin v1.0.1
```

### For Minor Release
```bash
# Update CHANGELOG.md and README.md first, then:
git add . && git commit -m "Release v1.1.0"
git tag v1.1.0 && git push origin main && git push origin v1.1.0
```

## ✅ Post-Release Verification

### 1. GitHub Releases
Visit: https://github.com/k-tsurumaki/fuselage/releases
- Verify new release is published
- Check release notes are generated correctly

### 2. Go Module Availability
```bash
# Check available versions
go list -m -versions github.com/k-tsurumaki/fuselage

# Test installation
mkdir test-release && cd test-release
go mod init test
go get github.com/k-tsurumaki/fuselage@vX.Y.Z
```

### 3. Documentation
- Verify pkg.go.dev updates (may take a few minutes)
- Check that examples work with new version

## 🚨 Rollback Process

If a release has issues:

### 1. Delete Tag
```bash
# Delete local tag
git tag -d vX.Y.Z

# Delete remote tag
git push origin :refs/tags/vX.Y.Z
```

### 2. Fix Issues and Re-release
```bash
# Fix the issues, then create new tag
git tag -a vX.Y.Z -m "Release vX.Y.Z - Fixed issues"
git push origin vX.Y.Z
```

## 🤖 Automated Process

The following are automated via GitHub Actions:
- ✅ Test execution on tag push
- ✅ GitHub Release creation
- ✅ Release notes generation
- ✅ Go module publication

## 📋 Release Checklist

- [ ] All tests pass
- [ ] CHANGELOG.md updated
- [ ] README.md version references updated
- [ ] Version history updated
- [ ] Changes committed and pushed
- [ ] Tag created and pushed
- [ ] GitHub Release verified
- [ ] Go module installation tested
- [ ] Documentation updated on pkg.go.dev

## 🎉 Release Complete!

After following these steps, your new version of Fuselage will be:
- Available via `go get`
- Published on GitHub Releases
- Documented on pkg.go.dev
- Ready for users to adopt

---

**Remember**: Always test thoroughly before releasing, and follow semantic versioning principles!