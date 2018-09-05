#!/bin/bash

send-cpu(){
    local host=$1
    local pct=$2
    docker-compose exec elasticsearch curl -XPOST -H 'Content-Type: application/json' localhost:9200/test/_doc -d "$(cat << EOF
{
    "@timestamp": "$(TZ=0 date '+%Y-%m-%dT%H:%M:%SZ')",
    "host": "$host",
    "system": {
        "cpu": {
            "total": {
                "pct": $pct
            }
        }
    }
}
EOF
)"
}

send-cpu "web-dev" "0.9"
send-cpu "web-prod" "0.9"
# send-cpu "web-dev" "0.7"
