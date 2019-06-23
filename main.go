package main

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

const sshUser = ""
const sshPass = ""
const domainWithPort = "" // ex: google.com:22

func main() {
	executeCommandPOC()
}

func executeCommandPOC() {
	conn := getConnectionSSH()
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("cd /var/www/html && ls"); err != nil {
		fmt.Println("Failed to run: " + err.Error())
	}

	if b.Len() == 0 {
		fmt.Println("No data!")
	}

	fmt.Println(b.String())
}

func getConnectionSSH() *ssh.Client {
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", domainWithPort, config)

	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	return conn
}
