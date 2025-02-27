#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./promote.sh <deployment version>
DEPLOYMENT_VERSION="$1"
PREVIOUS_DEPLOYMENT_VERSION=$(temporal worker deployment describe --name orders -o json | jq ".routingConfig.currentVersion" | tr -d '"')
sleep 3 # give enough time for pollers to arrive to the server

# Step 1: Ramp 10%
echo "➡️ Setting ramp to 10% for deployment version ${DEPLOYMENT_VERSION}"
temporal worker deployment set-ramping-version --version "$DEPLOYMENT_VERSION" --percentage 10 -y > /dev/null

echo "⏳  Waiting for 60 seconds before verification..."
sleep 60

# Step 2: Verification
echo "➡️ Verifying workflows for deployment version ${DEPLOYMENT_VERSION}..."
if ! ./ramp-verification.sh $DEPLOYMENT_VERSION; then
  echo "❌  Verification failed, cancelling the ramp."
  temporal worker deployment set-ramping-version --version "$DEPLOYMENT_VERSION" --percentage 0 --delete -y > /dev/null
  exit 1
fi

# Step 3: Set Current
echo "✅  Verification successful, setting $DEPLOYMENT_VERSION as Current."
temporal worker deployment set-current-version --version $DEPLOYMENT_VERSION -y > /dev/null

echo "➡️ Current Version changed from $PREVIOUS_DEPLOYMENT_VERSION to $DEPLOYMENT_VERSION. Next:"
echo "./decommission.sh $PREVIOUS_DEPLOYMENT_VERSION"
