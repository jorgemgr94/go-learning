package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

func channels() {
	fmt.Println("// == Channels ===========================================")
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	c := make(chan string)

	// Random timeout between 10-30 seconds
	randomTimeout := time.Duration(rand.IntN(20)+10) * time.Second
	fmt.Printf("Will exit after %v\n", randomTimeout)

	timer := time.NewTimer(randomTimeout)

	for _, link := range links {
		go checkLink(link, c)
	}

	for {
		select {
		case l := <-c:
			time.Sleep(5 * time.Second)
			go checkLink(l, c)
		case <-timer.C:
			fmt.Println("Randomly exiting due to timeout!")
			close(c)
			return
		}
	}

	// Alternative 1: using range to receive the message from the channel (single channel scenario)
	// for l := range c {
	// 	time.Sleep(5 * time.Second)
	// 	go checkLink(l, c)
	// }
}

func checkLink(link string, c chan string) {
	fmt.Println("Checking link:", link, "at", time.Now())
	// check if the link is up
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down!")
		c <- link
		return
	}

	fmt.Println(link, "is up!")
	// sending the link to the channel
	c <- link
}
