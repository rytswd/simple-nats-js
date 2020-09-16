package subscriber

import (
	"github.com/nats-io/nats.go"
)

// Connection holds NATS connection.
//
// Use ConnectToConsumer method to initialize connection correctly.
type Connection struct {
	conn *nats.Conn
}

// ConnectToConsumer sets up NATS JetStream connection given the endpoint address,
// Stream name, and Consumer name.
func ConnectToConsumer(addr, streamName, consumerName string) (Connection, error) {
	nc, err := nats.Connect(addr)
	if err != nil {
		return Connection{}, err
	}

	result := Connection{
		conn: nc,
	}

	return result, nil
}

// Close closes connection.
func (c *Connection) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// IsConnected returns the connection status.
func (c *Connection) IsConnected() bool {
	return c.conn.IsConnected()
}
