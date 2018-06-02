package thorn

import (
	"net"
	"log"
	"encoding/json"
	"io"
	"encoding/binary"
	"fmt"
)

// Client 做这些事情：
//   与服务端建立连接
//   处理服务端下发的命令
//     与本地的某个端口建立连接后，连接pipe给服务端指定的端口
type Client struct {
	config *ClientConfig
	buf []byte
	in chan *Packet
}

func NewClient(config *ClientConfig) *Client {
	return &Client{
		config: config,
		buf: make([]byte, 512),
		in: make(chan *Packet, 16),
	}
}

func (c *Client) Start() {
	conn, err := net.Dial("tcp", c.config.ServerAddr())
	if err != nil {
		log.Println(err)
		return
	}

	body := RegClient{c.config.ClientID()}
	bs, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return
	}

	p := Packet{
		Magic: MAGIC,
		Version: VERSION_1,
		Reserved: 0,
		Cmd: CMD_REG_CLIENT,
		BodySize: uint32(len(bs)),
		Body: bs,
	}
	_, err = conn.Write(p.ToBytes())
	if err != nil {
		log.Println(err)
		return
	}

	go c.readLoop(conn)
	go c.ioLoop(conn)
}

func (c *Client) Stop() {
	log.Println("client stop")
}

func (c *Client) readLoop(conn net.Conn) {
	for {
		header := c.buf[0:9]
		_, err := io.ReadFull(conn, header)
		if err != nil {
			log.Println(err)
			break
		}

		if header[0] != MAGIC {
			log.Println("magic wrong")
			break
		}

		version := header[1]
		reserved := binary.BigEndian.Uint16(header[2:4])
		cmd := header[4]
		size := binary.BigEndian.Uint32(header[5:9])
		body := c.buf[0:size]
		if _, err := io.ReadFull(conn, body); err != nil {
			log.Println("read body fail: " + err.Error())
			break
		}

		p := &Packet{
			Magic: MAGIC,
			Version: version,
			Reserved: reserved,
			Cmd: cmd,
			BodySize: size,
			Body: body,
		}
		log.Printf("got a new packet: %+v\n", p)
		c.in <- p
	}
}

func (c *Client) ioLoop(conn net.Conn) {
	for {
		select {
		case p := <- c.in:
			{
				log.Println(p.Cmd)
				if p.Cmd == CMD_OPEN_PORT {
					var b OpenPort
					if err := json.Unmarshal(p.Body, &b); err != nil {
						log.Println(err)
						return
					}

					conn1, err := net.Dial("tcp", b.RemoteAddr)
					if err != nil {
						log.Println(err)
						return
					}
					conn2, err := net.Dial("tcp", fmt.Sprintf(":%d", b.Port))
					if err != nil {
						log.Println(err)
						return
					}
					go func() {
						if _, err := io.Copy(conn1, conn2); err != nil {
							log.Println(err)
						}
					}()
					go func() {
						if _, err := io.Copy(conn2, conn1); err != nil {
							log.Println(err)
						}
					}()
				}
			}
		}
	}
}