package mocks

import (
	"github.com/streadway/amqp"
)

type Channel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}