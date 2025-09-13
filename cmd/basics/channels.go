package main

import (
	"fmt"
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

	// creating a channel
	c := make(chan string)

	// launch concurrent link checkers
	for _, link := range links {
		// concurrency is achieved by using go keyword.
		go checkLink(link, c)
	}

	// Alternative 1: using range to receive the message from the channel
	for l := range c {
		time.Sleep(5 * time.Second)
		go checkLink(l, c)
	}

	// Alternative 2: using a for loop to receive the message from the channel
	// for {
	// 	time.Sleep(5 * time.Second)
	// 	go checkLink(<-c, c)
	// }

	// Alternative 3: anonymous function to create a new go routine
	// for l := range c {
	// 	go func(link string) {
	// 		time.Sleep(5 * time.Second)
	// 		go checkLink(link, c)
	// 	}(l)
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
