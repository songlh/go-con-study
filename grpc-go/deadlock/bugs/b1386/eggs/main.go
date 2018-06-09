package main

//go:generate protoc --go_out=plugins=grpc:. eggs.proto

import (
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/tap"
	pb "bugs/eggs/eggs"
)

type service struct{}

func (s *service) Echo(ctx context.Context, in *wrappers.StringValue) (*wrappers.StringValue, error) {
	log.Print("Server sleeping, will timeout")
	time.Sleep(5 * time.Second)
	return &wrappers.StringValue{ Value: "Other"}, nil
}

func intap(ctx context.Context, info *tap.Info) (context.Context, error) {
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)
	return ctx, nil
}

func main() {
	const target = "127.0.0.1:1234"

	go func() {
		/*
			creds, err := credentials.NewServerTLSFromFile("plinthd.pem", "plinthd.key")
			if err != nil {
				panic(err)
			}
			grpcServer := grpc.NewServer(grpc.InTapHandle(intap), grpc.Creds(creds))
		*/
		grpcServer := grpc.NewServer(grpc.InTapHandle(intap))
		pb.RegisterEggsServer(grpcServer, &service{})
		lis, err := net.Listen("tcp", target)
		if err != nil {
			panic(err)
		}
		grpcServer.Serve(lis)
	}()

	time.Sleep(3 * time.Second) // Wait for startup
	/*
		creds, err := credentials.NewClientTLSFromFile("ca.pem", "localhost")
		if err != nil {
			log.Fatal(err)
		}
		conn, err := grpc.DialContext(context.Background(), target, grpc.WithTransportCredentials(creds))
	*/
	conn, err := grpc.DialContext(context.Background(), target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial target %s: %v", target, err)
	}
	defer conn.Close()
	client := pb.NewEggsClient(conn)
	log.Print("Calling")
	res, err := client.Echo(context.Background(), &wrappers.StringValue{ Value: "Hello"})
	if err != nil {
		log.Fatalf("client call failed: %v", err)
	}
	log.Printf("%+v", res)
}