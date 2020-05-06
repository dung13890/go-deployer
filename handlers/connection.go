package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/dung13890/go-deployer/utils"
	"golang.org/x/crypto/ssh"
)

type remoteScript struct {
	Conn       *ssh.Client
	Stdin      io.Reader
	Stdout     bytes.Buffer
	Stderr     bytes.Buffer
	ConnOpened bool
	Color      string
}

func (r *remoteScript) connection(addr string, user string, pathKey string, port ...string) error {
	if r.ConnOpened {
		log.Println("Error: Client already connected")
		return nil
	}
	host := fmt.Sprintf("%s:22", addr)
	if len(port) > 0 {
		host = fmt.Sprintf("%s:%s", addr, port[0])
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
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	r.Conn, err = ssh.Dial("tcp", host, config)
	if err != nil {
		return err
	}
	r.ConnOpened = true

	return nil
}

func (r *remoteScript) run(cmd string) error {
	sess, err := r.Conn.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	sess.Stdout = &r.Stdout
	sess.Stderr = &r.Stderr

	err = sess.Run(cmd)

	return err
}

func (r *remoteScript) close() error {
	if !r.ConnOpened {
		log.Println("Warning: Trying to close the already closed connection")
		return nil
	}
	r.ConnOpened = false
	err := r.Conn.Close()

	return err
}

func (r *remoteScript) showErr(name string, err error) string {
	stdErr := fmt.Sprintf("[%s]: [Failed] %v %v",
		name,
		err,
		string(r.Stderr.Bytes()),
	)
	return utils.FillColor(stdErr, utils.ColorRed)
}

func (r *remoteScript) showOut(name string) string {
	stdOut := fmt.Sprintf("[%s]: [OK] %v",
		name,
		string(r.Stdout.Bytes()),
	)
	return utils.FillColor(stdOut, r.Color)
}
