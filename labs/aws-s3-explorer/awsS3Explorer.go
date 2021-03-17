package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ListBucketResult struct {
	Contents []Contents `xml:"Contents`
}

type Contents struct {
	Key string `xml:"Key"`
}

type Send struct {
	BucketName       string
	ObjectsCount     int
	DirectoriesCount int
	Extensions       map[string]int
}

func writeJson(bucket string) []byte {
	resp, err := http.Get("http://s3.amazonaws.com/" + bucket)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var body ListBucketResult
	xml.Unmarshal(bodyBytes, &body)

	content, err := json.Marshal(body.Contents)
	if err != nil {
		log.Fatal(err)
	}

	contentStr := string(content[1 : len(content)-1])
	keys := strings.Split(contentStr, ",")

	objsCnt := len(keys)
	dirs := make(map[string]int)
	exts := make(map[string]int)

	for i := 0; i < objsCnt; i++ {
		key := string(keys[i][8 : len(keys[i])-2])
		dir := strings.Split(key, "/")
		if len(dir) > 1 {
			dirs[dir[0]] = 1
		}
		ext := strings.Split(key, ".")
		if len(ext) > 1 {
			exts[ext[len(ext)-1]] += 1
		}
	}
	s := Send{bucket, objsCnt, len(dirs), exts}
	send, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return send
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		bs := bufio.NewReader(c)
		bucket, _, _ := bs.ReadLine()
		_, err := io.WriteString(c, string(writeJson(string(bucket))))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	portPtr := flag.Int("port", 9090, "port to be used")
	flag.Parse()
	portStr := "localhost:" + strconv.Itoa(*portPtr)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
}
