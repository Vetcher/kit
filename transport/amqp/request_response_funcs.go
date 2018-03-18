package amqp

import (
	"context"

	"github.com/streadway/amqp"
)

// DeliveryFunc may take information from an amqp delivery and put it into a
// request context. In Subscribers, DeliveryFunc are executed prior to invoking the
// endpoint.
type DeliveryFunc func(context.Context, *amqp.Delivery) context.Context

// PublishingFunc may take information from an amqp publishing request and put it into a
// request context. In Subscribers, PublishingFunc are executed after invoking the
// endpoint.
type PublishingFunc func(context.Context, *amqp.Publishing) context.Context
