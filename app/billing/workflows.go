package billing

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Charge Workflow invoices and processes payment for a fulfillment.
func Charge(ctx workflow.Context, input *ChargeInput) (*ChargeResult, error) {
	logger := workflow.GetLogger(ctx)
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			ScheduleToCloseTimeout: 30 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumInterval: time.Second,
			},
		},
	)

	var invoice GenerateInvoiceResult
	err := workflow.ExecuteActivity(ctx,
		a.GenerateInvoice,
		GenerateInvoiceInput{
			CustomerID: input.CustomerID,
			Reference:  input.Reference,
			Items:      input.Items,
		},
	).Get(ctx, &invoice)
	if err != nil {
		return nil, err
	}

	var charge ChargeCustomerResult
	err = workflow.ExecuteActivity(ctx,
		a.ChargeCustomer,
		ChargeCustomerInput{
			CustomerID: input.CustomerID,
			Reference:  invoice.InvoiceReference,
			Charge:     invoice.Total,
		},
	).Get(ctx, &charge)
	if err != nil {
		logger.Warn("Charge failed", "customer_id", input.CustomerID, "error", err)
		charge.Success = false
	}

	return &ChargeResult{
		InvoiceReference: invoice.InvoiceReference,
		SubTotal:         invoice.SubTotal,
		Tax:              invoice.Tax,
		Shipping:         invoice.Shipping,
		Total:            invoice.Total,

		Success:  charge.Success,
		AuthCode: charge.AuthCode,
	}, nil
}
