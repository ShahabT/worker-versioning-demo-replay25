package acceptance_test

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/ShahabT/worker-versioning-demo-replay25/billing"
	"github.com/ShahabT/worker-versioning-demo-replay25/shipment"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

func AcceptanceTest(ctx workflow.Context) error {
	// Test Charge workflow
	chargeInput := generateTestChargeInput()
	var chargeResult *billing.ChargeResult
	workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{})
	err := workflow.ExecuteChildWorkflow(ctx, billing.Charge, chargeInput).Get(ctx, &chargeResult)
	if err != nil {
		return err
	}
	if err := validateChargeResult(chargeResult); err != nil {
		return err
	}

	// Test Shipment workflow
	shipmentInput := generateTestShipmentInput()
	var shipmentExecution workflow.Execution
	fut := workflow.ExecuteChildWorkflow(ctx, shipment.Shipment, shipmentInput)
	// Signal the workflow and say it's delivered
	err = fut.SignalChildWorkflow(ctx,
		shipment.ShipmentCarrierUpdateSignalName,
		&shipment.ShipmentCarrierUpdateSignal{
			Status: shipment.ShipmentStatusDelivered,
		}).Get(ctx, &shipmentExecution)
	if err != nil {
		return err
	}
	var shipmentResult *shipment.ShipmentResult
	err = fut.Get(ctx, &shipmentResult)
	if err != nil {
		return err
	}
	return nil
}

func validateChargeResult(result *billing.ChargeResult) error {
	if !result.Success {
		return errors.New("charge was not successful")
	}
	return nil
}

func generateTestChargeInput() *billing.ChargeInput {
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

func generateTestShipmentInput() *shipment.ShipmentInput {
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
