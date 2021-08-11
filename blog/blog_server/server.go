package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"go-grpc-mongo/blog/blogproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8010")
	if err != nil {
		panic(err)
	}

	certFile := "certs/directoryxx.com/cert.pem"
	keyFile := "certs/directoryxx.com/privkey.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		fmt.Println("Failed to load : ", sslErr)
		return
	}
	opts := grpc.Creds(creds)
	s := grpc.NewServer(opts)
	blogproto.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Blog service started...")
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("Stopping server")
	s.Stop()
	lis.Close()
	fmt.Println("Successfully stop server")

}
