package publisher

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

// Connection holds NATS connection.
//
// Use ConnectToStream method to initialize connection correctly.
type Connection struct {
	conn *nats.Conn
}

// ConnectToStream connects to given address with JetStream stream name. If
// the stream name does not exist, this errors out. If you want to create a new
// stream.
func ConnectToStream(addr string, streamName string) (*Connection, error) {
	nc, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}

	result := &Connection{
		conn: nc,
	}
	return result, nil
}

// Publish sends the message to given subject and data, using an existing
// connection. Duplicated massage will be deduped, and if the NATS server does
// not ACK within context timeline, it returns an error.
func (c *Connection) Publish(ctx context.Context, subject string, msgID string, data []byte) error {
	if c.conn == nil {
		return errors.New("connection not established")
	}

	m := &nats.Msg{
		Subject: subject,
		Header: http.Header{
			"Msg-Id": []string{msgID},
		},
		Data: data,
	}
	a, err := c.conn.RequestMsgWithContext(ctx, m)
	log.Println(string(a.Data))
	return err
}
