package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var url string

func main() {
	initServiceURL()

	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	callHelloEvery(5*time.Second, client)
}

func initServiceURL() {
	url = os.Getenv("SERVICE_URL")
	if len(url) == 0 {
		url = "http://simple-k8s-server:8080/info"
	}
}

func hello(t time.Time, client *http.Client) {
	// Call the greeter
	response, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	// print response
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s. Time is %v\n", body, t)
}

func callHelloEvery(d time.Duration, client *http.Client) {
	for x := range time.Tick(d) {
		hello(x, client)
	}
}
