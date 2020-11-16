package subscriber

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

// Subscribe uses the existing NATS connection to subscribe to given subject.
// The groupName is used to ensure other subscribers with the same groupName to
// not receive the same element.
//
// Without context handling, this can block forever.
func (c Connection) Subscribe(ctx context.Context, subj, groupName string) ([]byte, error) {
	sub, err := c.conn.QueueSubscribeSync(subj, groupName)
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

// // SubscribeWithChannel uses the existing NATS connection to subscribe to given
// // subject. The groupName is used to ensure other subscribers with the same
// // groupName to not receive the same element.
// func (c Connection) SubscribeWithChannel(subj, groupName string) (<-chan []byte, error) {
// 	msgCh := make(chan *nats.Msg, 1)
// 	resultCh := make(chan []byte, 1)

// 	_, err := c.conn.ChanQueueSubscribe(subj, groupName, msgCh)
// 	if err != nil {
// 		return nil, err
// 	}

// 	go func() {
// 		for {
// 			select {
// 			case m, ok := <-msgCh:
// 				if !ok { // means channel closed
// 					return
// 				}
// 				m.Respond(nil)
// 				resultCh <- m.Data
// 			}
// 		}
// 	}()

// 	return resultCh, nil
// }
