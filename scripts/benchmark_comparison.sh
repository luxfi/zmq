#!/bin/bash

# Benchmark comparison script for Pure Go vs CZMQ implementations

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "================================================"
echo "   ZMQ4 Performance Benchmark Comparison"
echo "================================================"
echo

# Check if CZMQ is available
if ! go build -tags czmq4 -o /dev/null . 2>/dev/null; then
    echo -e "${YELLOW}Warning: CZMQ not available. Install libczmq-dev to run comparison.${NC}"
    echo "Running Pure Go benchmarks only..."
    echo
fi

# Create results directory
RESULTS_DIR="benchmark_results"
mkdir -p "$RESULTS_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "Running benchmarks..."
echo "Results will be saved to: $RESULTS_DIR/benchmark_$TIMESTAMP.txt"
echo

# Run Pure Go benchmarks
echo -e "${BLUE}=== Pure Go Implementation ===${NC}"
echo "Running Pure Go benchmarks..."
go test -bench=BenchmarkPureGo -benchmem -benchtime=10s -count=3 . 2>/dev/null | tee "$RESULTS_DIR/pure_go_$TIMESTAMP.txt"

echo
echo -e "${BLUE}=== CZMQ Implementation ===${NC}"
if go build -tags czmq4 -o /dev/null . 2>/dev/null; then
    echo "Running CZMQ benchmarks..."
    go test -tags czmq4 -bench=BenchmarkCZMQ -benchmem -benchtime=10s -count=3 . 2>/dev/null | tee "$RESULTS_DIR/czmq_$TIMESTAMP.txt"
else
    echo -e "${YELLOW}Skipping CZMQ benchmarks (not available)${NC}"
fi

# Generate comparison report
echo
echo "================================================"
echo "           Benchmark Comparison Report"
echo "================================================"
echo

if [ -f "$RESULTS_DIR/pure_go_$TIMESTAMP.txt" ] && [ -f "$RESULTS_DIR/czmq_$TIMESTAMP.txt" ]; then
    echo "Comparing results..."
    
    # Extract and compare key metrics
    echo -e "${GREEN}Pure Go Performance:${NC}"
    grep "Benchmark" "$RESULTS_DIR/pure_go_$TIMESTAMP.txt" | grep -E "ns/op|B/op|allocs/op" | head -5
    
    echo
    echo -e "${GREEN}CZMQ Performance:${NC}"
    grep "Benchmark" "$RESULTS_DIR/czmq_$TIMESTAMP.txt" | grep -E "ns/op|B/op|allocs/op" | head -5
    
    # Create summary file
    cat > "$RESULTS_DIR/summary_$TIMESTAMP.txt" << EOF
Benchmark Comparison Summary
Generated: $(date)

Pure Go Results:
$(grep "Benchmark" "$RESULTS_DIR/pure_go_$TIMESTAMP.txt" | head -10)

CZMQ Results:
$(grep "Benchmark" "$RESULTS_DIR/czmq_$TIMESTAMP.txt" | head -10)

Notes:
- Lower ns/op is better (faster)
- Lower B/op is better (less memory)
- Lower allocs/op is better (fewer allocations)
EOF
    
    echo
    echo -e "${GREEN}Summary saved to: $RESULTS_DIR/summary_$TIMESTAMP.txt${NC}"
else
    echo "Only Pure Go results available:"
    grep "Benchmark" "$RESULTS_DIR/pure_go_$TIMESTAMP.txt" 2>/dev/null | head -10
fi

echo
echo "================================================"
echo "           Benchmark Complete"
echo "================================================"
echo
echo "To run specific benchmarks:"
echo "  Pure Go:  go test -bench=BenchmarkPureGo/PubSub ."
echo "  CZMQ:     go test -tags czmq4 -bench=BenchmarkCZMQ/PubSub ."
echo
echo "To compare specific patterns:"
echo "  ./scripts/benchmark_comparison.sh | grep -E '(PubSub|ReqRep)'"