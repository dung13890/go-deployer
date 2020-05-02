package handlers

import "io"

type client interface {
	Connect(hot string) error
	Run(tasks []string)
	Wait() error
	Close() error
	Stderr() io.Reader
	Stdout() io.Reader
}
