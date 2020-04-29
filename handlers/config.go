package handlers

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	WebServers
}

type WebServers struct {
	Hosts map[string]Server `yaml:"hosts"`
	Tasks []string          `yaml:"Tasks"`
}

type Server struct {
	Address string `yaml:"address"`
	User    string `yaml:"user"`
	Dir     string `yaml:"dir"`
}

func (c *Configuration) ReadFile() {
	file, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal("Error loading yml file")
	}
	errY := yaml.Unmarshal(file, &c)
	if errY != nil {
		log.Fatalf("error: %v", errY)
	}
}
