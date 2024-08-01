package main

import (
	"log"
	"os"
	"reverse-phone-lookup/bot"
	"reverse-phone-lookup/config"
)

func main() {
	cwd, err := os.Getwd()
	err = config.ReadConfig(cwd)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()

	<-make(chan struct{})
	return
}
