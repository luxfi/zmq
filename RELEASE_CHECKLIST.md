# Release Checklist for ZMQ4

## Pre-Release Verification

### Code Quality ✅
- [x] All code formatted (`go fmt ./...`)
- [x] No vet issues (`go vet ./...`)
- [x] go.mod is tidy (`go mod tidy`)
- [x] Build tags properly set in CZMQ files
- [x] No CZMQ imports in non-tagged files

### Testing ✅
- [x] Pure Go tests pass
- [x] Pure Go build works with CGO_ENABLED=0
- [x] CZMQ build works with tags (if available)
- [x] Examples build successfully
- [x] Benchmarks run without errors

### Documentation ✅
- [x] README.md updated with badges
- [x] BUILD_TAGS.md created
- [x] CZMQ_COMPATIBILITY.md created
- [x] QUICK_REFERENCE.md created
- [x] Examples documented

### Infrastructure ✅
- [x] Makefile targets for CZMQ testing
- [x] CI/CD workflow updated
- [x] Pre-commit hooks created
- [x] Validation scripts working

## Release Process

### 1. Final Verification
```bash
# Run CI simulation locally
./scripts/ci_simulation.sh

# Verify build isolation
./scripts/test_czmq_isolation.sh

# Run benchmarks (optional)
./scripts/benchmark_comparison.sh
```

### 2. Commit Changes
```bash
# Check status
git status

# Add all changes
git add -A

# Commit with descriptive message
git commit -m "feat: Isolate CZMQ compatibility layer with build tags

- Pure Go implementation as default (zero C dependencies)
- CZMQ compatibility behind 'czmq4' build tag
- Comprehensive testing and benchmarks
- Full documentation suite
- Developer tools and automation"
```

### 3. Push to GitHub
```bash
# Push to main branch
git push origin main

# Wait for CI to pass
# Check: https://github.com/luxfi/zmq/actions
```

### 4. Create Release Tag
```bash
# Generate release notes
./scripts/generate_release_notes.sh v4.2.2

# Create and push tag
git tag -a v4.2.2 -m "Release v4.2.2: CZMQ isolation and pure Go default"
git push origin v4.2.2
```

### 5. Create GitHub Release
1. Go to https://github.com/luxfi/zmq/releases/new
2. Select the tag: v4.2.2
3. Title: "v4.2.2: Pure Go Default with Optional CZMQ"
4. Copy content from RELEASE_NOTES_v4.2.2.md
5. Check "Set as the latest release"
6. Publish release

### 6. Post-Release
- [ ] Verify pkg.go.dev updated
- [ ] Monitor CI for any issues
- [ ] Update internal documentation if needed
- [ ] Announce to team/community

## Version Numbering

Current: v4.2.1
Next: v4.2.2

Justification for patch version:
- No API changes
- No protocol changes (still ZMQ v4.2.x)
- Build configuration improvement only
- Backward compatible
- Documentation and tooling enhancements

## Release Notes Summary

### 🎯 Key Features
- **Pure Go by Default**: Zero C dependencies
- **Optional CZMQ**: Build tag controlled compatibility
- **Comprehensive Testing**: Benchmarks, integration tests
- **Developer Tools**: Pre-commit hooks, automation scripts
- **Full Documentation**: Build guides, compatibility docs

### 🔄 Breaking Changes
None - API remains backward compatible

### 🆕 New Features
- Build tag isolation for CZMQ
- Performance benchmarks
- Pre-commit validation
- Release automation
- Comprehensive documentation

### 🐛 Bug Fixes
- Fixed import paths in examples
- Cleaned up test artifacts
- Resolved go.mod dependencies

### 📚 Documentation
- 8 new documentation files
- Updated README with badges
- Quick reference guide
- Build configuration guide

## Verification Commands

```bash
# Verify pure Go build
go build .

# Verify with CZMQ tag
go build -tags czmq4 .

# Run tests
make test

# Run CZMQ tests (if available)
make test-czmq

# Check isolation
./scripts/test_czmq_isolation.sh
```

## Contact for Issues

- GitHub Issues: https://github.com/luxfi/zmq/issues
- Documentation: https://pkg.go.dev/github.com/luxfi/zmq/v4

---

✅ **Ready for Release**: All checks passed, documentation complete, CI ready.