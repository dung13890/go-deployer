package handlers

import (
	"fmt"
	"log"
	"sync"

	"github.com/dung13890/go-deployer/config"
)

type deploy struct {
}

func (d *deploy) exec(c config.Configuration) {
	pathKey := c.GetPathKey()
	wg := sync.WaitGroup{}
	rs := make(chan string, 10)
	i := 0
	for k, s := range c.Hosts {
		i++
		wg.Add(1)
		go func(w *sync.WaitGroup, s config.Server, k string, i int) {
			defer wg.Done()
			var h remote = &host{}
			defer h.close()
			h.load(k, s, i)
			h.loadTask([]string{
				"uname -a",
				"sleep 5",
				"ls -al",
				"exit",
			})
			if err := h.connect(pathKey); err != nil {
				log.Print(err)
				return
			}

			h.shell()
			rs <- fmt.Sprintf("[%s] SUCCESSES!", k)
		}(&wg, s, k, i)
	}
	wg.Wait()
	for i := 0; i < len(c.Hosts); i++ {
		select {
		case rv := <-rs:
			fmt.Println(rv)
		default:
			fmt.Println("")
		}
	}
	close(rs)
}
