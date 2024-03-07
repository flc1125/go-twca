package mid

import "time"

type Action string

const (
	ValidateMSISDNAdvanceAction Action = "ValidateMSISDNAdvanced"
)

type MIDInputParams interface {
	isMIDInputParams()
}

type ValidateMSISDNAdvanceMIDInputParams struct {
	Msisdn     string    `json:"Msisdn"`
	Birthday   *string   `json:"Birthday"`
	ClauseVer  string    `json:"ClauseVer"`
	ClauseTime time.Time `json:"ClauseTime"`
}

func (*ValidateMSISDNAdvanceMIDInputParams) isMIDInputParams() {}
