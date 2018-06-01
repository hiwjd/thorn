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
	manageHTTPAddr = flag.String("mAddr", "127.0.0.1:9991", "http addr provide manage server")
	controlTCPAddr = flag.String("cAddr", "127.0.0.1:9992", "tcp addr provide control server")
)

func init()  {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	flag.Parse()
	config := thorn.NewServerConfig(*manageHTTPAddr, *controlTCPAddr)
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
