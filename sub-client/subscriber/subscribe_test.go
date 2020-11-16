package subscriber

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/nats-io/jsm.go"
)

func TestSubscribe(t *testing.T) {
	cases := map[string]struct {
		// input
		ctx              context.Context // optional
		streamName       string
		streamSubjects   []string
		prepopulate      [][]byte
		consumerName     string
		subscribeSubject string
		subscribeGroup   string

		// output
		want          []byte
		wantErr       error
		wantErrString string
	}{
		"simple test": {
			streamName: "testStream",
			streamSubjects: []string{
				"some.*",
				"dummy.subject",
			},
			prepopulate: [][]byte{
				[]byte("some random data 1"),
				[]byte("some random data 2"),
				[]byte("some random data 3"),
				[]byte("some random data 4"),
			},
			consumerName:     "testConsumer",
			subscribeSubject: "some.subject",
			subscribeGroup:   "someID",
			want:             []byte("abc"), // does not match, force error
		},
		// "fail: context timed out": {
		// 	ctx:            testctx.CancelledContext(),
		// 	streamName:     "testStream",
		// 	streamSubjects: []string{"some.subject"},
		// 	prepopulate: [][]byte{
		// 		[]byte("some data 1"),
		// 		[]byte("some data 2"),
		// 	},
		// 	subscribeSubject: "some.subject",
		// 	subscribeGroup:   "someID",
		// 	wantErr:          context.Canceled,
		// },
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// ----------------------------------------------------------------
			// Prepare NATS JetStream
			srv, nc, mgr := StartJetStreamServer(t)
			defer srv.Shutdown()
			defer nc.Flush()

			// Create Stream
			stream, err := mgr.NewStream(tc.streamName, jsm.FileStorage(), jsm.Subjects(tc.streamSubjects...))
			if err != nil {
				t.Fatalf("creating steram failed: %v", err)
			}

			// Create Consumer
			consumer, err := stream.NewConsumer(
				jsm.DurableName("some_consumer_name"),
				jsm.FilterStreamBySubject(tc.subscribeSubject),
				jsm.DeliverAllAvailable())
			if err != nil {
				t.Fatalf("creating consumer failed: %v", err)
			}
			_ = consumer

			// Prepopulate data
			for _, p := range tc.prepopulate {
				nc.Publish(tc.subscribeSubject, p)
			}

			// Test consumer subscirbe
			ms, _ := consumer.NextMsg()
			fmt.Printf("%#v\n", ms)
			fmt.Printf("Data: %s\n", ms.Data)

			// ----------------------------------------------------------------
			// Main test
			s := Connection{
				conn: nc,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if tc.ctx != nil {
				ctx = tc.ctx
			}

			result, err := s.Subscribe(ctx, tc.subscribeSubject, tc.subscribeGroup)
			if err != nil {
				t.Errorf("%v", err)
			}
			if !reflect.DeepEqual(tc.want, result) {
				t.Errorf("mismatch\n    want: %s\n    got:  %s\n", tc.want, result)
			}
		})
	}
}
