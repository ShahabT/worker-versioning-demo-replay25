#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./decommission.sh <deployment version>
DEPLOYMENT_VERSION="$1"

if [ "$DEPLOYMENT_VERSION" = "__unversioned__" ]; then
  # Unversioned workers can be decommissioned immediately once a Current Deployment Version is set.
  docker compose -p orders down
  echo "➡️ Decommissioned the unversioned workers."
  exit 0
fi

# Wait until the version is drained and then kill the workers
while true; do
  DRAINAGE_STATUS=$(temporal worker deployment describe-version --version "$DEPLOYMENT_VERSION" -o json | jq -r ".drainageInfo.drainageStatus")
  if [ "$DRAINAGE_STATUS" = "drained" ]; then
    echo "➡️ Deployment Version $DEPLOYMENT_VERSION is DRAINED now. Decommissioning the workers..."
    docker compose -p "${DEPLOYMENT_VERSION//./_}" down
    exit 0
  fi
  echo "⏳  Deployment Version $DEPLOYMENT_VERSION is NOT DRAINED. Checking again in 5 seconds..."
  sleep 5
done
