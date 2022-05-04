package jsonrpc

import (
	"fmt"
	"sync"

	"github.com/umbracle/ethgo/jsonrpc/transport"
)

// Client is the jsonrpc client
type Client struct {
	// transport transport.Transport
	transportPool *TransportPool
	endpoints     endpoints
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

type TransportPool struct {
	url     string
	headers map[string]string
	pool    sync.Pool
}

func NewTransportPool(url string, headers map[string]string) *TransportPool {
	pool := &TransportPool{url: url, headers: headers, pool: sync.Pool{
		New: func() interface{} {
			for true {
				transport, _ := transport.NewTransport(url, headers)
				if transport != nil {
					return transport
				}
			}
			panic(fmt.Errorf("can‘t get transport:", url))
		},
	}}
	return pool
}

func (cp *TransportPool) Require() transport.Transport {
	return cp.pool.Get().(transport.Transport)
}

func (cp *TransportPool) Release(c transport.Transport) {
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

	// t, err := transport.NewTransport(addr, config.headers)
	// if err != nil {
	// 	return nil, err
	// }
	// c.transport = t
	c.transportPool = NewTransportPool(addr, config.headers)
	return c, nil
}

// Close closes the transport
func (c *Client) Close() error {
	// return c.transport.Close()
	return nil
}

// Call makes a jsonrpc call
func (c *Client) Call(method string, out interface{}, params ...interface{}) error {
	pool := c.transportPool
	transport := pool.Require()
	defer pool.Release(transport)
	return transport.Call(method, out, params...)
}

// SetMaxConnsLimit sets the maximum number of connections that can be established with a host
func (c *Client) SetMaxConnsLimit(count int) {
	// c.transport.SetMaxConnsPerHost(count)
}
