package main

import (
	"context"
	"fmt"
	"go-grpc-mongo/blog/blogproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	certFile := "certs/directoryxx.com/cert.pem"
	creds, certErr := credentials.NewClientTLSFromFile(certFile, "")
	if certErr != nil {
		fmt.Println("Failed to load : ", certErr)
		return
	}
	opts := grpc.WithTransportCredentials(creds)

	conn, err := grpc.Dial("grpc.directoryxx.com:8010", opts)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	c := blogproto.NewBlogServiceClient(conn)

	// doCreate(c)

	doRead(c)
}

func doCreate(c blogproto.BlogServiceClient) {
	createBlogReq := &blogproto.CreateBlogRequest{
		Blog: &blogproto.Blog{
			AuthorId: "2",
			Title:    "Aku Nulis",
			Content:  "Test tes",
		},
	}
	resCreate, errCreate := c.CreateBlog(context.Background(), createBlogReq)
	if errCreate != nil {
		statErr, ok := status.FromError(errCreate)
		if ok {
			if statErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout")
			} else {
				fmt.Println("Unexpected Error", statErr)
			}
		}
	}

	fmt.Println(resCreate)
}

func doRead(c blogproto.BlogServiceClient) {
	readBlogReq := &blogproto.ReadBlogRequest{
		BlogId: "6113a6fb1ced4cdce6c4bccd",
	}
	resRead, errRead := c.ReadBlog(context.Background(), readBlogReq)
	if errRead != nil {
		statErr, ok := status.FromError(errRead)
		if ok {
			if statErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout")
			} else {
				fmt.Println("Unexpected Error", statErr)
			}
		}
	}

	fmt.Println(resRead)
}
