package main

import (
	"log"

	"github.com/ShahabT/worker-versioning-demo-replay25/billing"
	"github.com/ShahabT/worker-versioning-demo-replay25/shipment"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c := createTemporalClient()
	defer c.Close()

	w := worker.New(c, "orders", worker.Options{})
	w.RegisterWorkflow(billing.Charge)
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
