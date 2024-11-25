package mid

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
)

type Client struct {
	config     *Config
	httpClient *http.Client
}

type Option func(c *Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
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
}

func (c *Client) buildIdentifyNo(input string) string {
	input += c.config.HashKey

	utf16le := make([]byte, len(input)*2) //nolint:mnd
	for i, r := range input {
		binary.LittleEndian.PutUint16(utf16le[i*2:], uint16(r))
	}

	hash := sha256.Sum256(utf16le)
	return hex.EncodeToString(hash[:])
}
