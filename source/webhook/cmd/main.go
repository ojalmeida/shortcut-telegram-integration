package main

import (
	"os"
	"os/signal"
	"webhook/server"
)

func main() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go server.Start()

	<-ch
	server.Stop()

}
