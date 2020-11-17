package testhelper

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/nats-io/jsm.go"
	server "github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

// JS holds NATS JetStream server setup for testing.
type JS struct {
	// Conn is used for connecting to running NATS JetStream server.
	Conn *nats.Conn

	// Mgr is JetStream manager which allows manipulating Stream and Consumer.
	Mgr *jsm.Manager

	// Srv is the running Server. Usually this is not used for test client.
	Srv *server.Server
}

// Close closes NATS JetStream server.
func (js *JS) Close() {
	js.Conn.Flush()
	js.Srv.Shutdown()
}

// CreateSimpleStream creates a simple Stream based on stream name and
// subjects.
func (js *JS) CreateSimpleStream(t testing.TB, streamName string, streamSubjects []string) *jsm.Stream {
	stream, err := js.Mgr.NewStream(streamName, jsm.FileStorage(), jsm.Subjects(streamSubjects...))
	if err != nil {
		t.Fatalf("creating steram failed: %v", err)
	}
	return stream
}

// CreateSimpleConsumer creates a simple Consumer based on Stream, consumer
// name and its subject.
func (js *JS) CreateSimpleConsumer(t testing.TB, stream *jsm.Stream, consumerName, consumerSubject string) *jsm.Consumer {
	consumer, err := stream.NewConsumer(
		jsm.DurableName(consumerName),
		jsm.FilterStreamBySubject(consumerSubject),
		jsm.DeliverAllAvailable())
	if err != nil {
		t.Fatalf("creating consumer failed: %v", err)
	}
	return consumer
}

// Prepopulate stores data for further processing
func (js *JS) Prepopulate(t testing.TB, subject string, data [][]byte) {
	for _, p := range data {
		// m := nats.NewMsg(tc.subscribeSubject)
		// m.Data = p
		// nc.PublishMsg(m)

		js.Conn.Publish(subject, p)
	}
}

// StartJSWithStreamAndConsumer creates NATS JetStream server, with Stream and
// Consumer created at the same time. If you need to create more than one
// Stream or Consumer, you may want to use CreateSimpleStream and/or
// CreateSimpleConsumer along with this.
//
// Make sure to call Close().
// 	defer js.Close()
//
// This function contains the following setup. You can adjust according to your
// needs.
//
// 	js := testhelper.StartJS(t)
// 	stream := js.CreateSimpleStream(t, tc.streamName, tc.streamSubjects)
// 	js.CreateSimpleConsumer(t, stream, tc.consumerName, tc.subscribeSubject)
// 	js.Prepopulate(t, tc.subscribeSubject, tc.prepopulate)
func StartJSWithStreamAndConsumer(t testing.TB,
	streamName string, streamSubjects []string,
	consumerName string, consumerSubject string,
	prepopulate [][]byte) *JS {
	t.Helper()

	js := StartJS(t)

	// Set up Stream
	stream := js.CreateSimpleStream(t, streamName, streamSubjects)

	// Set up Consumer
	consumer := js.CreateSimpleConsumer(t, stream, consumerName, consumerSubject)
	_ = consumer

	// Prepopulate data
	js.Prepopulate(t, consumerSubject, prepopulate)

	return js
}

// StartJS creates NATS JetStream server for testing.
func StartJS(t testing.TB) *JS {
	t.Helper()

	d, err := ioutil.TempDir("", "jstest")
	if err != nil {
		t.Fatalf("temp dir could not be made: %s", err)
	}

	opts := natsserver.DefaultTestOptions
	opts.JetStream = true
	opts.StoreDir = d
	opts.Port = -1
	opts.NoLog = false
	opts.TraceVerbose = true
	opts.Trace = true
	opts.LogFile = "/tmp/nats.log"

	s, err := server.NewServer(&opts)
	if err != nil {
		t.Fatal("server start failed: ", err)
	}

	s.ConfigureLogger()
	go s.Start()

	if !s.ReadyForConnections(5 * time.Second) {
		t.Error("nats server did not start")
	}

	nc, err := nats.Connect(s.ClientURL())
	if err != nil {
		t.Fatalf("client start failed: %s", err)
	}

	mgr, err := jsm.New(nc, jsm.WithTimeout(time.Second))
	if err != nil {
		t.Fatalf("manager creation failed: %s", err)
	}

	js := &JS{
		Srv:  s,
		Conn: nc,
		Mgr:  mgr,
	}

	return js
}
