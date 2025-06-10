#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api/v1"
FAILED_TESTS=0

# Function to test endpoint
test_endpoint() {
    local endpoint=$1
    local expected_status=${2:-200}
    local description=$3
    
    echo -e "${YELLOW}Testing: $description${NC}"
    
    response=$(curl -s -w "%{http_code}" -o /tmp/response.json "$BASE_URL$endpoint")
    status_code="${response: -3}"
    
    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ PASS${NC} - Status: $status_code"
        if [ -s /tmp/response.json ]; then
            echo "Response preview: $(head -c 100 /tmp/response.json)..."
        fi
    else
        echo -e "${RED}✗ FAIL${NC} - Expected: $expected_status, Got: $status_code"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo "---"
}

# Run tests
echo "Starting API Tests..."
echo "====================="

test_endpoint "/rankings" 200 "Get rankings endpoint"
test_endpoint "/teams" 200 "Get teams endpoint" 
test_endpoint "/health" 200 "Health check endpoint"

# Test a non-existent endpoint
test_endpoint "/nonexistent" 404 "Non-existent endpoint (should return 404)"

# Summary
echo "====================="
if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}$FAILED_TESTS test(s) failed${NC}"
    exit 1
fi
