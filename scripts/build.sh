#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./build.sh

# STEP 1: generate a new build ID
#export BUILD_ID=$(date '+%Y%m%d%H%M%S')$(git rev-parse --short HEAD)
export BUILD_ID=v$(date '+%M')
echo "New Build ID: $BUILD_ID"

# STEP 2: create Docker image
IMAGE_TAG=oms-worker:$BUILD_ID
docker build --tag $IMAGE_TAG ../app

echo "Built Image, run the following command to deploy:"
echo "./deploy.sh $BUILD_ID"
