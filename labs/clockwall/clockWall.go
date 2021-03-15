// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func handleDial(c chan string, port int) {
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	bs := bufio.NewReader(conn)
	line, _, _ := bs.ReadLine()
	defer wg.Done()
	c <- string(line)
}

func main() {
	portPtr1 := flag.Int("NewYork", 9090, "port to be used")
	portPtr2 := flag.Int("Tokyo", 9090, "port to be used")
	portPtr3 := flag.Int("London", 9090, "port to be used")
	flag.Parse()
	c := make(chan string, 3)
	wg.Add(1)
	go handleDial(c, *portPtr1)
	wg.Add(1)
	go handleDial(c, *portPtr2)
	wg.Add(1)
	go handleDial(c, *portPtr3)
	wg.Wait()
	close(c)
	for item := range c {
		fmt.Println(item)
	}
}
