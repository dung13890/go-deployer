package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/dung13890/go-deployer/config"
	"github.com/dung13890/go-deployer/utils"

	"golang.org/x/crypto/ssh"
)

type host struct {
	name       string
	server     config.Server
	tasks      []string
	conn       *ssh.Client
	stdin      io.WriteCloser
	stdout     io.Reader
	stderr     io.Reader
	connOpened bool
	color      string
}

func (h *host) load(name string, server config.Server, index int) {
	h.name = name
	h.server = server
	h.color = utils.ClientColors[index%len(utils.ClientColors)]
}

func (h *host) loadTask(tasks []string) {
	h.tasks = tasks
}

func (h *host) makeString(str string) string {
	pre := fmt.Sprintf("[%s] %s", h.name, str)
	return utils.FillColor(pre, h.color)
}

func (h *host) connect(pathKey string, port ...string) error {
	if h.connOpened {
		log.Println("Error: Client already connected")
		return nil
	}
	addr := fmt.Sprintf("%s:22", h.server.Address)
	if len(port) > 0 {
		addr = fmt.Sprintf("%s:%s", h.server.Address, port[0])
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
		User: h.server.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	h.conn, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	h.connOpened = true

	return nil
}

func (h *host) shell(cf callbackFunc) error {
	sess, err := h.conn.NewSession()
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

	if h.stdin, err = sess.StdinPipe(); err != nil {
		return err
	}
	if h.stdout, err = sess.StdoutPipe(); err != nil {
		return err
	}
	if h.stderr, err = sess.StderrPipe(); err != nil {
		return err
	}
	if err = sess.Shell(); err != nil {
		return err
	}
	h.muxShell(cf)
	return sess.Wait()
}

func (h *host) run(cmd string) error {
	sess, err := h.conn.NewSession()
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
	if h.stdout, err = sess.StdoutPipe(); err != nil {
		return err
	}
	if h.stderr, err = sess.StderrPipe(); err != nil {
		return err
	}

	if err = sess.Run(cmd); err != nil {
		return err
	}

	buf := [65 * 1024]byte{}
	n, _ := h.stdout.Read(buf[:])
	fmt.Print(h.showOut(cmd, string(buf[:n])))

	return nil
}

func (h *host) showOut(in string, out string) string {
	stdOut := fmt.Sprintf("%s\n%s", in, out)
	return h.makeString(stdOut)
}

func (h *host) muxShell(cf callbackFunc) error {
	in := make(chan string, 10)
	out := make(chan string, 10)
	errC := make(chan string, 10)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for cmd := range in {
			io.WriteString(h.stdin, cmd+"\n")
		}
		wg.Done()
	}()
	go func() {
		buf := [65 * 1024]byte{}
		t := 0
		for {
			n, err := h.stdout.Read(buf[t:])
			if err == io.EOF {
				close(in)
				close(out)
				// close(errC)
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

	go func() {
		io.Copy(os.Stderr, h.stderr)
	}()
	<-out
	for _, cmd := range h.tasks {
		in <- cmd
		select {
		case stdOut := <-out:
			cf(h.showOut(cmd, stdOut))
		case stdError := <-errC:
			fmt.Println(stdError)
		}
	}

	wg.Wait()

	return nil
}

func (h *host) close() error {
	if !h.connOpened {
		log.Println("Warning: Trying to close the already closed connection")
		return nil
	}
	h.connOpened = false
	err := h.conn.Close()

	return err
}
