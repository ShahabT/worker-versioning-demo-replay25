package main

import (
	"context"
	"log"
	"os"

	"github.com/ShahabT/worker-versioning-demo-replay25/billing"
	"github.com/ShahabT/worker-versioning-demo-replay25/shipment"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/sync/errgroup"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c := createTemporalClient()
	defer c.Close()

	billingWorker := worker.New(c, billing.TaskQueue, worker.Options{})
	billingWorker.RegisterWorkflow(billing.Charge)
	billingWorker.RegisterActivity(&billing.Activities{})

	shipmentWorker := worker.New(c, shipment.TaskQueue, worker.Options{})
	shipmentWorker.RegisterWorkflow(shipment.Shipment)
	shipmentWorker.RegisterActivity(&shipment.Activities{})

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return billingWorker.Run(worker.InterruptCh())
	})
	g.Go(func() error {
		return shipmentWorker.Run(worker.InterruptCh())
	})
	if err := g.Wait(); err != nil {
		log.Fatalln("Error while running workers", err)
		return
	}
}

func getDeploymentOptions() worker.DeploymentOptions {
	return worker.DeploymentOptions{
		UseVersioning:             true,
		Version:                   os.Getenv("DEPLOYMENT_VERSION"),
		DefaultVersioningBehavior: workflow.VersioningBehaviorAutoUpgrade,
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
