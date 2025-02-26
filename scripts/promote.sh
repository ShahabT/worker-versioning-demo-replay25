#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./promote.sh <deployment version>
DEPLOYMENT_VERSION="$1"
PREVIOUS_DEPLOYMENT_VERSION=$(temporal worker deployment describe --name orders -o json | jq ".routingConfig.currentVersion" | tr -d '"')

temporal worker deployment set-current-version --version $DEPLOYMENT_VERSION -y > /dev/null

echo "➡️ Current Version changed from $PREVIOUS_DEPLOYMENT_VERSION to $DEPLOYMENT_VERSION. Next:"
echo "./decommission.sh $PREVIOUS_DEPLOYMENT_VERSION"
