package port

import "context"

type MessageBroker interface {
	PublishMessage(ctx context.Context, value []byte) error
	Close() error
}
