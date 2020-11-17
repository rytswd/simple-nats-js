package subscriber

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/rytswd/simple-nats-js/sub-client/testhelper"
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
		subscribeCount   int

		// output
		want    [][]byte
		wantErr error
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
			subscribeCount:   2,

			want: [][]byte{
				[]byte("some random data 1"),
				[]byte("some random data 2"),
			},
		},
		"fail: context timed out": {
			ctx:            cancelledContext(),
			streamName:     "testStream",
			streamSubjects: []string{"some.subject"},
			prepopulate: [][]byte{
				[]byte("some data 1"),
				[]byte("some data 2"),
			},
			consumerName:     "testConsumer",
			subscribeSubject: "some.subject",
			subscribeGroup:   "someID",
			subscribeCount:   1,
			wantErr:          context.Canceled,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Prepare NATS JetStream
			js := testhelper.StartJSWithStreamAndConsumer(t, tc.streamName, tc.streamSubjects, tc.consumerName, tc.subscribeSubject, tc.prepopulate)
			// js := testhelper.StartJS(t)
			// stream := js.CreateSimpleStream(t, tc.streamName, tc.streamSubjects)
			// js.CreateSimpleConsumer(t, stream, tc.consumerName, tc.subscribeSubject)
			// js.Prepopulate(t, tc.subscribeSubject, tc.prepopulate)
			defer js.Close()

			// Main test
			s := Connection{
				conn: js.Conn,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if tc.ctx != nil {
				ctx = tc.ctx
			}

			for i := 0; i < tc.subscribeCount; i++ {
				result, err := s.Subscribe(ctx, tc.streamName, tc.consumerName, tc.subscribeSubject)
				if err != nil {
					if !errors.Is(err, tc.wantErr) {
						t.Fatalf("error mismatch\n    want: %v\n    got:  %v\n", err, tc.wantErr)
					}
					return
				}
				if !reflect.DeepEqual(tc.want[i], result) {
					t.Errorf("mismatch\n    want: %s\n    got:  %s\n", tc.want, result)
				}
			}
		})
	}
}

func TestSubscribeWithMultipleSubscribers(t *testing.T) {
	cases := map[string]struct {
		// input
		ctx1             context.Context // optional
		ctx2             context.Context // optional
		streamName       string
		streamSubjects   []string
		prepopulate      [][]byte
		consumerName     string
		subscribeSubject string
		subscribeGroup   string
		subscribeCount   int

		// output
		want1    [][]byte
		wantErr1 error
		want2    [][]byte
		wantErr2 error
	}{
		"multiple subscriptions using same consumer": {
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
			subscribeCount:   2,

			want1: [][]byte{
				[]byte("some random data 1"),
				[]byte("some random data 3"),
			},
			want2: [][]byte{
				[]byte("some random data 2"),
				[]byte("some random data 4"),
			},
		},
		"fail: context timed out": {
			ctx1:           cancelledContext(),
			ctx2:           cancelledContext(),
			streamName:     "testStream",
			streamSubjects: []string{"some.subject"},
			prepopulate: [][]byte{
				[]byte("some data 1"),
				[]byte("some data 2"),
			},
			consumerName:     "testConsumer",
			subscribeSubject: "some.subject",
			subscribeGroup:   "someID",
			subscribeCount:   1,
			wantErr1:         context.Canceled,
			wantErr2:         context.Canceled,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Prepare NATS JetStream
			js := testhelper.StartJSWithStreamAndConsumer(t, tc.streamName, tc.streamSubjects, tc.consumerName, tc.subscribeSubject, tc.prepopulate)
			// js := testhelper.StartJS(t)
			// stream := js.CreateSimpleStream(t, tc.streamName, tc.streamSubjects)
			// js.CreateSimpleConsumer(t, stream, tc.consumerName, tc.subscribeSubject)
			// js.Prepopulate(t, tc.subscribeSubject, tc.prepopulate)
			defer js.Close()

			// Main test
			s := Connection{
				conn: js.Conn,
			}

			ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if tc.ctx1 != nil {
				ctx1 = tc.ctx1
			}
			ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if tc.ctx2 != nil {
				ctx2 = tc.ctx2
			}

			for i := 0; i < tc.subscribeCount; i++ {
				// Subscriber 1
				result1, err := s.Subscribe(ctx1, tc.streamName, tc.consumerName, tc.subscribeSubject)
				if err != nil {
					if !errors.Is(err, tc.wantErr1) {
						t.Fatalf("error mismatch\n    want: %v\n    got:  %v\n", err, tc.wantErr1)
					}
					return
				}
				if !reflect.DeepEqual(tc.want1[i], result1) {
					t.Errorf("mismatch\n    want: %s\n    got:  %s\n", tc.want1, result1)
				}

				// Subscriber 2
				result2, err := s.Subscribe(ctx2, tc.streamName, tc.consumerName, tc.subscribeSubject)
				if err != nil {
					if !errors.Is(err, tc.wantErr2) {
						t.Fatalf("error mismatch\n    want: %v\n    got:  %v\n", err, tc.wantErr2)
					}
					return
				}
				if !reflect.DeepEqual(tc.want2[i], result2) {
					t.Errorf("mismatch\n    want: %s\n    got:  %s\n", tc.want2, result2)
				}
			}
		})
	}
}

// CancelledContext returns context that has passed a deadline
func cancelledContext() context.Context {
	c, cancel := context.WithCancel(context.Background())
	cancel()
	return c
}
