all:
	GOOS=windows GOARCH=386  go build main.go serial_port.go case.go server.go

clean:
	rm -rf main.exe
