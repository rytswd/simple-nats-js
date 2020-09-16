package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rytswd/simple-nats-js/sub-client/subscriber"
)

func main() {
	streamName := fmt.Sprintf("AnotherStream")
	consumerName := fmt.Sprintf("SomeConsumer")

	log.Printf("connecting to Stream:   %s\n", streamName)
	log.Printf("connecting to Consumer: %s\n", consumerName)

	s, err := subscriber.ConnectToConsumer("localhost:4222", streamName, consumerName)
	if err != nil {
		log.Fatalf("error occurred, %v\n", err)
		return
	}
	defer s.Close()

	for {
		if !s.IsConnected() {
			log.Fatalf("connection error")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
		data, err := s.Subscribe(ctx, "*", "some_group_name")
		if err != nil {
			log.Printf("error occurred, %v\n", err)
			cancel()
			continue
		}
		log.Printf("data received, '%s'\n", data)
		cancel()
	}
}
