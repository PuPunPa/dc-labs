package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type ListBucketResult struct {
	Contents []Contents `xml:"Contents"`
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
		log.Fatal("ERROR1: " + err.Error())
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ERROR2: " + err.Error())
	}

	var body ListBucketResult
	xml.Unmarshal(bodyBytes, &body)

	content, err := json.Marshal(body.Contents)
	if err != nil {
		log.Fatal("ERROR3: " + err.Error())
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
	fmt.Println(string(send))
	if err != nil {
		log.Fatal("ERROR4" + err.Error())
	}
	return send
}

func handleConn(c net.Conn) {
	defer c.Close()
	bs := bufio.NewReader(c)
	bucket := getMessage(bs)
	_, err := io.WriteString(c, string(writeJson(string(bucket))))
	if err != nil {
		log.Fatal("ERROR5: " + err.Error())
		return // e.g., client disconnected
	}
	fmt.Println("Sent")
}

func getMessage(bs *bufio.Reader) []byte {
	var message []byte
	prefix := true
	for prefix {
		message, prefix, _ = bs.ReadLine()
	}
	return message
}

func main() {
	portPtr := flag.Int("port", 9090, "port to be used")
	flag.Parse()
	portStr := "localhost:" + strconv.Itoa(*portPtr)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Fatal("ERROR7: " + err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("ERROR8: " + err.Error()) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
}
