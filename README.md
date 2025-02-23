# Temporal Replay 2025 Demo: Getting Started with Worker Versioning

How to run:
1. Enable K8s on Docker Desktop
2. run `./cicd.sh` to build an image and deploy the workers.
3. run `cd app; go run load-gen/load-gen.go` to start generating load (ExecuteWorkflow requests)
4. Make any changes you want to the WF or activities, don't mind about compatibility!
5. run step again to deploy the new code and stop the old workers once all old WFs are closed.