package mid

import (
	"time"
)

type Action string

const (
	ValidateMSISDNAdvanceAction Action = "ValidateMSISDNAdvance"
)

type MIDInputParams interface {
	isMIDInputParams()
}

type ValidateMSISDNAdvanceMIDInputParams struct {
	Msisdn     string  `json:"Msisdn"`
	Birthday   *string `json:"Birthday"`
	ClauseVer  string  `json:"ClauseVer"`
	ClauseTime string  `json:"ClauseTime"`
}

func (*ValidateMSISDNAdvanceMIDInputParams) isMIDInputParams() {}

type MIDOutputParams interface {
	isMIDOutputParams()
}

type ValidateMSISDNAdvanceMIDOutputParams struct {
	MIDResp       *MIDResp `json:"-"`
	MIDRespString string   `json:"MIDResp"`
	VerifyCode    string   `json:"VerifyCode"`
	VerifyMsg     string   `json:"VerifyMsg"`
}

func (*ValidateMSISDNAdvanceMIDOutputParams) isMIDOutputParams() {}

type MIDResp struct {
	Code     string      `json:"code"`
	FullCode string      `json:"fullcode"`
	Message  string      `json:"message"`
	Msisdn   string      `json:"msisdn"`
	ReqSeq   string      `json:"reqSeq"`
	RspSeq   string      `json:"rspSeq"`
	RspTime  time.Time   `json:"rspTime"`
	SrvCode  string      `json:"srvCode"`
	Result   interface{} `json:"result"`
}
