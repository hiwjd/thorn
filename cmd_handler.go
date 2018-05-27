package thorn

import (
	"encoding/json"
	"net"
)

type CmdHandler func(p *Packet, conn net.Conn, logic Logic) error

func RegClientCmdHandler(p *Packet, conn net.Conn, logic Logic) error {
	var v RegClient
	err := json.Unmarshal(p.Body, &v)
	if err != nil {
		return err
	}
	logic.AddClient(v.ID, conn)
	return nil
}