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

func main() {
	c, err := client.Dial(
		client.Options{},
	)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	timer := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timer:
			fmt.Println("Starting new Charge workflow")
			_, err := c.ExecuteWorkflow(
				context.Background(),
				client.StartWorkflowOptions{
					TaskQueue: "orders",
				},
				billing.Charge,
				generateChargeInput(),
			)
			if err != nil {
				fmt.Errorf("failed to start a Charge %e\n", err)
			}
			fmt.Println("Starting new Shipment workflow")
			s := generateShipmentInput()
			_, err = c.ExecuteWorkflow(
				context.Background(),
				client.StartWorkflowOptions{
					ID:        s.ID,
					TaskQueue: "orders",
				},
				shipment.Shipment,
				s,
			)
			if err != nil {
				fmt.Errorf("failed to start a Shipment %e\n", err)
			}
		}
	}
}

func generateShipmentInput() *shipment.ShipmentInput {
	return &shipment.ShipmentInput{
		ID: uuid.NewString(),
		Items: []shipment.Item{
			{
				SKU:      strconv.Itoa(2000 + rand.Intn(10)),
				Quantity: int32(rand.Intn(20)),
			},
		},
	}
}

//func getShipmentID(secondsBefore int) string {
//	ts := time.Now().Unix() - int64(secondsBefore)
//	ts -= ts % intervalSec // ts is always a multiplier of intervalSec
//	return fmt.Sprintf("shipment-%d", ts)
//}

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
