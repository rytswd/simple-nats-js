package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/rytswd/simple-nats-js/pub-client/publisher"
)

func main() {
	conn, err := publisher.ConnectToStream("localhost:4222", "AnotherStream")
	if err != nil {
		log.Fatalf("error occurred, %v\n", err)
		return
	}

	for {
		// data := fmt.Sprintf("some random data --- %d", rand.Intn(20000))
		data := fmt.Sprintf("some random data --- %s", time.Now().Format(time.RFC3339))
		log.Printf("publishing data '%s'\n", data)

		err = pub(context.Background(), conn, []byte(data))
		if err != nil {
			log.Printf("error occurred, %v\n", err)
			return
		}
		log.Println("done")
		time.Sleep(5 * time.Second)
	}
}

// This function shows how the package should be consumed.
func pub(ctx context.Context, conn *publisher.Connection, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	// Hash the message data for message id
	h := sha1.New()
	_, err := h.Write(data)
	if err != nil {
		return err
	}
	id := fmt.Sprintf("%x", h.Sum(nil))

	err = conn.Publish(ctx, "xyz.test", id, []byte(data))
	return err
}
