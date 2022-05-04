package jsonrpc

import (
	"fmt"

	"github.com/umbracle/ethgo/jsonrpc/transport"
)

// SubscriptionEnabled returns true if the subscription endpoints are enabled
func (c *Client) SubscriptionEnabled() bool {
	r := c.transportPool.Require()
	_, ok := r.(transport.PubSubTransport)
	c.transportPool.Release(r)
	return ok
}

// Subscribe starts a new subscription
func (c *Client) Subscribe(method string, parmas interface{}, callback func(b []byte)) (func() error, error) {
	r := c.transportPool.Require()
	defer c.transportPool.Release(r)

	pub, ok := r.(transport.PubSubTransport)
	if !ok {
		return nil, fmt.Errorf("Transport does not support the subscribe method")
	}
	close, err := pub.Subscribe(method, parmas, callback)
	return close, err
}
