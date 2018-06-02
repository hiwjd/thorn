package thorn

import (
	"encoding/binary"
)

const (
	MAGIC = 0xF0
	VERSION_1 = 0x01
	CMD_REG_CLIENT = 0x01
	CMD_OPEN_PORT = 0x02
)

type Packet struct {
	Magic byte

	// 协议版本
	Version byte

	// 保留的2字节
	Reserved uint16

	// 命令码
	Cmd byte

	// 包体长度
	BodySize uint32

	// 包体
	Body []byte
}

func (p *Packet) ToBytes() []byte {
	buf := make([]byte, 9 + p.BodySize)
	buf[0] = p.Magic
	buf[1] = p.Version
	buf[2] = byte(p.Reserved >> 8)
	buf[3] = byte(p.Reserved)
	buf[4] = p.Cmd
	binary.BigEndian.PutUint32(buf[5:9], p.BodySize)
	copy(buf[9:], p.Body)

	return buf[:]
}

type RegClient struct {
	ID string
}

type OpenPort struct {
	Port int
	RemoteAddr string
}