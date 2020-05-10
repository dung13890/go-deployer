package main

import (
	"github.com/dung13890/go-deployer/config"
	"github.com/dung13890/go-deployer/handlers"
)

func main() {
	file := config.LogSetup()
	defer file.Close()
	c := config.Configuration{}
	c.ReadFile()

	handlers.Run(c)
}
