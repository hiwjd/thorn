package main

import (
	"github.com/hiwjd/thorn"
	"os"
	"os/signal"
	"syscall"
	"log"
	"flag"
	"github.com/teris-io/shortid"
	"os/user"
	"gopkg.in/ini.v1"
)

var (
	clientID = flag.String("clientID", "", "specify a unique client id or left blank to auto generate one")
	serverAddr = flag.String("serverAddr", "127.0.0.1:9992", "server addr")
)

func init()  {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	flag.Parse()

	if *clientID == "" {
		// look up from ~/.thorn file
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		iniPath := usr.HomeDir + "/.thorn"
		if _, err := os.Stat(iniPath); os.IsNotExist(err) {
			f, err := os.Create(iniPath)
			if err != nil {
				log.Fatal(err)
			}
			f.Close()
		}

		cfg, err := ini.Load(iniPath)
		if err != nil {
			log.Fatal(err)
		}

		*clientID = cfg.Section("").Key("clientID").String()

		// auto generate one
		if *clientID == "" {
			*clientID = shortid.MustGenerate()
			cfg.Section("").Key("clientID").SetValue(*clientID)
			cfg.SaveTo(iniPath)
		}
	}

	config := thorn.NewClientConfig(*clientID, *serverAddr)
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
