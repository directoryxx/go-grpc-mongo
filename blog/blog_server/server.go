package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"go-grpc-mongo/blog/blogproto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct{}

func (s *server) CreateBlog(ctx context.Context, req *blogproto.CreateBlogRequest) (*blogproto.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := &blogItem{
		AuthorID: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprint("Error : %v", err))
	}

	objId, errMongo := res.InsertedID.(primitive.ObjectID)

	if !errMongo {
		return nil, status.Errorf(codes.Internal, fmt.Sprint("Error : %v", errMongo))
	}

	return &blogproto.CreateBlogResponse{
		Blog: &blogproto.Blog{
			Id:       objId.Hex(),
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil
}

func (s *server) ReadBlog(ctx context.Context, req *blogproto.ReadBlogRequest) (*blogproto.ReadBlogResponse, error) {
	blogId := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprint("Error : %v", err))
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}

	find := collection.FindOne(context.Background(), filter)
	if err := find.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &blogproto.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *blogproto.UpdateBlogRequest) (*blogproto.UpdateBlogResponse, error) {
	blog := req.GetBlog()

	oid, err := primitive.ObjectIDFromHex(blog.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}

	find := collection.FindOne(context.Background(), filter)
	if err := find.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	//Update data
	data.AuthorID = blog.AuthorId
	data.Title = blog.Title
	data.Content = blog.Content

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)

	if updateErr != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot update blog with specified ID: %v", updateErr),
		)
	}

	return &blogproto.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func dataToBlogPb(data *blogItem) *blogproto.Blog {
	return &blogproto.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
		fmt.Println("Disconnecting MongoDB")
	}()

	collection = client.Database("mydb").Collection("blog")

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
