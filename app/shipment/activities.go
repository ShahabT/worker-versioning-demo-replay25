package shipment

import (
	"context"
	"time"
)

// Activities implements the shipment package's Activities.
// Any state shared by the worker among the activities is stored here.
type Activities struct {
}

var a Activities

// ShipmentStatus holds the status of a Shipment.
type ShipmentStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Items     []Item    `json:"items"`
}

// ShipmentStatusUpdate is used to update the status of a Shipment.
type ShipmentStatusUpdate struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ListShipmentEntry is an entry in the Shipment list.
type ListShipmentEntry struct {
	ID     string `json:"id" db:"id" bson:"id"`
	Status string `json:"status" db:"status" bson:"status"`
}

// BookShipmentInput is the input for the BookShipment operation.
// All fields are required.
type BookShipmentInput struct {
	Reference string
	Items     []Item
}

// BookShipmentResult is the result for the BookShipment operation.
// CourierReference is recorded where available, to allow tracking enquiries.
type BookShipmentResult struct {
	CourierReference string
}

// BookShipment engages a courier who can deliver the shipment to the customer
func (a *Activities) BookShipment(
	_ context.Context,
	input *BookShipmentInput,
) (*BookShipmentResult, error) {
	return &BookShipmentResult{
		CourierReference: input.Reference + ":1234",
	}, nil
}
