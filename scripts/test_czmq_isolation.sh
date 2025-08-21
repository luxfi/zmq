#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "================================================"
echo "       CZMQ Build Tag Isolation Test"
echo "================================================"
echo

# Function to print colored output
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
        exit 1
    fi
}

# Test 1: Verify pure Go build works without czmq
echo "Test 1: Pure Go Build (no czmq dependency)"
echo "-------------------------------------------"
GO111MODULE=on go build -o /tmp/zmq4_pure . 2>/dev/null
print_result $? "Pure Go build successful"
rm -f /tmp/zmq4_pure

# Test 2: Check that cxx_zmq4_compat.go is excluded by default
echo
echo "Test 2: Build Tag Exclusion"
echo "----------------------------"
FILES=$(go list -f '{{.GoFiles}}' . | grep -c "cxx_zmq4_compat.go")
if [ "$FILES" -eq 0 ]; then
    print_result 0 "cxx_zmq4_compat.go correctly excluded without tag"
else
    print_result 1 "cxx_zmq4_compat.go incorrectly included"
fi

# Test 3: Check that cxx_zmq4_compat.go is included with tag
echo
echo "Test 3: Build Tag Inclusion"
echo "----------------------------"
FILES=$(go list -f '{{.GoFiles}}' -tags=czmq4 . | grep -c "cxx_zmq4_compat.go")
if [ "$FILES" -eq 1 ]; then
    print_result 0 "cxx_zmq4_compat.go correctly included with czmq4 tag"
else
    print_result 1 "cxx_zmq4_compat.go not included with tag"
fi

# Test 4: Verify test files are properly tagged
echo
echo "Test 4: Test File Isolation"
echo "----------------------------"
TEST_FILES=$(go list -f '{{.XTestGoFiles}}' . | grep -c "czmq4_test.go")
if [ "$TEST_FILES" -eq 0 ]; then
    print_result 0 "czmq4_test.go correctly excluded without tag"
else
    print_result 1 "czmq4_test.go incorrectly included"
fi

TEST_FILES=$(go list -f '{{.XTestGoFiles}}' -tags=czmq4 . | grep -c "czmq4_test.go")
if [ "$TEST_FILES" -eq 1 ]; then
    print_result 0 "czmq4_test.go correctly included with czmq4 tag"
else
    print_result 1 "czmq4_test.go not included with tag"
fi

# Test 5: Verify package imports
echo
echo "Test 5: Import Dependencies"
echo "----------------------------"
IMPORTS=$(go list -f '{{.Imports}}' . | grep -c "czmq")
if [ "$IMPORTS" -eq 0 ]; then
    print_result 0 "No czmq imports in default build"
else
    print_result 1 "Found czmq imports in default build"
fi

# Test 6: Build with czmq4 tag (will fail if czmq not installed, which is expected)
echo
echo "Test 6: Build with CZMQ Tag"
echo "----------------------------"
GO111MODULE=on go build -tags czmq4 -o /tmp/zmq4_czmq . 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Build with czmq4 tag successful (czmq installed)"
    rm -f /tmp/zmq4_czmq
else
    echo -e "${YELLOW}⚠${NC} Build with czmq4 tag failed (czmq not installed - this is OK)"
fi

# Test 7: Verify Makefile targets exist
echo
echo "Test 7: Makefile Targets"
echo "------------------------"
if make help | grep -q "test-czmq"; then
    print_result 0 "Makefile has czmq test targets"
else
    print_result 1 "Makefile missing czmq test targets"
fi

# Summary
echo
echo "================================================"
echo -e "${GREEN}All isolation tests passed!${NC}"
echo "================================================"
echo
echo "Summary:"
echo "• Pure Go build works without czmq dependency"
echo "• Build tags correctly isolate czmq code"
echo "• Test files are properly tagged"
echo "• No czmq imports in default build"
echo "• Makefile includes czmq test targets"
echo
echo "To run tests with czmq compatibility:"
echo "  make test-czmq"
echo
echo "To run regular tests (pure Go):"
echo "  make test"