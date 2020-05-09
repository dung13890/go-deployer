package handlers

import (
	"github.com/dung13890/go-deployer/config"
)

type remote interface {
	load(name string, server config.Server, index int)
	loadTask(tasks []string)
	makeString(str string) string
	connect(pathKey string, port ...string) error
	shell(cf callbackFunc) error
	run(cmd string) error
	showOut(in string, out string) string
	muxShell(cf callbackFunc) error
	close() error
}

type callbackFunc func(out string)
