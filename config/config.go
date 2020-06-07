package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Hosts   map[string]Server `yaml:"hosts"`
	Tasks   []string          `yaml:"tasks"`
	Setting Setting           `yaml:"setting"`
}

type Server struct {
	Address string `yaml:"address"`
	User    string `yaml:"user"`
}

type Setting struct {
	PathKey string `yaml:"pathKey"`
	PathEnv string `yaml:"pathEnv"`
	Repo    string `yaml:"repo"`
	Name    string `yaml:"name"`
}

func (c *Configuration) ReadFile() {
	file, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal("Error: loading yml file")
	}
	errY := yaml.Unmarshal(file, &c)
	if errY != nil {
		log.Fatalf("Error: %v", errY)
	}

	if c.Setting.Repo == "" || c.Setting.Name == "" {
		log.Fatalf("Error: Please config repository")
	}
}

func (c *Configuration) GetPathKey() string {
	path := c.Setting.PathKey
	if path == "" {
		return filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}

	replacePath := strings.Replace(path, "~", os.Getenv("HOME"), 1)
	pathKey, err := filepath.Abs(replacePath)

	if err != nil {
		log.Fatal("Warning: path file of rsa key is not exists")
	}

	return pathKey
}
