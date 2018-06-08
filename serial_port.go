package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	
	"github.com/tarm/serial"
)

const (
	address string = "127.0.0.1"
	port    string = "1234"
	name    string = "COM9"
	baud    int    = 115200
)

var (
	content    chan string
	sendnotify chan int
	passitem   int
	failitem   int
)


func AtTestRoutine() {
	defer wg.Done()

	fmt.Println("start open COM")
	time.Sleep(time.Second * time.Duration(12))

	s := OpenCom()

	for _, cmd := range casetable {
		ExecAT(s, &cmd)
		CheckRes(s, &cmd)
	}

	fmt.Println("=========== Summary ============")
	fmt.Printf("Pass: %d | Failed: %d\n", passitem, failitem)
	time.Sleep(time.Second * time.Duration(5))
}

func CheckRes(s *serial.Port, cmd *atcmd) {
	res := <-content
	
	// if is urc, we just skip it
	if strings.Contains(res, "URC") {
		return
	}

	if cmd.hasURC {
		// here will block to wait the qiopen urc
		res = <-content
		if cmd.urcContent != res {
			fmt.Printf("Failed! CMD:%q | Expect: %q Get: %q\n",
				cmd.command,
				cmd.expect,
				res)
			failitem++
			return
		}
	}

	if res != cmd.expect {
		fmt.Printf("Failed! CMD:%q | Expect: %q Get: %q\n",
			cmd.command,
			cmd.expect,
			res)
		failitem++
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
	
	fmt.Printf("[CMD] %s\n", cmd.command)
	_, err := s.Write([]byte(cmd.command))
	if err != nil {
		fmt.Println(err)
	}

	if cmd.hasSend {
		_, err := s.Write([]byte(cmd.sendContent))
		if err != nil {
			fmt.Println(err)
		}
	}

	// notify server to send data to UE
	if cmd.bnextrd {
		sendnotify <- cmd.readitem
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
