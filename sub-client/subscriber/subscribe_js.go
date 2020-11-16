package subscriber

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

// SubscribeToJetStream subscribes to the given subject by creating Consumer
// to Stream.
func (c *Connection) SubscribeToJetStream(ctx context.Context, subj, groupName, streamName, consumerName string) ([]byte, error) {
	consumer := nats.Consumer(streamName, nats.ConsumerConfig{
		Durable:       consumerName,
		DeliverPolicy: nats.DeliverAll,
		AckPolicy:     nats.AckExplicit,
		AckWait:       5 * time.Second,
		ReplayPolicy:  nats.ReplayInstant,
	})

	// Subject needs to be empty string?
	sub, err := c.conn.QueueSubscribeSync("", groupName, consumer)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe() // Not sure if you need this

	msg, err := sub.NextMsgWithContext(ctx)
	if err != nil {
		return nil, err
	}
	msg.Respond(nil)

	return msg.Data, nil
}
