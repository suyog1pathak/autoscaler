package model

// ReplicasRequest represents the JSON structure for updating replica count.
type ReplicasRequest struct {
	Replicas int `json:"replicas"`
}

// CPU represents the CPU metrics structure.
type CPU struct {
	HighPriority float64 `json:"highPriority"`
}

// StatusResponse represents the JSON structure for status response.
type StatusResponse struct {
	CPU      CPU `json:"cpu"`
	Replicas int `json:"replicas"`
}
