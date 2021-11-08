package main

import (
	"fmt"
	"log"
	"net"

	protos "github.com/truesch/grpc_getting_started/protos/translation"
	"github.com/truesch/grpc_getting_started/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// create new gRPC server
	s := grpc.NewServer()

	// create new instance of Translation server
	trans := server.NewTranslation()

	// register reflection API https://github.com/grpc/grpc/blob/master/doc/server-reflection.md
	reflection.Register(s)

	// register it to the grpc server
	protos.RegisterTranslationServer(s, trans)

	// create socket to listen to requests
	tl, err := net.Listen("tcp", "localhost:8765")
	if err != nil {
		log.Fatal(fmt.Println("Error starting tcp listener on port 8765", err))
	}

	// start listening
	s.Serve(tl)
}
