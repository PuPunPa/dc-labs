package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

func handleDial(port int, bucket string, c chan string) {
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = io.WriteString(conn, bucket)
	if err != nil {
		log.Fatal(err)
	}
	bs := bufio.NewReader(conn)
	line, _, _ := bs.ReadLine()
	c <- string(line)
}

func main() {
	portPtr := flag.Int("proxy", 8000, "proxy to dial")
	bucketPtr := flag.String("bucket", "ryft-public-sample-data", "bucket to request")
	flag.Parse()
	c := make(chan string, 1)
	go handleDial(*portPtr, *bucketPtr, c)
	for item := range c {
		fmt.Println(item)
	}
	close(c)
}
