package thorn

type ClientConfig struct {
	clientID string
	serverAddr string
}

func NewClientConfig(clientID string, serverAddr string) *ClientConfig {
	return &ClientConfig{
		clientID: clientID,
		serverAddr: serverAddr,
	}
}

func (cc *ClientConfig) ClientID() string {
	return cc.clientID
}

func (cc *ClientConfig) ServerAddr() string {
	return cc.serverAddr
}