package main

import (
	"flag"
	"grpc/etcd"
	"google.golang.org/grpc"
	"context"
	"time"
	"strconv"
	"fmt"
	"grpc/helloworld"
)

var (
	service = flag.String("service", "hello_service", "service name")
	registerAddress = flag.String("reg", "http://localhost:2379", "register etcd address")
)

func main()  {
	flag.Parse()
	r := etcd.NewResolver(*service)
	b := grpc.RoundRobin(r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *registerAddress, grpc.WithInsecure(), grpc.WithBalancer(), grpc.WithBlock())
	cancel()

	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	for t := range ticker.C {
		client := helloworld.NewGreeterClient(conn)
		resp, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		}
	}
}
