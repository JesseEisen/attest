package main

import (
	"fmt"
	"net"
	"time"
	_"io"
	"os"
)

func ServerRoutine(serverType string) {
	defer wg.Done()

	l, err := net.Listen(serverType, address+":"+port)
	if err != nil {
		fmt.Println("Error listening: ", err)
		time.Sleep(time.Second * time.Duration(5))
		os.Exit(1)
	}
	defer l.Close()
	fmt.Printf("Server is Listening on %s...\n", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accept ", err)
			os.Exit(1)
		}

		fmt.Printf("Server: Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

		go handleRequest(conn)
		go listenToSend(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 2048)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		checkReadContent(buf)
	}
}

func checkReadContent(buf []byte) {
	index := int(buf[0]) - 30

	if string(buf) != data[index] {
		fmt.Printf("[server] data not consistent, index: %d\n", index)
	}
}

// This function is used to send content to conn
// it is block by a unbuffered channel. when get the
// notify send related data to client
func listenToSend(conn net.Conn) {
	defer conn.Close()

	for {
		item := <- sendnotify
	
		_, err := conn.Write([]byte(data[item]))
		if err != nil {
			fmt.Println("send data error")
		}
	}
}
