package mid

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ServerSideTransactionRequest struct {
	MemberNo       string         `json:"MemberNo"`
	Action         Action         `json:"Action"`
	MIDInputParams MIDInputParams `json:"MIDInputParams"`
}

type ServerSideTransactionOutputParams struct {
	MemberNo             string          `json:"MemberNo"`
	Token                string          `json:"Token"`
	TimeStamp            time.Time       `json:"TimeStamp"`
	MIDOutputParams      MIDOutputParams `json:"-"`
	MIDOutputParamsBytes json.RawMessage `json:"MIDOutputParams"` // for unmarshal
}

type ServerSideTransactionResponse struct {
	BusinessNo         string                             `json:"BusinessNo"`
	ApiVersion         string                             `json:"ApiVersion"`
	HashKeyNo          string                             `json:"HashKeyNo"`
	VerifyNo           string                             `json:"VerifyNo"`
	ResultCode         string                             `json:"ResultCode"`
	ReturnCode         string                             `json:"ReturnCode"`
	ReturnCodeDesc     string                             `json:"ReturnCodeDesc"`
	IdentifyNo         string                             `json:"IdentifyNo"`
	OutputParams       *ServerSideTransactionOutputParams `json:"-"`
	OutputParamsString string                             `json:"OutputParams"` // for unmarshal
}

func (c *Client) ServerSideTransaction(ctx context.Context, req ServerSideTransactionRequest, opts ...RequestOption) (response ServerSideTransactionResponse, err error) {
	bytes, err := json.Marshal(req)
	if err != nil {
		return
	}

	opts = append(opts,
		withSystemParamsRequestOption(),
		withIdentifyNoRequestOption([]string{"BusinessNo", "ApiVersion", "HashKeyNo", "VerifyNo", "InputParams"}),
		setPostValueRequestOption("InputParams", string(bytes)),
	)

	request, err := c.newRequest(ctx, http.MethodPost, "/IDPortal/ServerSideTransaction", opts...)
	if err != nil {
		return
	}

	if err = c.doRequest(request, &response); err != nil {
		return
	}

	if response.OutputParamsString == "" {
		return
	}

	// unmarshal OutputParamsBytes to OutputParams
	var outputParams ServerSideTransactionOutputParams
	if err = json.Unmarshal([]byte(response.OutputParamsString), &outputParams); err != nil {
		return
	}
	response.OutputParams = &outputParams

	// validate msisdn advance
	if req.Action == ValidateMSISDNAdvanceAction {
		if outputParams.MIDOutputParamsBytes != nil {
			var midOutputParams ValidateMSISDNAdvanceMIDOutputParams
			if err = json.Unmarshal([]byte(outputParams.MIDOutputParamsBytes), &midOutputParams); err != nil {
				return
			}
			outputParams.MIDOutputParams = &midOutputParams

			if midOutputParams.MIDRespString != "" {
				var midResp MIDResp
				if err = json.Unmarshal([]byte(midOutputParams.MIDRespString), &midResp); err != nil {
					return
				}
				midOutputParams.MIDResp = &midResp
			}
		}
	}

	return
}
