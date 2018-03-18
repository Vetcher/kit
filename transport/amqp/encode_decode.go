package amqp

import (
	"context"

	"github.com/streadway/amqp"
)

// EncodePublishingFunc constructs publishing from passed object. The result is used
// in subscribers before replying on delivery and in publishers to send messages.
type EncodePublishingFunc func(context.Context, interface{}) (amqp.Publishing, error)

// DecodeDeliveryFunc extracts a user-domain request from amqp Delivery object.
// It is used in subscriber.
type DecodeDeliveryFunc func(context.Context, amqp.Delivery) (interface{}, error)
