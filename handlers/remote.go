package handlers

import (
	"github.com/dung13890/go-deployer/config"
)

type remote interface {
	load(name string, server config.Server, index int)
	loadTask(tasks []string)
	connect(pathKey string, port ...string) error
	shell() error
	run(cmd string) error
	printOut(in string, out string)
	muxShell() error
	close() error
}
