package shipment

import (
	"math/rand"
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

// Item represents an item being shipped.
type Item struct {
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

// ShipmentInput is the input for a Shipment workflow.
type ShipmentInput struct {
	ID    string
	Items []Item
}

// StatusQuery is the name of the query to use to fetch a Shipment's status.
const StatusQuery = "status"

// ShipmentCarrierUpdateSignalName is the name for a signal to update a shipment's status from the carrier.
const ShipmentCarrierUpdateSignalName = "ShipmentCarrierUpdate"

// ShipmentStatusUpdatedSignalName is the name for a signal to notify of an update to a shipment's status.
const ShipmentStatusUpdatedSignalName = "ShipmentStatusUpdated"

const (
	// ShipmentStatusPending represents a shipment that has not yet been booked with a carrier
	ShipmentStatusPending = "pending"
	// ShipmentStatusBooked represents a shipment acknowledged by a carrier, but not yet picked up
	ShipmentStatusBooked = "booked"
	// ShipmentStatusDispatched represents a shipment picked up by a carrier, but not yet delivered to the customer
	ShipmentStatusDispatched = "dispatched"
	// ShipmentStatusDelivered represents a shipment that has been delivered to the customer
	ShipmentStatusDelivered = "delivered"
	ShipmentStatusExpired   = "expired"
)

// ShipmentCarrierUpdateSignal is used by a carrier to update a shipment's status.
type ShipmentCarrierUpdateSignal struct {
	Status string `json:"status"`
}

// ShipmentStatusUpdatedSignal is used to notify the requestor of an update to a shipment's status.
type ShipmentStatusUpdatedSignal struct {
	ShipmentID string    `json:"shipmentID"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ShipmentResult is the result of a Shipment workflow.
type ShipmentResult struct {
	CourierReference string
}

type shipmentImpl struct {
	id        string
	status    string
	updatedAt time.Time

	logger log.Logger
}

// Shipment implements the Shipment workflow.
func Shipment(ctx workflow.Context, input *ShipmentInput) (*ShipmentResult, error) {
	wf := new(shipmentImpl)

	if err := wf.setup(ctx, input); err != nil {
		return nil, err
	}

	return wf.run(ctx, input)
}

func (s *shipmentImpl) setup(ctx workflow.Context, input *ShipmentInput) error {
	s.id = input.ID
	s.status = ShipmentStatusPending

	s.logger = log.With(
		workflow.GetLogger(ctx),
		"shipmentId", s.id,
	)

	return workflow.SetQueryHandler(ctx, StatusQuery, func() (*ShipmentStatus, error) {
		return &ShipmentStatus{
			ID:        s.id,
			Status:    s.status,
			UpdatedAt: s.updatedAt,
			Items:     input.Items,
		}, nil
	})
}

func (s *shipmentImpl) run(ctx workflow.Context, input *ShipmentInput) (*ShipmentResult, error) {
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Second,
		},
	)

	var result BookShipmentResult

	err := workflow.ExecuteActivity(ctx,
		a.BookShipment,
		BookShipmentInput{
			Reference: s.id,
			Items:     input.Items,
		},
	).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	s.updateStatus(ctx, ShipmentStatusBooked)

	err = s.handleCarrierUpdates(ctx)

	return &ShipmentResult{
		CourierReference: result.CourierReference,
	}, err
}

func (s *shipmentImpl) handleCarrierUpdates(ctx workflow.Context) error {
	timer := workflow.NewTimer(ctx, time.Duration(10+rand.Intn(10))*time.Minute)
	signalCh := workflow.GetSignalChannel(ctx, ShipmentCarrierUpdateSignalName)

	sel := workflow.NewSelector(ctx)
	sel.AddReceive(signalCh, func(ch workflow.ReceiveChannel, more bool) {
		var signal ShipmentCarrierUpdateSignal
		ch.Receive(ctx, &signal)
		s.logger.Info("Received carrier update", "status", signal.Status)
		s.updateStatus(ctx, signal.Status)
	})
	sel.AddFuture(timer, func(_ workflow.Future) {
		s.updateStatus(ctx, ShipmentStatusExpired)
		// Shipment expired. Nothing to do, just allow workflow to close.
	})

	for s.status != ShipmentStatusDelivered && s.status != ShipmentStatusExpired {
		sel.Select(ctx)
	}

	return nil
}

func (s *shipmentImpl) updateStatus(ctx workflow.Context, status string) {
	s.status = status
	s.updatedAt = workflow.Now(ctx)
}
