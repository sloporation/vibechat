package transactionlog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"vibechat/internal/config"
)

const (
	streamName = "VIBECHAT"
	subject    = "vibechat.events"
)

// Event is the unit of the transaction log.
type Event struct {
	ID        string          `json:"id"`
	Action    string          `json:"action"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

// TransactionLog wraps a NATS JetStream connection and provides
// Commit and Replay for the server event log.
type TransactionLog struct {
	js     jetstream.JetStream
	stream jetstream.Stream
}

func New(cfg config.NATSConfig) (*TransactionLog, error) {
	opts := []nats.Option{}

	switch {
	case cfg.Token != "":
		opts = append(opts, nats.Token(cfg.Token))
	case cfg.User != "" && cfg.Password != "":
		opts = append(opts, nats.UserInfo(cfg.User, cfg.Password))
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("jetstream init: %w", err)
	}

	ctx := context.Background()
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  jetstream.FileStorage,
	})
	if err != nil {
		return nil, fmt.Errorf("create stream: %w", err)
	}

	return &TransactionLog{js: js, stream: stream}, nil
}

// Commit publishes an event to the transaction log.
func (tl *TransactionLog) Commit(ctx context.Context, e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	if _, err := tl.js.Publish(ctx, subject, data); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}

// Replay calls fn for each event in the log starting from the given
// sequence number. Pass 1 to replay from the beginning.
func (tl *TransactionLog) Replay(ctx context.Context, fromSeq uint64, fn func(Event) error) error {
	consumer, err := tl.stream.OrderedConsumer(ctx, jetstream.OrderedConsumerConfig{
		DeliverPolicy: jetstream.DeliverByStartSequencePolicy,
		OptStartSeq:   fromSeq,
	})
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	info, err := tl.stream.Info(ctx)
	if err != nil {
		return fmt.Errorf("stream info: %w", err)
	}
	lastSeq := info.State.LastSeq

	msgs, err := consumer.Messages()
	if err != nil {
		return fmt.Errorf("consumer messages: %w", err)
	}
	defer msgs.Stop()

	for {
		msg, err := msgs.Next()
		if err != nil {
			return fmt.Errorf("next message: %w", err)
		}

		var e Event
		if err := json.Unmarshal(msg.Data(), &e); err != nil {
			msg.Nak()
			return fmt.Errorf("unmarshal event: %w", err)
		}

		if err := fn(e); err != nil {
			msg.Nak()
			return err
		}
		msg.Ack()

		meta, err := msg.Metadata()
		if err == nil && meta.Sequence.Stream >= lastSeq {
			break
		}
	}

	return nil
}
