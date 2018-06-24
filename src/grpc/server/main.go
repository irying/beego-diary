package main

import (
	"flag"
	"net"
	"grpc/etcd"
	"time"
	"os"
	"os/signal"
	"syscall"
	"log"
	"google.golang.org/grpc"
	"context"
	"fmt"
	"grpc/helloworld"
)

var (
	service = flag.String("service", "hello_service", "service name")
	host = flag.String("host", "localhost", "listening host")
	port = flag.String("port", "50001", "listening port")
	registerAddress = flag.String("reg", "http://localhost:2379", "register etcd address")
)

type server struct{
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", net.JoinHostPort(*host, *port))

	err = etcd.Register(*service, *host, *port, *registerAddress, time.Second*10, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		etcd.UnRegister()
		os.Exit(1)
	}()

	log.Printf("starting hello service at %s", *port)
	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	s.Serve(lis)

}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest)(*helloworld.HelloReply, error)  {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)

	return &helloworld.HelloReply{Message: "Hello " + in.Name + " from " + net.JoinHostPort(*host, *port)}, nil
}