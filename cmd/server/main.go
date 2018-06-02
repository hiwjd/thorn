package main

import (
	"github.com/hiwjd/thorn"
	"os/signal"
	"os"
	"syscall"
	"log"
	"flag"
)

var (
	sip = flag.String("sip", "127.0.0.1", "server ip")
	mPort = flag.Int("mPort", 9991, "http port provide manage server")
	cPort = flag.Int("cPort", 9992, "tcp addr provide control server")
)

func init()  {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	flag.Parse()
	config := thorn.NewServerConfig(*sip, *mPort, *cPort)
	server := thorn.NewServer(config)

	server.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sc

	server.Stop()
}
