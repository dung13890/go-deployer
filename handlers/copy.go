package handlers

import (
	"fmt"
	"time"

	"github.com/gosuri/uiprogress"
)

func copy() {
	uiprogress.Start()
	var steps = []string{"downloading source", "installing deps", "compiling", "packaging", "seeding database", "deploying", "staring servers"}
	bar := uiprogress.AddBar(len(steps)).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return "app: " + steps[b.Current()-1]
	})
	for bar.Incr() {
		time.Sleep(time.Millisecond * 400)
	}
	uiprogress.Stop()
	fmt.Println("OK")
}
