// https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

func main() {
	var u =  flag.String("user", "", "Username")
	var pw = flag.String("password", "", "Password")
	var f  = flag.String("file", "", "File to copy")
	var h = flag.String("host", "", "Host name")
	flag.Parse()
	
	clientConfig := &ssh.ClientConfig{
		User: *u,
		Auth: []ssh.AuthMethod{
			ssh.Password(*pw),
		},
	}

	var hoststring bytes.Buffer
	hoststring.WriteString(*h)
	hoststring.WriteString(":22")

	client, err := ssh.Dial("tcp", hoststring.String(), clientConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		
		file, err := os.Open(*f)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			// Could not obtain stat, handle error
		}

		fmt.Fprintln(w, "D0755", 0, "testdir") // mkdir
		fmt.Fprintln(w, "C0644", fi.Size(), "testfile1")
		io.Copy(w, file);
		fmt.Fprint(w, "\x00") // transfer end with \x00
	}()
	if err := session.Run("/usr/bin/scp -tr ./"); err != nil {
		panic("Failed to run: " + err.Error())
	}
}