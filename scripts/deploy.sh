#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./deploy.sh <build ID>
export BUILD_ID="$1"

# Temporal expects Deployment Versions to be of the following format: <Deployment Name>.<Build ID>
export DEPLOYMENT_VERSION=oms-worker.${BUILD_ID}

docker-compose -f ../docker-compose.yaml -p ${DEPLOYMENT_VERSION//./_} up -d

echo "Deployed Version $DEPLOYMENT_VERSION. Run the following command to promote:"
echo "./promote.sh $DEPLOYMENT_VERSION"