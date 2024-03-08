package mid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClause(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/IDPortal/MIDClause", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded; charset=utf-8", r.Header.Get("Content-Type"))

		_, err := w.Write([]byte(`{
  "fullcode": "0",
  "srvCode": "364",
  "code": "0000",
  "lastUpdate": "2024-03-07T16:00:01Z",
  "clausever": "1121215",
  "rspTime": "2024-03-07T16:25:30Z",
  "html": "Mobile ID",
  "lang": "zh-Hant",
  "message": "success",
  "local": false
}`))
		assert.NoError(t, err)
	}))
	defer srv.Close()

	client := New(&Config{
		Addr:       srv.URL,
		BusinessNo: MockBusinessNo,
		HashKeyNo:  MockHashKeyNo,
		HashKey:    MockHashKey,
		APIVersion: MockAPIVersion,
	})

	resp, err := client.Clause(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "0", resp.FullCode)
	assert.Equal(t, "364", resp.SrvCode)
	assert.Equal(t, "0000", resp.Code)
	assert.Equal(t, "2024-03-07T16:00:01Z", resp.LastUpdate.Format(time.RFC3339))
	assert.Equal(t, "1121215", resp.ClauseVer)
	assert.Equal(t, "Mobile ID", resp.HTML)
	assert.Equal(t, "zh-Hant", resp.Lang)
	assert.Equal(t, "2024-03-07T16:25:30Z", resp.RspTime.Format(time.RFC3339))
	assert.Equal(t, false, resp.Local)
	assert.Equal(t, "success", resp.Message)
}
