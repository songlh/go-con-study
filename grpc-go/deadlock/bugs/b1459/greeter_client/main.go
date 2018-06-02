/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/golang/glog"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/b1459/helloworld"
	"google.golang.org/grpc/keepalive"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    time.Second,
			Timeout: 2 * time.Second,
		}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	glog.Info("sleeping for 2 keepalives")
	time.Sleep(2400 * time.Millisecond)

	doit := func() {
		s, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("could not open stream: %v", err)
		}
		r, err := s.Recv()
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		glog.Infof("Greeting: %s", r.Message)
		err = s.CloseSend()
		glog.Info("stream closed")
		if err != nil {
			log.Fatalf("failed to close stream: %v", err)
		}
	}

	doit()

	glog.Info("sleeping for 3-4 keepalives")
	time.Sleep(6000 * time.Millisecond)

	doit()

	glog.Info("sleeping for 3-4 keepalives")
	time.Sleep(6000 * time.Millisecond)

	panic("fuck")
}
