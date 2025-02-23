#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./promote.sh <deployment version>
DEPLOYMENT_VERSION="$1"

PREVIOUS_DEPLOYMENT_VERSION=$(../../cli/temporal worker deployment describe --name oms-worker -o json | jq ".routingConfig.currentVersion" | tr -d '"')

../../cli/temporal worker deployment version set-current --version $DEPLOYMENT_VERSION -y

echo "Current Version changed from $PREVIOUS_DEPLOYMENT_VERSION to $DEPLOYMENT_VERSION."
echo "Run the following command to decommission $PREVIOUS_DEPLOYMENT_VERSION:"
echo "./decommission.sh $PREVIOUS_DEPLOYMENT_VERSION"