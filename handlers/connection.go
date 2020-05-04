package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

type remoteScript struct {
	Conn       *ssh.Client
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	ConnOpened bool
	Color      string
}

func (r *remoteScript) connection(addr string, user string, pathKey string, port ...string) error {
	if r.ConnOpened {
		log.Fatal("Error: Client already connected")
	}
	host := fmt.Sprintf("%s:22", addr)
	if len(port) > 0 {
		host = fmt.Sprintf("%s:%s", addr, port[0])
	}

	key, err := ioutil.ReadFile(pathKey)
	if err != nil {
		log.Fatal("Error: Have Not private key")
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("Error: Wrong format of private key")
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
	sess.Stdout = r.Stdout
	sess.Stderr = r.Stderr

	err = sess.Run(cmd)

	return err
}

func (r *remoteScript) close() error {
	if !r.ConnOpened {
		log.Fatal("Warning: Trying to close the already closed connection")
	}
	r.ConnOpened = false
	err := r.Conn.Close()

	return err
}

func test() {
	c := Configuration{}
	c.ReadFile()
	pathKey := c.Setting.PathKey
	if pathKey != "" {
		pathKey = strings.Replace(pathKey, "~", os.Getenv("HOME"), 1)
		pathKey, err := filepath.Abs(pathKey)
		fmt.Println(pathKey)
		if err != nil {
			log.Fatal("Warning: path file of rsa key is not exists")
		}
	} else {
		pathKey = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}

	fmt.Println(pathKey)
	for k, s := range c.Hosts {
		out := bytes.Buffer{}
		r := &remoteScript{}
		r.Stdout = &out
		if err := r.connection(s.Address, s.User, pathKey); err != nil {
			log.Fatalf("Error: %v", err)
		}
		r.run("uname -a")
		r.Color = colorMagenta
		fmt.Printf("%s:%s", fillColor(k, r.Color), fillColor(string(out.Bytes()), r.Color))
	}
}
