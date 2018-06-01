package main

import (
	"github.com/hiwjd/thorn"
	"os"
	"os/signal"
	"syscall"
	"log"
	"flag"
)

var (
	clientID = flag.String("clientID", "clientID001", "")
	serverAddr = flag.String("serverAddr", "127.0.0.1:9992", "server addr")
)

func init()  {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	flag.Parse()
	config := thorn.ClientConfig{*clientID, *serverAddr}
	client := thorn.NewClient(config)

	client.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sc

	client.Stop()
}
