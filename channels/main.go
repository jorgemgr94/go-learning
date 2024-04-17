package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
		// "https://lavelada.es",
	}

	// creating a channel
	c := make(chan string)

	for _, link := range links {
		// concurrency is achieved by using go keyword.
		go checkLink(link, c)
	}

	// this for loop purpose is to receive the message from the channel when the function is done.
	// meaning that the function is done and the message is sent to the channel.
	// for {
	// 	go checkLink(<-c, c)
	// }

	// alternative way to write the above for loop
	// for l := range c {
	// 	go checkLink(l, c)
	// }

	for l := range c {
		// anonymous function: function literal to create a new go routine
		go func(link string) {
			time.Sleep(5 * time.Second)
			go checkLink(link, c)
		}(l)
	}
}

func checkLink(link string, c chan string) {
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
