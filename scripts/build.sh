#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./build.sh

# STEP 1: generate a new build ID
export BUILD_ID=v$(date '+%M')
echo "Generating new build with ID: $BUILD_ID  ..."

# STEP 2: create Docker image
IMAGE_TAG=orders-worker:$BUILD_ID
docker build --tag $IMAGE_TAG ../app > /dev/null

DEPLOYMENT_VERSION=orders.${BUILD_ID}
PREVIOUS_DEPLOYMENT_VERSION=$(temporal worker deployment describe --name orders -o json | jq ".routingConfig.currentVersion" | tr -d '"')
echo "➡️ Built the Docker Image. Next:"
echo "./deploy.sh $BUILD_ID"
