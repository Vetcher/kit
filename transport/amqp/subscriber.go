package amqp

import (
	"context"

	"github.com/streadway/amqp"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Subscriber wraps an endpoint.
type Subscriber struct {
	e            endpoint.Endpoint
	dec          DecodeDeliveryFunc
	enc          EncodePublishingFunc
	before       []DeliveryFunc
	after        []PublishingFunc
	logger       log.Logger
	errorEncoder ErrorEncoder
}

// NewSubscriber constructs a new subscriber, which wraps endpoint and provides amqp handler.
func NewSubscriber(e endpoint.Endpoint, dec DecodeDeliveryFunc, enc EncodePublishingFunc, options ...SubscriberOption) *Subscriber {
	sub := &Subscriber{
		e:            e,
		dec:          dec,
		enc:          enc,
		logger:       log.NewNopLogger(),
		errorEncoder: NopErrorEncoder,
	}
	for _, option := range options {
		option(sub)
	}
	return sub
}

// SubscriberOption sets an optional parameter for subscribers.
type SubscriberOption func(*Subscriber)

// SubscriberErrorLogger is used to log non-terminal errors. By default, no errors
// are logged. This is intended as a diagnostic measure. Finer-grained control
// of error handling, including logging in more detail, should be performed in a
// custom SubscriberErrorEncoder which has access to the context.
func SubscriberErrorLogger(logger log.Logger) SubscriberOption {
	return func(subscriber *Subscriber) { subscriber.logger = logger }
}

func SubscriberBefore(before ...DeliveryFunc) SubscriberOption {
	return func(subscriber *Subscriber) { subscriber.before = append(subscriber.before, before...) }
}

func SubscriberAfter(after ...PublishingFunc) SubscriberOption {
	return func(subscriber *Subscriber) { subscriber.after = append(subscriber.after, after...) }
}

// SubscriberErrorEncoder sets ErrorEncoder for subscriber.
func SubscriberErrorEncoder(encoder ErrorEncoder) SubscriberOption {
	return func(subscriber *Subscriber) { subscriber.errorEncoder = encoder }
}

// ErrorEncoder is a function, that should construct and encode error and push it to broker via PublishFunc.
type ErrorEncoder func(ctx context.Context, err error, publish PublishFunc)

// NopErrorEncoder is a default Subscriber's ErrorEncoder, which does nothing.
func NopErrorEncoder(context.Context, error, PublishFunc) { return }

// ServeAMQP executes subscriber endpoint.
// If ReplyTo field is empty, function does not invokes encode and returns nil publishing and error.
func (s Subscriber) ServeAMQP(delivery amqp.Delivery) (*amqp.Publishing, error) {
	ctx := context.Background()
	for _, f := range s.before {
		f(ctx, &delivery)
	}
	request, err := s.dec(ctx, delivery)
	if err != nil {
		s.logger.Log("err", err)
		return nil, err
	}
	response, err := s.e(ctx, request)
	if err != nil {
		s.logger.Log("err", err)
		return nil, err
	}

	if delivery.ReplyTo == "" {
		return nil, nil
	}

	pub, err := s.enc(ctx, response)
	if err != nil {
		s.logger.Log("err", err)
		return nil, err
	}
	return &pub, nil
}
