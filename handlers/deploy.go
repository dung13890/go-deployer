package handlers

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/dung13890/go-deployer/config"
	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
)

const (
	dir string = "/data/sites"
)

type deploy struct {
	pathKey string
	hosts   map[string]config.Server
	tag     string
	tasks   []string
	branch  string
	logged  bool
}

func deploySetup(c config.Configuration, tag string, branch string, logged bool) *deploy {
	pathKey := c.GetPathKey()
	d := &deploy{
		hosts:   c.Hosts,
		pathKey: pathKey,
		tag:     tag,
		branch:  branch,
		logged:  logged,
	}
	d.loadTask(c)

	return d
}

func (d *deploy) loadTask(c config.Configuration) []string {
	projectName := strings.Replace(c.Setting.Name, " ", "", -1)
	shareDir := fmt.Sprintf("%s/%s/shared", dir, projectName)
	tasks := []string{
		"sudo mkdir -p " + shareDir,
		"sudo chown -R $USER:$USER " + shareDir,
		"cd " + shareDir,
	}
	d.tasks = append(tasks, c.Tasks...)

	return tasks
}

func (d *deploy) exec() {
	wg := sync.WaitGroup{}
	out := make(chan string, 10)
	arrOut := make([][]string, len(d.hosts))
	uiprogress.Start()
	i := 0
	for k, s := range d.hosts {
		wg.Add(1)
		go d.running(&wg, s, k, i, out, &arrOut[i])
		i += 1
	}
	wg.Wait()
	uiprogress.Stop()
	for i := 0; i < len(d.hosts); i++ {
		select {
		case o := <-out:
			fmt.Println(o)
		default:
			fmt.Println()
		}
	}
	if d.logged {
		for _, loS := range arrOut {
			for _, lo := range loS {
				fmt.Print(lo)
			}
		}
	}
	close(out)
}

func (d *deploy) running(wg *sync.WaitGroup, se config.Server, name string, index int, ch chan string, strOut *[]string) {
	defer wg.Done()
	// Make Interface remote
	var h remote = &host{}
	h.load(name, se, index)

	// Setting progress bar
	bar := uiprogress.AddBar(len(d.tasks)).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		strCommand := "init"
		if b.Current() != 0 {
			strCommand = d.tasks[b.Current()-1]
		}
		return strutil.Resize(h.makeString(strCommand), 30)
	})

	h.loadTask(d.tasks)
	// Connect Server
	if err := h.connect(d.pathKey); err != nil {
		log.Print(err)
		return
	}
	h.shell(func(out string) {
		bar.Incr()
		*strOut = append(*strOut, out)
	})
	ch <- fmt.Sprintf("[%s] SUCCESSES!", name)
	h.close()
}
