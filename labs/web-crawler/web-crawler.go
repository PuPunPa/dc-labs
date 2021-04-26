// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

//!-sema

//!+
func main() {
	worklist := make(chan []string)
	depthPtr := flag.Int("depth", 0, "depth to crawl")
	resultsPtr := flag.String("results", "results.txt", "save to txt file")
	flag.Parse()
	depth := *depthPtr
	results := *resultsPtr
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++

	argSize := len(os.Args)
	go func() { worklist <- os.Args[argSize-1:] }()

	// make output file
	fo, err := os.Create(results)
	if err != nil {
		log.Fatal(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		if depth > 0 {
			list := <-worklist
			for _, link := range list {
				if !seen[link] {
					seen[link] = true
					n++
					go func(link string) {
						worklist <- crawl(link)
					}(link)

					// write a chunk
					if _, err := fo.Write([]byte(link)); err != nil {
						log.Fatal(err)
					}
					if _, err := fo.Write([]byte("\n")); err != nil {
						log.Fatal(err)
					}

				}
			}
			depth--
		} else {
			break
		}
	}
}

//!-
