package mid

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	MockBusinessNo = "11111111"
	MockHashKeyNo  = "10"
	MockHashKey    = "1234567890"
	MockAPIVersion = "1.0"
	MockClauseVer  = "1121215"
)

type MockMIDInputParams struct {
	Msisdn     string    `json:"Msisdn"`
	Birthday   *string   `json:"Birthday"`
	ClauseVer  string    `json:"ClauseVer"`
	ClauseTime time.Time `json:"ClauseTime"`
}

type MockInputParams struct {
	MemberNo       string              `json:"MemberNo"`
	Action         string              `json:"Action"`
	MIDInputParams *MockMIDInputParams `json:"MIDInputParams"`
}

func init() {
	time.Local = time.UTC
}

func TestClient(t *testing.T) {
	client := New(&Config{
		Addr:       "http://localhost",
		BusinessNo: MockBusinessNo,
		HashKeyNo:  MockHashKeyNo,
		HashKey:    MockHashKey,
	}, WithHTTPClient(http.DefaultClient))

	assert.IsType(t, &Client{}, client)
	assert.IsType(t, &http.Client{}, client.httpClient)
}
