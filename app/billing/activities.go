package billing

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
)

// Item represents an item being ordered.
type Item struct {
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

// ChargeInput is the input for the Charge workflow.
type ChargeInput struct {
	CustomerID     string `json:"customerId"`
	Reference      string `json:"orderReference"`
	Items          []Item `json:"items"`
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
}

// ChargeResult is the result for the Charge workflow.
type ChargeResult struct {
	InvoiceReference string `json:"invoiceReference"`
	SubTotal         int32  `json:"subTotal"`
	Shipping         int32  `json:"shipping"`
	Tax              int32  `json:"tax"`
	Total            int32  `json:"total"`

	Success  bool   `json:"success"`
	AuthCode string `json:"authCode"`
}

// GenerateInvoiceInput is the input for the GenerateInvoice activity.
type GenerateInvoiceInput struct {
	CustomerID string `json:"customerId"`
	Reference  string `json:"orderReference"`
	Items      []Item `json:"items"`
}

type ShippingCost struct {
}

// GenerateInvoiceResult is the result for the GenerateInvoice activity.
type GenerateInvoiceResult struct {
	InvoiceReference string `json:"invoiceReference"`
	SubTotal         int32  `json:"subTotal"`
	Shipping         int32  `json:"shipping"`
	Tax              int32  `json:"tax"`
	Total            int32  `json:"total"`
	Credit           int32  `json:"credit"`
}

// ChargeCustomerInput is the input for the ChargeCustomer activity.
type ChargeCustomerInput struct {
	CustomerID string `json:"customerId"`
	Reference  string `json:"reference"`
	Charge     int32  `json:"charge"`
}

// ChargeCustomerResult is the result for the GenerateInvoice activity.
type ChargeCustomerResult struct {
	Success  bool   `json:"success"`
	AuthCode string `json:"authCode"`
}

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

// SpendCredits adjust the invoice base on the credits this customer has so credits get spent first.
func (a *Activities) SpendCredits(
	ctx context.Context,
	invoice *GenerateInvoiceResult,
) (*GenerateInvoiceResult, error) {
	// Assuming user always has a credit of 1200.
	creditToSpend := min(invoice.Total, 1200)
	invoice.Credit = creditToSpend
	invoice.Total -= creditToSpend
	return invoice, nil
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

	// Simulate some short delay
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)

	activity.GetLogger(ctx).Info(
		"Charge",
		"Customer", input.CustomerID,
		"Amount", input.Charge,
		"Reference", input.Reference,
		"Success", result.Success,
	)

	return &result, nil
}
