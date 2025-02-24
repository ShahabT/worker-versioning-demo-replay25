# Temporal Replay 2025 Demo: Getting Started with Worker Versioning

## Prepare
1. Install and run Docker Desktop.
2. Clone CLI repo: https://github.com/temporalio/cli. Make sure the cli root directory and this project's root directory are adjacent.
3. Checkout this branch locally: `shahab/versioning-3.1`
4. Build CLI (run from cli root directory):
    ```shell
    go build ./cmd/temporal
    ```
5. Run dev server using the binary in the cli home folder:
   ./temporal server start-dev \
   --dynamic-config-value frontend.workerVersioningWorkflowAPIs=true \
   --dynamic-config-value system.enableDeploymentVersions=true \
   --dynamic-config-value 'matching.wv.VersionDrainageStatusVisibilityGracePeriod="30s"' \
   --dynamic-config-value 'matching.wv.VersionDrainageStatusRefreshInterval="5s"'

## Run Scripts
1. Go to `scripts/` dir in this project. From now on, all commands assume you are in this folder.
2. Run `./build.sh` to build a worker version.
3. The build command prints the next command to run: deploy. Run that.
4. Now you should have one version polling. The following command should show a Deployment with one version:
    ```shell
    ../../cli/temporal worker deployment describe -d oms-worker
    ```
5. The deploy script in previous step prints the next command to run: promote. Run that.
6. Now, the describe command above should show the current version is not `__unversioned__` anymore.
7. Run build, deploy, and promote commands again to create and promote more versions.
8. Everytime the promote script is called, the previous version becomes `draining` and eventually `drained`
9. For ramp, there is no script here yet, but you can run the following command between the deploy and promote steps to
ramp the version:
   ```shell
   ../../cli/temporal worker deployment version set-ramping --version $DEPLOYMENT_VERSION --percentage 10.0
   ```

## Generate Load
Go to the `app` folder and run the following command so the load-gen script starts generating new wfs every few 
seconds.
   ```shell
   go run load-gen/load-gen.go
   ```

This script generates two workflow types: Charge (short-running, PINNED) and Shipment (long-running, AutoUpgrade). Note 
that for workflows to get a Pinned/AutoUpgrade versioning behavior, there needs to be a deployment version running and 
set as current first. Load-gen only starts the workflows, the versioning behavior is set only after the first task of 
the workflow is processed.
