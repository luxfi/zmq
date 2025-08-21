#!/bin/bash

# Install git hooks for the zmq4 project

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "================================"
echo "   Git Hooks Installation"
echo "================================"
echo

# Check if we're in a git repository
if [ ! -d .git ]; then
    echo -e "${YELLOW}Error: Not in a git repository${NC}"
    echo "Please run this script from the repository root."
    exit 1
fi

# Set git hooks path
echo "Configuring git to use .githooks directory..."
git config core.hooksPath .githooks

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Git hooks installed successfully${NC}"
    echo
    echo "Hooks installed:"
    echo "  - pre-commit: Validates build tags and CZMQ isolation"
    echo
    echo "The pre-commit hook will:"
    echo "  1. Check build tags in CZMQ files"
    echo "  2. Ensure no CZMQ imports in non-tagged files"
    echo "  3. Verify pure Go build works"
    echo "  4. Check code formatting"
    echo "  5. Verify go.mod is tidy"
    echo "  6. Check for test files"
    echo
    echo "To skip hooks temporarily (not recommended):"
    echo "  git commit --no-verify"
    echo
    echo "To uninstall hooks:"
    echo "  git config --unset core.hooksPath"
else
    echo -e "${YELLOW}Failed to install git hooks${NC}"
    exit 1
fi