package main

import (
	"fmt"
	"io"
	"net"
	"os"
	_ "strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

const (
	address string = "127.0.0.1"
	port    string = "1234"
	name    string = "COM9"
	baud    int    = 115200
)

var (
	wg       sync.WaitGroup
	content  chan string
	passitem int
)

func main() {
	wg.Add(2)

	content = make(chan string)
	fmt.Printf("start test | All %d items\n", len(casetable))
	go ServerRoutine("tcp")
	go AtTestRoutine()

	wg.Wait()
}

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
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		io.Copy(conn, conn)
	}
}

func AtTestRoutine() {
	defer wg.Done()

	fmt.Println("start open COM")
	time.Sleep(time.Second * time.Duration(12))

	s := OpenCom()

	for _, cmd := range casetable {
		ExecAT(s, &cmd)
		CheckRes(&cmd)
	}

	fmt.Println("=========== Summary ============")
	fmt.Printf("Pass: %d | Failed: %d\n", passitem, len(casetable)-passitem)
	time.Sleep(time.Second * time.Duration(5))
}

func CheckRes(cmd *atcmd) {
	res := <-content
	if cmd.hasURC {
		// here will block to wait the qiopen urc
		res = <-content
		if cmd.URCContent != res {
			fmt.Printf("Failed! CMD:%q | Expect: %q Get: %q\n",
				cmd.command,
				cmd.expect,
				string(res))
		}
		return
	}

	if string(res) != cmd.expect {
		fmt.Printf("Failed! CMD:%q | Expect: %q Get: %q\n",
			cmd.command,
			cmd.expect,
			string(res))
	}

	passitem++
}

func OpenCom() *serial.Port {
	c := &serial.Config{
		Name: name,
		Baud: baud,
	}

	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println("Open error: ", err)
		os.Exit(1)
	}

	go ReadCOM(s)

	return s
}

func ExecAT(s *serial.Port, cmd *atcmd) {

	_, err := s.Write([]byte(cmd.command))
	if err != nil {
		fmt.Println(err)
	}

}

// ReadCOM should always read the serial port
// but here may contain a porblem: if there are more
// than one piece of data come to port, the result may become
// dirty.so the check function will be more complicated.
func ReadCOM(s *serial.Port) {

	for {
		buf := make([]byte, 2048)
		n, err := s.Read(buf)
		if err != nil {
			fmt.Println(err)
		}

		content <- string(buf[:n])
	}
}
