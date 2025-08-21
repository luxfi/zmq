#!/bin/bash

# Automated release notes generator for zmq4
# Usage: ./scripts/generate_release_notes.sh [version]

VERSION=${1:-"next"}
SINCE_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
DATE=$(date +"%Y-%m-%d")

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "================================================"
echo "    ZMQ4 Release Notes Generator"
echo "================================================"
echo

# Determine commit range
if [ -z "$SINCE_TAG" ]; then
    echo "No previous tags found. Generating notes for all commits."
    COMMIT_RANGE=""
else
    echo "Generating notes since: $SINCE_TAG"
    COMMIT_RANGE="$SINCE_TAG..HEAD"
fi

# Start generating release notes
cat > "RELEASE_NOTES_${VERSION}.md" << EOF
# Release Notes - v${VERSION}

**Date:** ${DATE}

## 🎯 Highlights

### Pure Go Implementation (Default)
- Zero C dependencies
- Cross-platform compatibility
- Easy deployment

### Optional CZMQ Compatibility
- Build tag: \`czmq4\`
- Full protocol compliance testing
- Performance benchmarking support

## 📊 Build Configuration

| Configuration | Build Command | Description |
|--------------|---------------|-------------|
| Pure Go | \`go build .\` | Default, no dependencies |
| With CZMQ | \`go build -tags czmq4 .\` | Compatibility testing |

## 🚀 What's New

EOF

# Add commit summary
echo "### Commits Since Last Release" >> "RELEASE_NOTES_${VERSION}.md"
echo "" >> "RELEASE_NOTES_${VERSION}.md"

if [ -n "$COMMIT_RANGE" ]; then
    git log $COMMIT_RANGE --pretty=format:"- %s (%an)" >> "RELEASE_NOTES_${VERSION}.md"
else
    git log --pretty=format:"- %s (%an)" --max-count=20 >> "RELEASE_NOTES_${VERSION}.md"
fi

echo "" >> "RELEASE_NOTES_${VERSION}.md"
echo "" >> "RELEASE_NOTES_${VERSION}.md"

# Add feature breakdown
cat >> "RELEASE_NOTES_${VERSION}.md" << EOF

## ✨ Features

### Build System
- CZMQ compatibility layer isolated with \`czmq4\` build tag
- Pure Go remains the default implementation
- Enhanced Makefile with CZMQ test targets

### Testing
- Comprehensive integration tests for CZMQ compatibility
- Performance benchmarks comparing implementations
- Validation scripts for build isolation

### Documentation
- BUILD_TAGS.md - Complete build constraints guide
- CZMQ_COMPATIBILITY.md - Compatibility documentation
- QUICK_REFERENCE.md - Developer quick reference
- Updated examples with compatibility demonstrations

## 🔧 Changes by Category

EOF

# Categorize changes
echo "### Added" >> "RELEASE_NOTES_${VERSION}.md"
git log $COMMIT_RANGE --pretty=format:"%s" 2>/dev/null | grep -i "^add\|^feat" | sed 's/^/- /' >> "RELEASE_NOTES_${VERSION}.md" || echo "- No additions" >> "RELEASE_NOTES_${VERSION}.md"

echo "" >> "RELEASE_NOTES_${VERSION}.md"
echo "### Changed" >> "RELEASE_NOTES_${VERSION}.md"
git log $COMMIT_RANGE --pretty=format:"%s" 2>/dev/null | grep -i "^update\|^change\|^modify" | sed 's/^/- /' >> "RELEASE_NOTES_${VERSION}.md" || echo "- No changes" >> "RELEASE_NOTES_${VERSION}.md"

echo "" >> "RELEASE_NOTES_${VERSION}.md"
echo "### Fixed" >> "RELEASE_NOTES_${VERSION}.md"
git log $COMMIT_RANGE --pretty=format:"%s" 2>/dev/null | grep -i "^fix\|^bug" | sed 's/^/- /' >> "RELEASE_NOTES_${VERSION}.md" || echo "- No fixes" >> "RELEASE_NOTES_${VERSION}.md"

# Add file statistics
echo "" >> "RELEASE_NOTES_${VERSION}.md"
echo "## 📈 Statistics" >> "RELEASE_NOTES_${VERSION}.md"
echo "" >> "RELEASE_NOTES_${VERSION}.md"

# Count changes
if [ -n "$COMMIT_RANGE" ]; then
    COMMITS=$(git rev-list --count $COMMIT_RANGE)
    FILES_CHANGED=$(git diff --name-only $COMMIT_RANGE | wc -l)
    INSERTIONS=$(git diff --stat $COMMIT_RANGE | tail -1 | grep -oE '[0-9]+ insertion' | grep -oE '[0-9]+' || echo "0")
    DELETIONS=$(git diff --stat $COMMIT_RANGE | tail -1 | grep -oE '[0-9]+ deletion' | grep -oE '[0-9]+' || echo "0")
else
    COMMITS=$(git rev-list --count HEAD)
    FILES_CHANGED=$(git ls-files | wc -l)
    INSERTIONS="N/A"
    DELETIONS="N/A"
fi

cat >> "RELEASE_NOTES_${VERSION}.md" << EOF
- Commits: $COMMITS
- Files Changed: $FILES_CHANGED
- Insertions: $INSERTIONS
- Deletions: $DELETIONS

## 🧪 Testing

### Run Tests
\`\`\`bash
# Pure Go tests
make test

# CZMQ compatibility tests
make test-czmq

# Performance benchmarks
./scripts/benchmark_comparison.sh
\`\`\`

### Validation
\`\`\`bash
# Verify build isolation
./scripts/test_czmq_isolation.sh
\`\`\`

## 📦 Installation

### Pure Go (Recommended)
\`\`\`bash
go get github.com/luxfi/zmq/v4@v${VERSION}
\`\`\`

### With CZMQ Support
\`\`\`bash
# Install CZMQ first
sudo apt-get install libczmq-dev  # Ubuntu/Debian
brew install czmq                 # macOS

# Then get the package
go get github.com/luxfi/zmq/v4@v${VERSION}
\`\`\`

## 🔄 Migration Guide

If upgrading from a previous version:

1. The default build now uses pure Go (no C dependencies)
2. CZMQ compatibility requires explicit \`czmq4\` build tag
3. API remains unchanged for backward compatibility

## 📚 Documentation

- [README.md](README.md) - Main documentation
- [BUILD_TAGS.md](BUILD_TAGS.md) - Build configuration details
- [CZMQ_COMPATIBILITY.md](CZMQ_COMPATIBILITY.md) - Compatibility guide
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Quick reference

## 🤝 Contributors

EOF

# Add contributors
if [ -n "$COMMIT_RANGE" ]; then
    git log $COMMIT_RANGE --pretty=format:"%an" | sort -u | sed 's/^/- /' >> "RELEASE_NOTES_${VERSION}.md"
else
    echo "- All contributors" >> "RELEASE_NOTES_${VERSION}.md"
fi

cat >> "RELEASE_NOTES_${VERSION}.md" << EOF

## 📝 License

BSD-3-Clause - See [LICENSE](LICENSE) file for details.

---

*Generated on ${DATE}*
EOF

echo -e "${GREEN}✓ Release notes generated: RELEASE_NOTES_${VERSION}.md${NC}"
echo
echo "To publish this release:"
echo "1. Review and edit RELEASE_NOTES_${VERSION}.md"
echo "2. Create a git tag: git tag -a v${VERSION} -m \"Release v${VERSION}\""
echo "3. Push the tag: git push origin v${VERSION}"
echo "4. Create GitHub release with these notes"