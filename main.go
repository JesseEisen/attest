package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup


func main() {
	wg.Add(2)

	content = make(chan string)
	fmt.Printf("start test | All %d items\n", len(casetable))
	go ServerRoutine("tcp")
	go AtTestRoutine()

	wg.Wait()
}
