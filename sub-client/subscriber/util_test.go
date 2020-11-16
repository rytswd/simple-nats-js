package subscriber

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/nats-io/jsm.go"
	server "github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

// StartJetStreamServer starts NATS JetStream server for testing.
//
// Taken from https://github.com/nats-io/jsm.go/blob/ecfafe3d16278cffd637f604ee25f28e3dbcc8d5/manager_test.go#L15
//
// You can follow the below pattern:
//
// 	srv, nc, mgr := startJSServer(t)
// 	defer srv.Shutdown()
// 	defer nc.Flush()
//
// This returns all toolings necessary for interacting with JetStream server.
func StartJetStreamServer(t *testing.T) (*server.Server, *nats.Conn, *jsm.Manager) {
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

	return s, nc, mgr
}
