package thorn

import (
	"net"
	"github.com/pkg/errors"
)

type Logic interface {
	AddClient(clientID string, conn net.Conn) error
	RemoveClient(clientID string) error
	GetClient(clientID string) (net.Conn, error)
	SendPacket(clientID string, p *Packet) error
}

type dftLogic struct {
	s *Server
	cs map[string]net.Conn
}

func NewLogic(s *Server) Logic {
	cs := make(map[string]net.Conn, 4)
	return &dftLogic{s, cs}
}

func (l *dftLogic) AddClient(clientID string, conn net.Conn) error {
	l.cs[clientID] = conn
	return nil
}

func (l *dftLogic) RemoveClient(clientID string) error {
	delete(l.cs, clientID)
	return nil
}

func (l *dftLogic) SendPacket(clientID string, p *Packet) error {
	l.s.out <- &packetWithClientID{p, clientID}
	return nil
}

func (l *dftLogic) GetClient(clientID string) (net.Conn, error) {
	if conn, ok := l.cs[clientID]; ok {
		return conn, nil
	}
	return nil, errors.New("clientID not exists")
}