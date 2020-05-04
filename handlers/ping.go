package handlers

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/dung13890/go-deployer/config"
	"github.com/dung13890/go-deployer/utils"
)

type ping struct {
	Client []client
}

func pingClient(w *sync.WaitGroup, s config.Server, k string, i int, pathKey string) {
	defer w.Done()
	out := bytes.Buffer{}
	r := &remoteScript{}
	r.Color = utils.ClientColors[i%len(utils.ClientColors)]
	r.Stdout = &out
	if err := r.connection(s.Address, s.User, pathKey); err != nil {
		log.Fatalf("Error: %v", err)
	}
	if err := r.run("uname -a"); err != nil {
		log.Fatalf("Error: %v", err)
	}
	stdOut := fmt.Sprintf("[%s]: %s\r[%s]: Status OK",
		k,
		string(out.Bytes()),
		k,
	)
	fmt.Println(utils.FillColor(stdOut, r.Color))
}

func (p *ping) exec(c config.Configuration) {
	pathKey := c.GetPathKey()
	wg := sync.WaitGroup{}
	i := 0
	for k, s := range c.Hosts {
		i++
		wg.Add(1)
		go pingClient(&wg, s, k, i, pathKey)
	}
	wg.Wait()
	return
}
