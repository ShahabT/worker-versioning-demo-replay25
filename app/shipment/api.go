package shipment

import (
	"strings"
	"time"
)

// TaskQueue is the default task queue for the Shipment system.
const TaskQueue = "shipments"

// StatusQuery is the name of the query to use to fetch a Shipment's status.
const StatusQuery = "status"

// ShipmentWorkflowID returns the workflow ID for a Shipment.
func ShipmentWorkflowID(id string) string {
	return "Shipment:" + id
}

// ShipmentIDFromWorkflowID returns the ID for a Shipment from a WorkflowID.
func ShipmentIDFromWorkflowID(id string) string {
	return strings.TrimPrefix(id, "Shipment:")
}

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
