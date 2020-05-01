package main

import (
	"log"
	"os"

	"github.com/dung13890/go-deployer/handlers"
)

func main() {
	// Config Log path
	file, err := os.OpenFile("./logs/deploy.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	handlers.Run()
}
