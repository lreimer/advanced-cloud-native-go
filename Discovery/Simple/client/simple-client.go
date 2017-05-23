package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

var url string

func main() {
	lookupServiceWithConsul()

	fmt.Println("Starting Simple Client.")
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	callHelloEvery(5*time.Second, client)
}

func lookupServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, error := consulapi.NewClient(config)
	if error != nil {
		fmt.Println(error)
	}

	services, error := consul.Agent().Services()
	if error != nil {
		fmt.Println(error)
	}

	service := services["simple-server"]
	address := service.Address
	port := service.Port

	url = fmt.Sprintf("http://%s:%v/info", address, port)
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
