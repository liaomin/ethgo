package jsonrpc

import (
	"sync"

	"github.com/umbracle/ethgo/jsonrpc/transport"
)

// Client is the jsonrpc client
type Client struct {
	transport transport.Transport
	endpoints endpoints
}

type endpoints struct {
	w *Web3
	e *Eth
	n *Net
	d *Debug
}

type Config struct {
	headers map[string]string
}

type ConfigOption func(*Config)

func WithHeaders(headers map[string]string) ConfigOption {
	return func(c *Config) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

type ClientPool struct {
	addr string
	opts []ConfigOption
	pool sync.Pool
}

func NewClientPool(addr string, opts ...ConfigOption) *ClientPool {
	pool := &ClientPool{addr: addr, opts: opts, pool: sync.Pool{
		New: func() interface{} {
			for true {
				client, _ := NewClient(addr, opts...)
				if client != nil {
					return client
				}
			}
			return nil
		},
	}}
	return pool
}

func (cp *ClientPool) Require() *Client {
	return cp.pool.Get().(*Client)
}

func (cp *ClientPool) Release(c *Client) {
	cp.pool.Put(c)
}

func NewClient(addr string, opts ...ConfigOption) (*Client, error) {
	config := &Config{headers: map[string]string{}}
	for _, opt := range opts {
		opt(config)
	}

	c := &Client{}
	c.endpoints.w = &Web3{c}
	c.endpoints.e = &Eth{c: c}
	c.endpoints.n = &Net{c}
	c.endpoints.d = &Debug{c}

	t, err := transport.NewTransport(addr, config.headers)
	if err != nil {
		return nil, err
	}
	c.transport = t
	return c, nil
}

// Close closes the transport
func (c *Client) Close() error {
	return c.transport.Close()
}

// Call makes a jsonrpc call
func (c *Client) Call(method string, out interface{}, params ...interface{}) error {
	return c.transport.Call(method, out, params...)
}

// SetMaxConnsLimit sets the maximum number of connections that can be established with a host
func (c *Client) SetMaxConnsLimit(count int) {
	c.transport.SetMaxConnsPerHost(count)
}
