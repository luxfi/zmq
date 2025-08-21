#!/bin/bash

# CI Simulation Script
# Simulates the GitHub Actions CI pipeline locally

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "================================================"
echo "         CI Pipeline Simulation"
echo "================================================"
echo

FAILED=0

# Function to run a test and report results
run_test() {
    local name=$1
    local cmd=$2
    echo -e "${BLUE}Running: $name${NC}"
    if eval "$cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ $name passed${NC}"
    else
        echo -e "${RED}✗ $name failed${NC}"
        FAILED=1
    fi
    echo
}

# 1. Verify formatting (from lint job)
run_test "Format check" "test -z \$(gofmt -l .)"

# 2. Run go vet (from lint job)
run_test "Go vet" "go vet ./..."

# 3. Check go mod tidy (from lint job)
cp go.mod go.mod.backup
cp go.sum go.sum.backup 2>/dev/null || touch go.sum.backup
go mod tidy 2>/dev/null
if diff -q go.mod go.mod.backup >/dev/null && diff -q go.sum go.sum.backup >/dev/null; then
    echo -e "${GREEN}✓ go.mod is tidy${NC}"
else
    echo -e "${YELLOW}⚠ go.mod was modified by tidy${NC}"
    FAILED=1
fi
mv go.mod.backup go.mod
mv go.sum.backup go.sum 2>/dev/null || rm -f go.sum.backup
echo

# 4. Test pure Go build (from test-pure-go job)
echo -e "${BLUE}Testing Pure Go implementation${NC}"
run_test "Pure Go build" "CGO_ENABLED=0 go build -o /dev/null ."

# 5. Run pure Go tests
run_test "Pure Go tests" "CGO_ENABLED=0 go test -short -race ./..."

# 6. Test with CGO (from test-cgo job)
echo -e "${BLUE}Testing CGO implementation${NC}"
run_test "CGO build" "CGO_ENABLED=1 go build -o /dev/null ."

# 7. Test with CZMQ tag if available
echo -e "${BLUE}Testing CZMQ compatibility${NC}"
if go build -tags czmq4 -o /dev/null . 2>/dev/null; then
    run_test "CZMQ build" "go build -tags czmq4 -o /dev/null ."
    echo -e "${GREEN}✓ CZMQ compatibility available${NC}"
else
    echo -e "${YELLOW}⚠ CZMQ not available (this is OK)${NC}"
fi
echo

# 8. Build examples (from examples job)
echo -e "${BLUE}Building examples${NC}"
EXAMPLE_FAILED=0
for example in example/*.go; do
    if [ "$example" = "example/compatibility_test.go" ]; then
        # Skip CZMQ-only example
        continue
    fi
    if ! go build -o /dev/null "$example" 2>/dev/null; then
        echo -e "${RED}✗ Failed to build $example${NC}"
        EXAMPLE_FAILED=1
    fi
done
if [ $EXAMPLE_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All examples build successfully${NC}"
else
    FAILED=1
fi
echo

# Summary
echo "================================================"
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}       All CI checks passed!${NC}"
    echo "================================================"
    echo
    echo "Ready for release. Next steps:"
    echo "1. Commit all changes"
    echo "2. Push to GitHub"
    echo "3. Wait for CI to pass"
    echo "4. Create release tag"
    exit 0
else
    echo -e "${RED}       Some CI checks failed${NC}"
    echo "================================================"
    echo
    echo "Please fix the issues above before pushing."
    exit 1
fi