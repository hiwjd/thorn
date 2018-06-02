package thorn

import "fmt"

type ServerConfig struct {
	// server ip
	sip string

	// http port for manage
	mPort  int

	// tcp port for client to connect
	cPort int
}

func NewServerConfig(sip string, mPort, cPort int) ServerConfig {
	return ServerConfig{
		sip:  sip,
		mPort: mPort,
		cPort: cPort,
	}
}

func (sc *ServerConfig) SIP() string {
	return sc.sip
}

func (sc *ServerConfig) ManageAddr() string {
	return fmt.Sprintf("%s:%d", sc.sip, sc.mPort)
}

func (sc *ServerConfig) ControlAddr() string {
	return fmt.Sprintf("%s:%d", sc.sip, sc.cPort)
}