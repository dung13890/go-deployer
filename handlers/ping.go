package handlers

import (
	"bytes"
	"fmt"
	"log"
	"sync"
)

type ping struct {
	Client []client
}

func (p *ping) load() {

}

func (p *ping) run() {
	c := Configuration{}
	c.ReadFile()
	pathKey := c.GetPathKey()
	wg := sync.WaitGroup{}
	i := 0
	for k, s := range c.Hosts {
		i++
		wg.Add(1)
		go func(w *sync.WaitGroup, s Server, k string, i int) {
			defer w.Done()
			out := bytes.Buffer{}
			r := &remoteScript{}
			r.Stdout = &out
			if err := r.connection(s.Address, s.User, pathKey); err != nil {
				log.Fatalf("Error: %v", err)
			}
			if err := r.run("uname -a"); err != nil {
				log.Fatalf("Error: %v", err)
			}
			r.Color = clientColors[i%len(clientColors)]
			fmt.Printf("%s:%s", fillColor(k, r.Color), fillColor(string(out.Bytes()), r.Color))
			fmt.Printf("%s:%s\n", fillColor(k, r.Color), fillColor(":=====> OK", r.Color))
		}(&wg, s, k, i)
	}
	wg.Wait()
	return
}
