package main

import (
	"fmt"
	"log"

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

	ch, err := s.SubscribeWithChannel("*", "some_group_name")
	if err != nil {
		log.Printf("subscription failed, %v\n", err)
	}
	for {
		select {
		case m, ok := <-ch:
			if !ok { // means closed
				return
			}
			log.Printf("data received, '%s'\n", m)
		}
	}
}
