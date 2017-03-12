package chord

// NotifyRequest represents a request sent to successor to notify it about local node
type NotifyRequest struct {
	ID   string
	host string
}

// NotifyResponse represents a response to a NotifyRequest
type NotifyResponse struct {
	ID   string
	host string
}
