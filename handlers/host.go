package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"github.com/dung13890/go-deployer/config"
	"github.com/dung13890/go-deployer/utils"

	"golang.org/x/crypto/ssh"
)

type host struct {
	Name       string
	Server     config.Server
	Tasks      []string
	Conn       *ssh.Client
	Stdin      io.WriteCloser
	Stdout     io.Reader
	Stderr     io.Reader
	ConnOpened bool
	Color      string
}

func (h *host) load(name string, server config.Server, index int) {
	h.Name = name
	h.Server = server
	h.Color = utils.ClientColors[index%len(utils.ClientColors)]
}

func (h *host) loadTask(tasks []string) {
	h.Tasks = tasks
}

func (h *host) connect(pathKey string, port ...string) error {
	if h.ConnOpened {
		log.Println("Error: Client already connected")
		return nil
	}
	addr := fmt.Sprintf("%s:22", h.Server.Address)
	if len(port) > 0 {
		addr = fmt.Sprintf("%s:%s", h.Server.Address, port[0])
	}
	key, err := ioutil.ReadFile(pathKey)
	if err != nil {
		log.Println("Error: Have Not private key")
		return err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Error: Wrong format of private key")
		return err
	}
	config := &ssh.ClientConfig{
		User: h.Server.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	h.Conn, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	h.ConnOpened = true

	return nil
}

func (h *host) shell() error {
	sess, err := h.Conn.NewSession()
	defer sess.Close()
	if err != nil {
		return err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err = sess.RequestPty("vt220", 80, 40, modes); err != nil {
		return err
	}

	if h.Stdin, err = sess.StdinPipe(); err != nil {
		return err
	}
	if h.Stdout, err = sess.StdoutPipe(); err != nil {
		return err
	}
	if h.Stderr, err = sess.StderrPipe(); err != nil {
		return err
	}
	if err = sess.Shell(); err != nil {
		return err
	}
	h.muxShell()
	return sess.Wait()
}

func (h *host) run(cmd string) error {
	sess, err := h.Conn.NewSession()
	defer sess.Close()
	if err != nil {
		return err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err = sess.RequestPty("vt220", 80, 40, modes); err != nil {
		return err
	}

	if h.Stdout, err = sess.StdoutPipe(); err != nil {
		return err
	}
	if h.Stderr, err = sess.StderrPipe(); err != nil {
		return err
	}

	if err = sess.Run(cmd); err != nil {
		return err
	}

	buf := [65 * 1024]byte{}
	n, _ := h.Stdout.Read(buf[:])
	h.printOut(cmd, string(buf[:n]))

	return nil
}

func (h *host) printOut(in string, out string) {
	stdOut := fmt.Sprintf("[%s] RUN %s \n%s\n", h.Name, in, out)
	fmt.Print(utils.FillColor(stdOut, h.Color))
}

func (h *host) muxShell() error {
	in := make(chan string, 10)
	out := make(chan string, 10)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for cmd := range in {
			io.WriteString(h.Stdin, cmd+"\r")
		}
		wg.Done()
	}()
	go func() {
		buf := [65 * 1024]byte{}
		t := 0
		for {
			n, err := h.Stdout.Read(buf[t:])
			if err == io.EOF {
				close(in)
				close(out)
				break
			}
			t += n
			if buf[t-2] == '$' {
				out <- string(buf[:t])
				t = 0
			}
		}
		wg.Done()
	}()
	<-out
	for _, cmd := range h.Tasks {
		in <- cmd
		h.printOut(cmd, <-out)
	}

	wg.Wait()

	return nil
}

func (h *host) close() error {
	if !h.ConnOpened {
		log.Println("Warning: Trying to close the already closed connection")
		return nil
	}
	h.ConnOpened = false
	err := h.Conn.Close()

	return err
}
