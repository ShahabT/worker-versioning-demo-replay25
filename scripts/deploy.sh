#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./deploy.sh <build ID>
export BUILD_ID="$1"

# Temporal expects Deployment Versions to be of the following format: <Deployment Name>.<Build ID>
export DEPLOYMENT_VERSION=orders.${BUILD_ID}

# Rainbow deployment. Each Deployment Version is separate from others.
docker compose -f ../docker-compose.yaml -p ${DEPLOYMENT_VERSION//./_} up -d

echo "➡️ Deployed $DEPLOYMENT_VERSION. Next:"
echo "./promote.sh $DEPLOYMENT_VERSION"