package handlers

import (
	"fmt"
	"log"
	"sync"

	"github.com/dung13890/go-deployer/config"
)

type ping struct {
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
			var h remote = &host{}
			defer h.close()
			h.load(k, s, i)
			if err := h.connect(pathKey); err != nil {
				er <- fmt.Sprintf("[%s] [Failed]", k)
				log.Print(err)
				return
			}
			h.run("uname -a")
			rs <- fmt.Sprintf("[%s] OK!", k)
		}(&wg, s, k, i)

	}
	wg.Wait()
	for i := 0; i < len(c.Hosts); i++ {
		select {
		case rv := <-rs:
			fmt.Println(rv)
		case e := <-er:
			fmt.Println(e)
		default:
			fmt.Println("")
		}
	}
	close(rs)
	close(er)
	return
}
