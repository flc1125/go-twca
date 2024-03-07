package mid

import (
	"context"
	"net/http"
	"time"
)

type ClauseResponse struct {
	FullCode   string    `json:"fullcode"`
	SrvCode    string    `json:"srvCode"`
	Code       string    `json:"code"`
	LastUpdate time.Time `json:"lastUpdate"`
	ClauseVer  string    `json:"clausever"`
	RspTime    time.Time `json:"rspTime"`
	HTML       string    `json:"html"`
	Lang       string    `json:"lang"`
	Message    string    `json:"message"`
	Local      bool      `json:"local"`
}

func (c *Client) Clause(ctx context.Context) (response ClauseResponse, err error) {
	request, err := c.newRequest(ctx, http.MethodPost, "/IDPortal/MIDClause")
	if err != nil {
		return
	}

	err = c.doRequest(request, &response)
	return
}
