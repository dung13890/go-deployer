package handlers

import (
	"fmt"
	"log"
	"sync"

	"github.com/dung13890/go-deployer/config"
	"github.com/dung13890/go-deployer/utils"
)

type ping struct {
	Client []client
}

func (p *ping) exec(c config.Configuration) {
	pathKey := c.GetPathKey()
	wg := sync.WaitGroup{}
	rs := make(chan string, 10)
	er := make(chan string, 10)
	i := 0
	for k, s := range c.Hosts {
		i++
		wg.Add(1)
		go func(w *sync.WaitGroup, s config.Server, k string, i int) {
			defer w.Done()
			r := &remoteScript{}
			defer r.close()
			r.Color = utils.ClientColors[i%len(utils.ClientColors)]
			if err := r.connection(s.Address, s.User, pathKey); err != nil {
				er <- r.showErr(k, err)
				log.Print(err)
				return
			}
			r.run("uname -a")
			rs <- r.showOut(k)
		}(&wg, s, k, i)

	}
	wg.Wait()
	close(rs)
	for i := 0; i < len(c.Hosts); i++ {
		select {
		case rv := <-rs:
			fmt.Print(rv)
		case e := <-er:
			fmt.Println(e)
		default:
			fmt.Println("")
		}
	}
	return
}
