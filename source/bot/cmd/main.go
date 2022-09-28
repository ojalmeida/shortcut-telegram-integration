package main

import (
	"bot/lib/telegram"
	"os"
	"os/signal"
)

func main() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go telegram.Start()

	<-ch

}
