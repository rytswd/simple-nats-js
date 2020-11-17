package subscriber

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

// Subscribe uses the existing NATS connection to subscribe to given Stream and
// Consumer. If Consumer does not exist with the given name, it creates a new
// one with the additionally provided subject.
//
// Without context handling, this can block forever.
func (c *Connection) Subscribe(ctx context.Context, streamName, consumerName, subj string) ([]byte, error) {
	mgr, err := jsm.New(c.conn)
	if err != nil {
		return nil, err
	}

	stream, err := mgr.LoadStream(streamName)
	if err != nil {
		return nil, err
	}

	consumer, err := stream.LoadOrNewConsumer(consumerName,
		jsm.DurableName(consumerName),
		jsm.FilterStreamBySubject(subj),
		jsm.DeliverAllAvailable(),
	)
	if err != nil {
		return nil, err
	}
	_ = consumer

	jsSubject := fmt.Sprintf(`$JS.API.CONSUMER.MSG.NEXT.%s.%s`, streamName, consumerName)
	_ = jsSubject

	// msg, err := mgr.NextMsg(streamName, consumerName)
	// msg, err := consumer.NextMsgContext(ctx)
	msg, err := c.conn.RequestWithContext(ctx, jsSubject, []byte("a"))
	if errors.Is(err, nats.ErrTimeout) {
		return nil, nil
	}
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
