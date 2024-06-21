package model

type ReplicasRequest struct {
	Replicas int `json:"replicas"`
}

type CPU struct {
	HighPriority float64 `json:"highPriority"`
}
type StatusResponse struct {
	CPU      CPU `json:"cpu"`
	Replicas int `json:"replicas"`
}
