package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	hystrix "github.com/afex/hystrix-go/hystrix"
	proto "github.com/lreimer/advanced-cloud-native-go/Communication/Go-Micro/proto"
	micro "github.com/micro/go-micro"
	breaker "github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"golang.org/x/net/context"
)

// The Greeter API.
type Greeter struct{}

// Hello is a Greeter API method.
func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func callEvery(d time.Duration, greeter proto.GreeterClient, f func(time.Time, proto.GreeterClient)) {
	for x := range time.Tick(d) {
		f(x, greeter)
	}
}

func hello(t time.Time, greeter proto.GreeterClient) {
	// Call the greeter
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "Leander, calling at " + t.String()})
	if err != nil {
		if err.Error() == "hystrix: timeout" {
			fmt.Printf("%s. Insert fallback logic here.\n", err.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	// Print response
	fmt.Printf("%s\n", rsp.Greeting)
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),
	)

	// Init will parse the command line flags. Any flags set will
	// override the above settings.
	// specify a Hystrix breaker client wrapper here
	service.Init(
		micro.WrapClient(breaker.NewClientWrapper()),
	)

	// override some default values for the Hystrix breaker
	hystrix.DefaultVolumeThreshold = 3
	hystrix.DefaultErrorPercentThreshold = 75
	hystrix.DefaultTimeout = 500
	hystrix.DefaultSleepWindow = 3500

	// export Hystrix stream
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8081"), hystrixStreamHandler)

	// Create new greeter client and call hello
	greeter := proto.NewGreeterClient("greeter", service.Client())
	callEvery(3*time.Second, greeter, hello)
}
