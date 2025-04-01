package mid

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Action string

const (
	ValidateMSISDNAdvanceAction Action = "ValidateMSISDNAdvance"
)

type MIDInputParams struct { //nolint:revive
	Msisdn     string  `json:"Msisdn"`
	Birthday   *string `json:"Birthday"`
	ClauseVer  string  `json:"ClauseVer"`
	ClauseTime string  `json:"ClauseTime"`
}

type ServerSideTransactionRequest struct {
	VerifyNo       string          `json:"-"`
	MemberNo       string          `json:"MemberNo"`
	Action         Action          `json:"Action"`
	MIDInputParams *MIDInputParams `json:"MIDInputParams"`
}

type MIDResp struct { //nolint:revive
	Code     string    `json:"code"`
	FullCode string    `json:"fullcode"`
	Message  string    `json:"message"`
	Msisdn   string    `json:"msisdn"`
	ReqSeq   string    `json:"reqSeq"`
	RspSeq   string    `json:"rspSeq"`
	RspTime  time.Time `json:"rspTime"`
	SrvCode  string    `json:"srvCode"`
	Result   any       `json:"result"`
}

type MIDOutputParams struct { //nolint:revive
	MIDResp       *MIDResp `json:"-"`
	MIDRespString string   `json:"MIDResp"`
	VerifyCode    string   `json:"VerifyCode"`
	VerifyMsg     string   `json:"VerifyMsg"`
}

type OutputParams struct {
	MemberNo             string           `json:"MemberNo"`
	Token                string           `json:"Token"`
	TimeStamp            time.Time        `json:"TimeStamp"`
	MIDOutputParams      *MIDOutputParams `json:"-"`
	MIDOutputParamsBytes json.RawMessage  `json:"MIDOutputParams"` // for unmarshal
}

type ServerSideResponse struct {
	BusinessNo         string        `json:"BusinessNo"`
	APIVersion         string        `json:"ApiVersion"`
	HashKeyNo          string        `json:"HashKeyNo"`
	VerifyNo           string        `json:"VerifyNo"`
	ResultCode         string        `json:"ResultCode"`
	ReturnCode         string        `json:"ReturnCode"`
	ReturnCodeDesc     string        `json:"ReturnCodeDesc"`
	IdentifyNo         string        `json:"IdentifyNo"`
	OutputParams       *OutputParams `json:"-"`
	OutputParamsString string        `json:"OutputParams"` // for unmarshal
}

func (c *Client) ServerSideTransaction(
	ctx context.Context,
	req ServerSideTransactionRequest,
) (response ServerSideResponse, err error) {
	bytes, err := json.Marshal(req)
	if err != nil {
		return
	}

	opts := []RequestOption{
		withSystemParamsRequestOption(),
		withIdentifyNoRequestOption([]string{"BusinessNo", "ApiVersion", "HashKeyNo", "VerifyNo", "InputParams"}),
		setPostValueRequestOption("InputParams", string(bytes)),
		setPostValueRequestOption("VerifyNo", req.VerifyNo),
	}

	request, err := c.newRequest(ctx, http.MethodPost, "/IDPortal/ServerSideTransaction", opts...)
	if err != nil {
		return
	}

	if err = c.doRequest(request, &response); err != nil {
		return
	}

	err = c.decodeServerSideResponse(&response)
	return
}

type ServerSideVerifyResultRequest struct {
	VerifyNo string `json:"VerifyNo"`
	MemberNo string `json:"MemberNo"`
	Token    string `json:"Token"`
}

func (c *Client) ServerSideVerifyResult(
	ctx context.Context,
	req ServerSideVerifyResultRequest,
) (response ServerSideResponse, err error) {
	opts := []RequestOption{
		withSystemParamsRequestOption(),
		withIdentifyNoRequestOption([]string{"BusinessNo", "ApiVersion", "HashKeyNo", "VerifyNo", "MemberNo", "Token"}),
		setPostValueRequestOption("MemberNo", req.MemberNo),
		setPostValueRequestOption("Token", req.Token),
		setPostValueRequestOption("VerifyNo", req.VerifyNo),
	}

	request, err := c.newRequest(ctx, http.MethodPost, "/IDPortal/ServerSideVerifyResult", opts...)
	if err != nil {
		return
	}

	if err = c.doRequest(request, &response); err != nil {
		return
	}

	err = c.decodeServerSideResponse(&response)
	return
}

func (c *Client) decodeServerSideResponse(response *ServerSideResponse) (err error) {
	if response.OutputParamsString == "" {
		return
	}

	// unmarshal OutputParamsBytes to OutputParams
	var outputParams OutputParams
	if err = json.Unmarshal([]byte(response.OutputParamsString), &outputParams); err != nil {
		return
	}
	response.OutputParams = &outputParams

	// validate msisdn advance
	if outputParams.MIDOutputParamsBytes != nil {
		var midOutputParams MIDOutputParams
		if err = json.Unmarshal(outputParams.MIDOutputParamsBytes, &midOutputParams); err != nil {
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

	return
}
