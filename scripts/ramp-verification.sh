#!/bin/bash
set -uo pipefail
IFS=$'\n\t'

# Usage: ./ramp-verification.sh <deployment version>
DEPLOYMENT_VERSION="$1"

# Checks:
# - At least one wf completed on DEPLOYMENT_VERSION.
# - No error logs are observed in the workers of DEPLOYMENT_VERSION.

# Count successful workflows using Temporal CLI
QUERY="TemporalWorkerDeploymentVersion = '$DEPLOYMENT_VERSION' AND ExecutionStatus = 'Completed'"
success_count=$(temporal workflow count --query "$QUERY" -o json | jq ".count // 0" | tr -d '"')

if [ $success_count -eq 0 ]; then
  echo "     ❌  Verification failed: Found no completed workflow(s)."
  echo "http://localhost:8233/namespaces/default/workflows?query=$(jq -Rr @uri <<< "$QUERY")"
  exit 1
else
  echo "     ✅  $success_count workflow(s) completed on $DEPLOYMENT_VERSION."
fi

# Check logs for errors.
error_logs=$(docker compose -p ${DEPLOYMENT_VERSION//./_} logs -n 1000 | grep ERROR | grep WorkflowID)
count_errors=$(echo $error_logs | grep ERROR | grep WorkflowID | wc -l | xargs)
if [ $count_errors -eq 0 ]; then
  echo "     ✅  No error logged."
else
  echo "     ❌  Verification failed: Found unexpected error logs. Example workflow:"
  example_log=$(echo $error_logs | grep ERROR | grep WorkflowID | head -n 1)
  example_wf_id=$(echo $example_log | sed -n 's/.*WorkflowID \([^ ]*\).*/\1/p')
  echo "http://localhost:8233/namespaces/default/workflows/$example_wf_id"
  exit 1
fi
