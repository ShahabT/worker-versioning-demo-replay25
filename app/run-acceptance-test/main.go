package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	acceptance_test "github.com/ShahabT/worker-versioning-demo-replay25/acceptance-test"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func main() {
	c := createTemporalClient()
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	newDeploymentVersion := os.Args[1]
	fmt.Printf("➡️ Starting AcceptanceTest workflow on %s\n", newDeploymentVersion)
	run, err := c.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			VersioningOverride: client.VersioningOverride{
				Behavior:      workflow.VersioningBehaviorPinned,
				PinnedVersion: newDeploymentVersion,
			},
			WorkflowExecutionTimeout: 10 * time.Second,
			TaskQueue:                "orders",
		},
		acceptance_test.AcceptanceTest,
	)
	if err != nil {
		fmt.Printf("❌  Failed to start the AcceptanceTest workflow \n%e\n", err)
		os.Exit(1)
	}

	//var res *any
	err = c.GetWorkflow(ctx, run.GetID(), run.GetRunID()).Get(ctx, nil)
	if err != nil {
		fmt.Printf("❌  Acceptance Test failed \n%e\n", err)
		fmt.Printf("http://localhost:8233/namespaces/default/workflows/%s\n", run.GetID())
		os.Exit(1)
	}
	fmt.Printf("✅  Acceptance Test passed!\n")
}

func createTemporalClient() client.Client {
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return c
}
