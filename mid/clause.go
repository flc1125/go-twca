package mid

import (
	"context"
	"net/http"
	"time"
)

type MIDClauseResponse struct {
	Fullcode   string    `json:"fullcode"`
	SrvCode    string    `json:"srvCode"`
	Code       string    `json:"code"`
	LastUpdate time.Time `json:"lastUpdate"`
	Clausever  string    `json:"clausever"`
	RspTime    time.Time `json:"rspTime"`
	Html       string    `json:"html"`
	Lang       string    `json:"lang"`
	Message    string    `json:"message"`
	Local      bool      `json:"local"`
}

func (c *Client) MIDClause(ctx context.Context) (response MIDClauseResponse, err error) {
	request, err := c.newRequest(ctx, http.MethodPost, "/IDPortal/MIDClause")
	if err != nil {
		return
	}

	err = c.doRequest(request, &response)
	return
}
