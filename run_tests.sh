#!/bin/bash

# 🧪 Comprehensive Test Runner for E-Commerce API
# Runs tests with coverage analysis and generates reports

set -e

echo "🚀 Starting E-Commerce API Test Suite..."
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Clean previous artifacts
echo -e "${BLUE}🧹 Cleaning previous test artifacts...${NC}"
rm -f coverage.out coverage.html *.prof

# Run tests with coverage
echo -e "${BLUE}🧪 Running tests with coverage...${NC}"
if go test -coverprofile=coverage.out ./api/... ./cache/... ./config/... ./models/... ./services/... ./tools/... ./utils/... -timeout=60s; then
    echo -e "${GREEN}✅ Tests completed successfully${NC}"
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi

# Generate coverage report
if [ -f coverage.out ]; then
    echo -e "${BLUE}📊 Generating coverage reports...${NC}"
    
    # Generate HTML report
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}📄 HTML coverage report: coverage.html${NC}"
    
    # Display coverage summary
    echo -e "${BLUE}� Coverage Summary:${NC}"
    go tool cover -func=coverage.out | grep -E "(handlers|utils|cache|services|total)" | while read line; do
        if echo "$line" | grep -q "total:"; then
            coverage=$(echo "$line" | awk '{print $3}' | sed 's/%//')
            if (( $(echo "$coverage >= 80" | bc -l) )); then
                echo -e "${GREEN}$line${NC}"
            elif (( $(echo "$coverage >= 70" | bc -l) )); then
                echo -e "${YELLOW}$line${NC}"
            else
                echo -e "${RED}$line${NC}"
            fi
        else
            echo "$line"
        fi
    done
    
    # Check coverage threshold
    total_coverage=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')
    echo ""
    if (( $(echo "$total_coverage >= 75" | bc -l) )); then
        echo -e "${GREEN}🎉 Coverage target met: ${total_coverage}%${NC}"
    else
        echo -e "${YELLOW}⚠️  Coverage below target: ${total_coverage}% (target: 75%)${NC}"
    fi
    
else
    echo -e "${RED}❌ No coverage file generated${NC}"
fi

# Run benchmarks for performance validation
echo -e "${BLUE}⚡ Running performance benchmarks...${NC}"
echo "JWT Benchmarks:"
go test -bench=BenchmarkGenerateToken ./utils -benchtime=1s
go test -bench=BenchmarkValidateToken ./utils -benchtime=1s

echo ""
echo "Handler Benchmarks:"
go test -bench=. ./api/handlers -benchtime=1s 2>/dev/null | grep Benchmark || echo "No handler benchmarks available"

echo ""
echo -e "${GREEN}🏁 Test suite completed!${NC}"
echo -e "${BLUE}📊 View detailed coverage: open coverage.html${NC}"
