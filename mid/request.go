package mid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type requestOptions struct {
	postValues       url.Values
	withSystemParams bool

	withIdentifyNo   bool
	identifyNoParams []string
}

func newRequestOptions(opts ...RequestOption) *requestOptions {
	opt := &requestOptions{
		postValues: url.Values{},
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

type RequestOption func(*requestOptions)

func setPostValueRequestOption(key, value string) RequestOption {
	return func(opt *requestOptions) {
		opt.postValues.Set(key, value)
	}
}

func withSystemParamsRequestOption() RequestOption {
	return func(opt *requestOptions) {
		opt.withSystemParams = true
	}
}

func withPostValueRequestOption(values url.Values) RequestOption { //nolint:unused
	return func(opt *requestOptions) {
		opt.postValues = values
	}
}

func withIdentifyNoRequestOption(params []string) RequestOption {
	return func(opt *requestOptions) {
		opt.withIdentifyNo = true
		opt.identifyNoParams = params
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, opts ...RequestOption) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.config.Addr+path, nil)
	if err != nil {
		return nil, err
	}

	opt := newRequestOptions(opts...)

	for _, o := range opts {
		o(opt)
	}

	// withSystemParams
	if opt.withSystemParams {
		opt.postValues.Set("BusinessNo", c.config.BusinessNo)
		opt.postValues.Set("ApiVersion", c.config.APIVersion)
		opt.postValues.Set("HashKeyNo", c.config.HashKeyNo)
	}

	// withIdentifyNoRequestOption
	if opt.withIdentifyNo {
		var buffer bytes.Buffer
		for _, p := range opt.identifyNoParams {
			buffer.WriteString(opt.postValues.Get(p))
		}

		opt.postValues.Set("IdentifyNo", c.buildIdentifyNo(buffer.String()))
	}

	if opt.postValues != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Body = io.NopCloser(strings.NewReader(opt.postValues.Encode()))
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, dest any) error {
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
