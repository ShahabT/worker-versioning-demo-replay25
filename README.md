# Temporal Replay 2025 Demo: Getting Started with Worker Versioning

## Prepare
1. Install and run Docker. We'll use `docker compose`.
2. Run dev server in a separate terminal using Temporal CLI:
   ```shell
   temporal server start-dev \
   --dynamic-config-value frontend.workerVersioningWorkflowAPIs=true \
   --dynamic-config-value system.enableDeploymentVersions=true \
   --dynamic-config-value 'matching.wv.VersionDrainageStatusVisibilityGracePeriod="30s"' \
   --dynamic-config-value 'matching.wv.VersionDrainageStatusRefreshInterval="5s"'
   ```
3. Run load-gen in a separate terminal: `cd app; go run load-gen/load-gen.go`


# Demo Script

## Brownfield Migration

### Baseline
1. Switch to `baseline` branch.
2. Run `./build.sh` to build the code.
3. Run `./deploy.sh <build id>` to deploy the code in place (no Rainbow deployments).
4. Go to UI and ensure workflows are running.
5. Note the workflows: one short running (Charge), one long-lived (Shipment).

### Add Rainbow Deployment Support
1. Switch to `use-versioning` branch.
2. Note the changes to `deploy.sh` that uses Rainbow deployment strategy now.
3. Note the new steps: `promote.sh` and `decommission.sh`.

### Configure Workers
1. Now that the deployment system is ready we can enable Versioning in the workers.
2. Note that now we pass Deployment Version to the deployed worker via and env var in `docker-compose.yaml`.
3. Note the changes made in `main.go` to pass DeploymentOptions and Versioning Behaviors.

### First Versioned Deploy
1. Unversioned workers are running.
2. Only worker config is changed to enable Versioning but no other code changes.
3. Run `build.sh` followed by `deploy.sh`.
4. Now both unversioned and versioned workers should be running. Can see using `docker-compose ls`.
5. In the UI, we now see a Worker Deployment. And a Worker Deployment Version in `Inactive` state.
6. The Current Version is still `__unversioned__`.
7. Run `promote.sh` to set the Current Version.
8. Now workflows should start getting values in their Versioning Behavior and Deployment Version columns.
9. All workflow should continue smoothly.
10. You can decommission the unversioned workers immediately. No need to wait for any "drainage".

### Incompatible Change to Pinned Workflows
1. Switch to `spend-credit` branch.
2. Notice the new activity added in the Charge workflow without Patching.
3. Run `cicd.sh` to run the full build and deploy process.
4. Notice that the promote step now has a ramp step that checks the following at a 10% ramp:
   - At least one wf succeeded in the new build.
   - No errors are logged by the workers.
5. Once the new Version is promoted, notice the old Version going through `Draining` and `Drained` status.
6. Once `Drained`, the decommission script should kill the workers.

## Sad Path

1. Switch to `bad-build` branch.
2. Run `cicd.sh` to do the f.
3. Notice the ramp verification fails for errors. And the ramp is cancelled.
4. There are some pinned workflows on the new Version.
5. We don't want to reset them because we don't want them to lose their state.
6. How can we roll forward?
7. Switch to `fixed-build` branch.
8. Build and deploy, but not promote the patched version.
9. Run `move-pinned.sh` to move the stuck workflows to the fixed Version.
10. Verify that the stuck workflow is completed now and the bad Version is Drained.
11. Now that the patch is confirmed on the broken workflows, run `promote.sh` to use it everywhere. 
12. Run `decommission.sh` to decommission both the bad and old Versions.