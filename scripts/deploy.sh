#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./deploy.sh <build ID>
export BUILD_ID="$1"

# Upgrade previous deployment in-place.
docker-compose -f ../docker-compose.yaml -p orders up -d

echo "➡️ Upgraded workers."
