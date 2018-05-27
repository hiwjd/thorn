package main

import (
	"github.com/hiwjd/thorn"
	"os/signal"
	"os"
	"syscall"
	"log"
)

func init()  {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	config := thorn.NewServerConfig("127.0.0.1:9991", "127.0.0.1:9992")
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
