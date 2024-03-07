package mid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type requestOptions struct {
	withPost   bool
	postValues url.Values
}

func newRequestOptions(opts ...requestOption) *requestOptions {
	opt := &requestOptions{
		postValues: url.Values{},
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

type requestOption func(*requestOptions)

func withPostRequestOption() requestOption {
	return func(opt *requestOptions) {
		opt.withPost = true
	}
}

func addPostValueRequestOption(key, value string) requestOption {
	return func(opt *requestOptions) {
		opt.postValues.Add(key, value)
	}
}

func setPostValueRequestOption(key, value string) requestOption {
	return func(opt *requestOptions) {
		opt.postValues.Set(key, value)
	}
}

func withPostValueRequestOption(values url.Values) requestOption {
	return func(opt *requestOptions) {
		opt.postValues = values
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, opts ...requestOption) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.config.Addr+path, nil)
	if err != nil {
		return nil, err
	}

	opt := newRequestOptions(opts...)

	for _, o := range opts {
		o(opt)
	}

	if opt.withPost {
		opt.postValues.Set("BusinessNo", c.config.BusinessNo)
		opt.postValues.Set("ApiVersion", c.config.ApiVersion)
		opt.postValues.Set("HashKeyNo", c.config.HashKeyNo)
		opt.postValues.Set("VerifyNo", "1")
		opt.postValues.Set("IdentifyNo", "1")
		req.PostForm = opt.postValues
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, dest interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mid: http status code: %d", resp.StatusCode)
	}

	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(resp.Body); err != nil {
		return err
	}

	return json.Unmarshal(buffer.Bytes(), dest)
}
