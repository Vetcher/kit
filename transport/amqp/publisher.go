package amqp

import (
	"context"

	"github.com/streadway/amqp"

	"github.com/go-kit/kit/endpoint"
)

// PublishFunc is a function that publish amqp.Publishing to broker.
// PublishFuncs are used in Publishers and Subscribers.
// Publisher use it to send message to broker.
// Subscriber use it to reply deliveries and push errors.
type PublishFunc func(ctx context.Context, publishing amqp.Publishing) error

type Publisher struct {
	enc     EncodePublishingFunc
	before  []PublishingFunc
	publish PublishFunc
}

type PublisherOption func(*Publisher)

func PublisherBefore(before ...PublishingFunc) PublisherOption {
	return func(publisher *Publisher) { publisher.before = append(publisher.before, before...) }
}

func NewPublisher(publish PublishFunc, enc EncodePublishingFunc, options ...PublisherOption) *Publisher {
	p := &Publisher{
		publish: publish,
		enc:     enc,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p Publisher) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		pub, err := p.enc(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, f := range p.before {
			f(ctx, &pub)
		}
		err = p.publish(ctx, pub)
		if err != nil {
			return nil, err
		}
		// fixme: How should we handle response from subscriber?
		if pub.ReplyTo == "" {
			return nil, nil
		}
		return nil, nil
	}
}
