#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Usage: ./move-pinned.sh <from version> <to version>
FROM_VERSION="$1"
TO_VERSION="$2"

# Step 1: Ramp 10%
#echo "➡️ Setting ramp to 10% for deployment version ${DEPLOYMENT_VERSION}"
temporal workflow update-options \
    --query="TemporalWorkerDeploymentVersion='$FROM_VERSION' AND TemporalWorkflowVersioningBehavior='Pinned'" \
    --versioning-override-behavior=pinned \
    --versioning-override-pinned-version="$TO_VERSION" -y -o json > ~/batch_id

echo Batch Job: http://localhost:8233/namespaces/default/batch-operations/$(cat ~/batch_id | jq ".batchJobId" | tr -d '"')
