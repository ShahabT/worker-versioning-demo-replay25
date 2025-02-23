package billing

import (
	"context"
	"fmt"
	"math/rand"

	"go.temporal.io/sdk/activity"
)

// Activities implements the billing package's Activities.
// Any state shared by the worker among the activities is stored here.
type Activities struct {
}

var a Activities

// GenerateInvoice activity creates an invoice for a fulfillment.
func (a *Activities) GenerateInvoice(
	ctx context.Context,
	input *GenerateInvoiceInput,
) (*GenerateInvoiceResult, error) {
	var result GenerateInvoiceResult

	if input.CustomerID == "" {
		return nil, fmt.Errorf("CustomerID is required")
	}
	if input.Reference == "" {
		return nil, fmt.Errorf("OrderReference is required")
	}
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("invoice must have items")
	}

	result.InvoiceReference = input.Reference

	for _, item := range input.Items {
		cost, tax := calculateCosts(item)
		result.SubTotal += cost
		result.Tax += tax
		result.Shipping += calculateShippingCost(item)
		result.Total += result.SubTotal + result.Tax + result.Shipping
	}

	activity.GetLogger(ctx).Info(
		"Invoice",
		"Customer", input.CustomerID,
		"Total", result.Total,
		"Reference", result.InvoiceReference,
	)

	return &result, nil
}

// calculateCosts calculates the cost and tax for an item.
func calculateCosts(item Item) (cost int32, tax int32) {
	// This is just a simulation, so make up a cost
	// Normally this would be looked up on the SKU
	costPerUnit := 3500 + rand.Int31n(8500)
	// Return tax at 20%
	return costPerUnit * int32(item.Quantity), costPerUnit * int32(item.Quantity) / 5
}

// calculateShippingCost calculates the shipping cost for an item.
func calculateShippingCost(item Item) int32 {
	// This is just a simulation, so make up a cost
	// Normally this would be looked up on the SKU
	costPerUnit := 500 + rand.Int31n(500)
	return costPerUnit * int32(item.Quantity)
}

// ChargeCustomer activity charges a customer for a fulfillment.
func (a *Activities) ChargeCustomer(
	ctx context.Context,
	input *ChargeCustomerInput,
) (*ChargeCustomerResult, error) {
	var result ChargeCustomerResult

	result.Success = true
	result.AuthCode = "1234"

	activity.GetLogger(ctx).Info(
		"Charge",
		"Customer", input.CustomerID,
		"Amount", input.Charge,
		"Reference", input.Reference,
		"Success", result.Success,
	)

	return &result, nil
}
