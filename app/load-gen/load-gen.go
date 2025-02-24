package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/ShahabT/worker-versioning-demo-replay25/billing"
	"github.com/ShahabT/worker-versioning-demo-replay25/shipment"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

const intervalSec = 2

func main() {
	c, err := client.Dial(
		client.Options{},
	)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	timer := time.Tick(intervalSec * time.Second)

	for {
		select {
		case <-timer:
			fmt.Println("Starting new Charge workflow")
			_, err := c.ExecuteWorkflow(
				context.Background(),
				client.StartWorkflowOptions{
					TaskQueue: billing.TaskQueue,
				},
				billing.Charge,
				generateChargeInput(),
			)
			if err != nil {
				panic(err)
			}
			fmt.Println("Starting new Shipment workflow")
			s := generateShipmentInput()
			_, err = c.ExecuteWorkflow(
				context.Background(),
				client.StartWorkflowOptions{
					ID:        s.ID,
					TaskQueue: shipment.TaskQueue,
				},
				shipment.Shipment,
				s,
			)
			if err != nil {
				panic(err)
			}
			err = c.SignalWorkflow(
				context.Background(),
				getShipmentID(intervalSec*30),
				"",
				shipment.ShipmentCarrierUpdateSignalName,
				&shipment.ShipmentCarrierUpdateSignal{
					Status: shipment.ShipmentStatusDelivered,
				},
			)
			if err == nil {
				fmt.Println("Shipment delivered")
			} else {
				fmt.Println(err)
			}
		}
	}
}

func generateShipmentInput() *shipment.ShipmentInput {
	return &shipment.ShipmentInput{
		ID: getShipmentID(0),
		Items: []shipment.Item{
			{
				SKU:      strconv.Itoa(2000 + rand.Intn(10)),
				Quantity: int32(rand.Intn(20)),
			},
		},
	}
}

func getShipmentID(secondsBefore int) string {
	ts := time.Now().Unix() - int64(secondsBefore)
	ts -= ts % intervalSec // ts is always a multiplier of intervalSec
	return fmt.Sprintf("shipment-%d", ts)
}

func generateChargeInput() *billing.ChargeInput {
	return &billing.ChargeInput{
		CustomerID: strconv.Itoa(rand.Intn(100)),
		Reference:  uuid.New().String(),
		Items: []billing.Item{
			{
				SKU:      strconv.Itoa(2000 + rand.Intn(10)),
				Quantity: int32(rand.Intn(20)),
			},
		},
	}
}
