package leads

// Scable is an interface to handle lead logic at smarcenter phase
type Scable interface {
	Active(Lead) bool
	// Send(Lead) (ScResponse, error)
	Send(Lead) ScResponse
}

// ScResponse represents an Smart Center response
type ScResponse struct {
	Success    bool  `json:"success"`
	StatusCode int64 `json:"status"`
	ID         int64 `json:"id"`
	Error      error `json:"error,omitempty"`
}
