package main

import (
	"log"
	"net"
	"os"
	"time"

	protos "github.com/truesch/grpc_getting_started/protos/translation"

	"github.com/getsentry/sentry-go"
	"github.com/truesch/grpc_getting_started/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// set up sentry logging
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://d35315e8d9814fea91340ad9d2785a11@o1062395.ingest.sentry.io/6052816",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("Hello from Translation API!")

	// create new server
	gs := grpc.NewServer()

	// create new instance of Translation server
	ts := server.NewTranslation()

	// register reflection API
	reflection.Register(gs)

	// register it to the grpc server
	protos.RegisterTranslationServer(gs, ts)

	// create socket to listen to requests
	tl, err := net.Listen("tcp", "localhost:9092")
	if err != nil {
		sentry.CaptureException(err)
		os.Exit(1)
	}

	// start listening
	gs.Serve(tl)
}
