package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const sshUser = ""
const sshPass = ""
const domainWithPort = "" // ex: google.com:22
const remotePath = "/var/www/html/"

func main() {
	files := getTxtFilesFromRemoteServer()
	getFilesFromRemoteServer(files)
}

func getTxtFilesFromRemoteServer() []string {
	conn := getConnectionSSH()
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(`cd ` + remotePath + ` && ls *.txt`); err != nil {
		fmt.Println("Failed to run: " + err.Error())
	}

	files := strings.Split(b.String(), "\n")
	return files
}

func getFilesFromRemoteServer(files []string) {
	conn := getConnectionSSH()
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	dstPath, _ := os.Getwd()
	pathErr := os.MkdirAll(strings.Join([]string{dstPath, "data"}, "/"), 0777)

	if pathErr != nil {
		fmt.Println(pathErr)
	}

	for _, file := range files {
		if len(file) > 0 {
			// create destination file
			dstFile, err := os.Create(strings.Join([]string{dstPath, "data", file}, "/"))
			if err != nil {
				log.Fatal(err)
			}
			defer dstFile.Close()

			// open source file
			srcFile, err := client.Open(remotePath + file)
			if err != nil {
				log.Fatal(err)
			}

			// copy source file to destination file
			bytes, err := io.Copy(dstFile, srcFile)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s: %d bytes copied\n", file, bytes)

			// flush in-memory copy
			err = dstFile.Sync()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
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
