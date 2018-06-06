package thorn

import (
	"net"
	"fmt"
	"log"
	"encoding/json"
	"io"
)

type JobRunner interface {
	Run(*Job) error
}

func NewJobRunner(server *Server) JobRunner {
	return &jobRunner{
		server: server,
	}
}

type jobRunner struct {
	server *Server
}

func (j *jobRunner) Run(job *Job) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", job.VirtualPort))
	if err != nil {
		return err
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
				Port:  job.Port,
				RemoteAddr: fmt.Sprintf("%s:%v", j.server.config.SIP(), job.VirtualPort),
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

			j.server.out <- &packetWithClientID{p, job.ClientID}

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
	return nil
}