package mid

import (
	"net/http"
)

type Client struct {
	config            *Config
	httpClient        *http.Client
	verifyNoGenerator VerifyNoGenerator
}

type Option func(c *Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithVerifyNoGenerator(generator VerifyNoGenerator) Option {
	return func(c *Client) {
		c.verifyNoGenerator = generator
	}
}

func New(config *Config, opts ...Option) *Client {
	config.init()

	c := &Client{
		config: config,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.init()

	return c
}

func (c *Client) init() {
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	if c.verifyNoGenerator == nil {
		c.verifyNoGenerator = DefaultVerifyNoGenerator
	}
}
