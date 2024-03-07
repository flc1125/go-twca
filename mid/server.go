package mid

import (
	"context"
	"net/http"
)

type ServerSideTransactionRequest struct {
	MemberNo       string         `json:"MemberNo"`
	Action         Action         `json:"Action"`
	MIDInputParams MIDInputParams `json:"MIDInputParams"`
}

type ServerSideTransactionResponse struct {
}

func (c *Client) ServerSideTransaction(ctx context.Context, req ServerSideTransactionRequest) (response ServerSideTransactionResponse, err error) {
	request, err := c.newRequest(ctx, http.MethodPost, "/IDPortal/ServerSideTransaction")
	if err != nil {
		return
	}

	err = c.doRequest(request, &response)
	return
}
