#!/bin/bash

BASE_URL="http://localhost:8088/api"
REQUESTS=10

echo "🚀 GoFiver API Benchmark"
echo "========================"
echo ""

benchmark() {
    local name=$1
    local url=$2
    local total=0
    local min=9999
    local max=0
    
    echo "⏱️  Testing: $name"
    
    for i in $(seq 1 $REQUESTS); do
        time=$(curl -s -o /dev/null -w "%{time_total}" "$url")
        time_ms=$(echo "$time * 1000" | bc)
        total=$(echo "$total + $time_ms" | bc)
        
        if (( $(echo "$time_ms < $min" | bc -l) )); then
            min=$time_ms
        fi
        if (( $(echo "$time_ms > $max" | bc -l) )); then
            max=$time_ms
        fi
    done
    
    avg=$(echo "scale=0; $total / $REQUESTS" | bc)
    min=$(printf "%.0f" $min)
    max=$(printf "%.0f" $max)
    
    printf "   Avg: %sms | Min: %sms | Max: %sms\n\n" "$avg" "$min" "$max"
}

benchmark "GET /blogs (paginated)" "$BASE_URL/blogs?page=1&limit=20"
benchmark "GET /blogs/:id (single)" "$BASE_URL/blogs/1"
benchmark "GET /users" "$BASE_URL/users?page=1&per_page=20"
benchmark "GET /users/:id" "$BASE_URL/users/1"

echo "✅ Benchmark complete!"
