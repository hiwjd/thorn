package thorn

import (
	"net/http"
	"log"
	"net"
	"io"
	"encoding/binary"
	"strings"
	"encoding/json"
	"strconv"
	"fmt"
)

// Server 处理这些事：
//   监听服务端口，等待用户请求
//     用户请求与某个客户端建立xxx连接
//   监听管理端口，等待客户端连接
//     告知客户端启动某个端口的pipe
//
// 与客户端交互协议：
//   [1约定头][1版本][2保留][1命令码][4消息体长度][n消息体]
type Server struct {
	config ServerConfig
	buf []byte
	in chan *packetWithConn
	out chan *packetWithClientID
	bufSize int
	cmdHandlers map[byte]CmdHandler
	logic Logic
}

type packetWithConn struct {
	p *Packet
	c net.Conn
}

type packetWithClientID struct {
	p *Packet
	id string
}

func NewServer(config ServerConfig) *Server {
	cmdHandlers := make(map[byte]CmdHandler, 4)
	cmdHandlers[CMD_REG_CLIENT] = RegClientCmdHandler
	s := &Server{
		config: config,
		buf: make([]byte, 1024),
		in: make(chan *packetWithConn, 16),
		out: make(chan *packetWithClientID, 16),
		bufSize: 256,
		cmdHandlers: cmdHandlers,
	}
	s.logic = NewLogic(s)

	return s
}

func (s *Server) Start() {
	go s.startHTTP()
	go s.startTCP()
	go s.ioLoop()
}

func (s *Server) Stop() {
	log.Println("server stop")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.ToLower(r.Method + "#" + r.URL.Path)
	switch path {
	case "post#/openport":
		portStr := r.FormValue("port")
		vportStr := r.FormValue("vport")
		clientID := r.FormValue("clientID")

		port, err := strconv.ParseInt(portStr, 10, 32)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		vport, err := strconv.ParseInt(vportStr, 10, 32)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", vport))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		go func() {
			// 正常情况有且只有两个连接
			//   client的连接
			//   使用方的连接
			// 把这两个连接pipe起来
			for {
				conn1, err := ln.Accept()
				if err != nil {
					log.Println(err)
				}

				body := OpenPort{
					Port:  int(port),
					RemoteAddr: fmt.Sprintf("%s:%v", s.config.SIP(), vport),
				}
				bs, _ := json.Marshal(body)
				p := &Packet{
					Magic: MAGIC,
					Version:  VERSION_1,
					Reserved: uint16(0),
					Cmd:      CMD_OPEN_PORT,
					BodySize: uint32(len(bs)),
					Body:     bs,
				}

				s.out <- &packetWithClientID{p, clientID}

				conn2, err := ln.Accept()
				if err != nil {
					log.Println(err)
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
		}()
	}
}

func (s *Server) startHTTP() {
	log.Fatalln(http.ListenAndServe(s.config.ManageAddr(), s))
}

func (s *Server) startTCP() {
	ln, err := net.Listen("tcp", s.config.ControlAddr())
	if err != nil {
		log.Fatalln(err)
	}

	for {
		var conn net.Conn
		if conn, err = ln.Accept(); err != nil {
			log.Println(err)
			continue
		}
		log.Println("got new client connection")

		go s.readLoop(conn)
	}
}

func (s *Server) ioLoop() {
	for {
		select {
		case p := <-s.in:
			log.Println("process in")
			if h, ok := s.cmdHandlers[p.p.Cmd]; ok {
				h(p.p, p.c, s.logic)
			} else {
				log.Printf("CmdHandler[%d] not found\n", p.p.Cmd)
			}
		case p := <-s.out:
			log.Println("process out")
			conn, err := s.logic.GetClient(p.id)
			log.Println("clientID="+p.id)
			if err != nil {
				log.Println(err)
			} else {
				_, err := conn.Write(p.p.ToBytes())
				if err != nil {
					log.Println(err)
				}
			}
			// sent p to
		}
	}
}

func (s *Server) readLoop(conn net.Conn) {
	buf := make([]byte, s.bufSize)
	for {
		header := buf[0:9]
		if _, err := io.ReadFull(conn, header); err != nil {
			log.Println("read header fail: " + err.Error())
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
		body := buf[0:size]
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
		s.in <- &packetWithConn{p, conn}
	}
}