package main

import (
	"github.com/hiwjd/thorn"
	"os"
	"os/signal"
	"syscall"
	"log"
)

func init()  {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	config := thorn.ClientConfig{"clientID001", "127.0.0.1:9992"}
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
