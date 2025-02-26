package main

import (
	"log"
	"os"

	"github.com/ShahabT/worker-versioning-demo-replay25/billing"
	"github.com/ShahabT/worker-versioning-demo-replay25/shipment"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	c := createTemporalClient()
	defer c.Close()

	w := worker.New(c, "orders", worker.Options{
		DeploymentOptions: worker.DeploymentOptions{
			UseVersioning:             true,
			Version:                   os.Getenv("DEPLOYMENT_VERSION"),
			DefaultVersioningBehavior: workflow.VersioningBehaviorAutoUpgrade,
		},
	})
	w.RegisterWorkflowWithOptions(billing.Charge, workflow.RegisterOptions{
		VersioningBehavior: workflow.VersioningBehaviorPinned,
	})
	w.RegisterWorkflow(shipment.Shipment)
	w.RegisterActivity(&billing.Activities{})
	w.RegisterActivity(&shipment.Activities{})

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Error while running workers", err)
		return
	}
}

func createTemporalClient() client.Client {
	c, err := client.Dial(client.Options{
		HostPort: "host.docker.internal:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return c
}
