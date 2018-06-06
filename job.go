package thorn

// Job represents a proxy job
type Job struct {
	// Job ID
	ID int64

	// Job name
	Name string

	// server ip
	ServerIP string

	// virtual port
	VirtualPort int

	// port to be proxy
	Port int

	// on which client port to be proxy
	ClientID string
}