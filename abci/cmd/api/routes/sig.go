package routes

// Signature represents a signature in transaction APIs
type Signature struct {
	Type  string `json:"type"`
	Value string `json:"value" binding:"required,eth_sig"`
}
