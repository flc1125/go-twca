package mid

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestClient_Clause(t *testing.T) {
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

func TestClient_ServerSideTransaction(t *testing.T) {
	tests := []struct {
		MemberNo     string
		Msisdn       string
		ResultCode   string
		ResponseJSON string
	}{
		{
			MemberNo:     "A123456789",
			Msisdn:       "0900005002",
			ResultCode:   "5002",
			ResponseJSON: `{"BusinessNo":"` + MockBusinessNo + `","ApiVersion":"` + MockAPIVersion + `","HashKeyNo":"` + MockHashKeyNo + `","VerifyNo":"5002","ResultCode":"F","ReturnCode":"3645002","ReturnCodeDesc":"mno response: invalid msisdn","IdentifyNo":"5fe8cad322f7edd8ca34be925c39de1af461037d2fc4ab886033bc0e0d1dba17","OutputParams":"{\"MemberNo\":\"A123456789\",\"Token\":\"srv8b13776bf752463eab91863cc0b1ff9999008\",\"TimeStamp\":\"2024-03-08T00:30:17.3120000Z\",\"MIDOutputParams\":{\"MIDResp\":\"{\\\"code\\\":\\\"5002\\\",\\\"fullcode\\\":\\\"3645002\\\",\\\"message\\\":\\\"mno response: invalid msisdn\\\",\\\"msisdn\\\":\\\"0900005002\\\",\\\"reqSeq\\\":\\\"01170985781723350828466605\\\",\\\"rspSeq\\\":\\\"3257402\\\",\\\"rspTime\\\":\\\"2024-03-08T00:30:17Z\\\",\\\"srvCode\\\":\\\"364\\\",\\\"result\\\":{}}\",\"VerifyCode\":\"5002\",\"VerifyMsg\":\"mno response: invalid msisdn\"}}"}`, //nolint:lll
		},
		{
			MemberNo:     "A123456789",
			Msisdn:       "0900005011",
			ResultCode:   "5011",
			ResponseJSON: `{"BusinessNo":"` + MockBusinessNo + `","ApiVersion":"` + MockAPIVersion + `","HashKeyNo":"` + MockHashKeyNo + `","VerifyNo":"5011","ResultCode":"F","ReturnCode":"3645011","ReturnCodeDesc":"mno response: invalid rocid","IdentifyNo":"70434d12c90deaaf84481092c3abc82fe976b4e58c8552343a2ccee536366894","OutputParams":"{\"MemberNo\":\"A123456789\",\"Token\":\"srv888d21c798e14df69025220c7c48cafeba3ed\",\"TimeStamp\":\"2024-03-08T00:30:17.4620000Z\",\"MIDOutputParams\":{\"MIDResp\":\"{\\\"code\\\":\\\"5011\\\",\\\"fullcode\\\":\\\"3645011\\\",\\\"message\\\":\\\"mno response: invalid rocid\\\",\\\"msisdn\\\":\\\"0900005011\\\",\\\"reqSeq\\\":\\\"01170985781739450928466605\\\",\\\"rspSeq\\\":\\\"3257403\\\",\\\"rspTime\\\":\\\"2024-03-08T00:30:17Z\\\",\\\"srvCode\\\":\\\"364\\\",\\\"result\\\":{}}\",\"VerifyCode\":\"5011\",\"VerifyMsg\":\"mno response: invalid rocid\"}}"}`, //nolint:lll
		},
		{
			MemberNo:     "A123456789",
			Msisdn:       "0900005016",
			ResultCode:   "5016",
			ResponseJSON: `{"BusinessNo":"` + MockBusinessNo + `","ApiVersion":"` + MockAPIVersion + `","HashKeyNo":"` + MockHashKeyNo + `","VerifyNo":"5016","ResultCode":"F","ReturnCode":"3645016","ReturnCodeDesc":"mno response: Valid ID but invalid birthday","IdentifyNo":"9c7abb440b2f36ef606ba7c9801c631bafd084c60b4c06cb3088218f9736cd57","OutputParams":"{\"MemberNo\":\"A123456789\",\"Token\":\"srv8599802aad9e4c679168f2cce8e6add92ac5a\",\"TimeStamp\":\"2024-03-08T00:30:17.6350000Z\",\"MIDOutputParams\":{\"MIDResp\":\"{\\\"code\\\":\\\"5016\\\",\\\"fullcode\\\":\\\"3645016\\\",\\\"message\\\":\\\"mno response: Valid ID but invalid birthday\\\",\\\"msisdn\\\":\\\"0900005016\\\",\\\"reqSeq\\\":\\\"01170985781754351028466605\\\",\\\"rspSeq\\\":\\\"3257404\\\",\\\"rspTime\\\":\\\"2024-03-08T00:30:17Z\\\",\\\"srvCode\\\":\\\"364\\\",\\\"result\\\":{}}\",\"VerifyCode\":\"5016\",\"VerifyMsg\":\"mno response: Valid ID but invalid birthday\"}}"}`, //nolint:lll
		},
		{
			MemberNo:     "A123456789",
			Msisdn:       "0900005021",
			ResultCode:   "5021",
			ResponseJSON: `{"BusinessNo":"` + MockBusinessNo + `","ApiVersion":"` + MockAPIVersion + `","HashKeyNo":"` + MockHashKeyNo + `","VerifyNo":"5021","ResultCode":"F","ReturnCode":"3645021","ReturnCodeDesc":"mno response: status bad","IdentifyNo":"0ab94a79968e5475455053dd96f3b74521c711af58c4a0b1a6980e49e827c559","OutputParams":"{\"MemberNo\":\"A123456789\",\"Token\":\"srvf58cd91c5dcb4ee3b17263af9bb445d911eed\",\"TimeStamp\":\"2024-03-08T00:30:17.7790000Z\",\"MIDOutputParams\":{\"MIDResp\":\"{\\\"code\\\":\\\"5021\\\",\\\"fullcode\\\":\\\"3645021\\\",\\\"message\\\":\\\"mno response: status bad\\\",\\\"msisdn\\\":\\\"0900005021\\\",\\\"reqSeq\\\":\\\"01170985781770951128466605\\\",\\\"rspSeq\\\":\\\"3257405\\\",\\\"rspTime\\\":\\\"2024-03-08T00:30:17Z\\\",\\\"srvCode\\\":\\\"364\\\",\\\"result\\\":{}}\",\"VerifyCode\":\"5021\",\"VerifyMsg\":\"mno response: status bad\"}}"}`, //nolint:lll
		},
		{
			MemberNo:     "A123456789",
			Msisdn:       "0900005024",
			ResultCode:   "5024",
			ResponseJSON: `{"BusinessNo":"` + MockBusinessNo + `","ApiVersion":"` + MockAPIVersion + `","HashKeyNo":"` + MockHashKeyNo + `","VerifyNo":"5024","ResultCode":"F","ReturnCode":"3645024","ReturnCodeDesc":"mno response: msisdn not accept","IdentifyNo":"d1ef73f9d7900b28f5ad9772070f362bbd6e990267c2b3f9866684b4ddaf555d","OutputParams":"{\"MemberNo\":\"A123456789\",\"Token\":\"srv2292b2ed7e1d47c3a086d53684fb2c21bb03a\",\"TimeStamp\":\"2024-03-08T00:30:17.9540000Z\",\"MIDOutputParams\":{\"MIDResp\":\"{\\\"code\\\":\\\"5024\\\",\\\"fullcode\\\":\\\"3645024\\\",\\\"message\\\":\\\"mno response: msisdn not accept\\\",\\\"msisdn\\\":\\\"0900005024\\\",\\\"reqSeq\\\":\\\"01170985781785151228466605\\\",\\\"rspSeq\\\":\\\"3257406\\\",\\\"rspTime\\\":\\\"2024-03-08T00:30:17Z\\\",\\\"srvCode\\\":\\\"364\\\",\\\"result\\\":{}}\",\"VerifyCode\":\"5024\",\"VerifyMsg\":\"mno response: msisdn not accept\"}}"}`, //nolint:lll
		},
		{
			MemberNo:     "A130882986",
			Msisdn:       "0965277111",
			ResultCode:   "0000",
			ResponseJSON: `{"BusinessNo":"` + MockBusinessNo + `","ApiVersion":"` + MockAPIVersion + `","HashKeyNo":"` + MockHashKeyNo + `","VerifyNo":"0000","ResultCode":"S","ReturnCode":"0","ReturnCodeDesc":"成功","IdentifyNo":"21cc8cbe50e1427b5df004ffcc6f671ec2330baebd8f448e02e8c0870af348e8","OutputParams":"{\"MemberNo\":\"A130882986\",\"Token\":\"srv0a7b2f24edc34aa1b890a1d25368f86f42547\",\"TimeStamp\":\"2024-03-08T00:30:18.1040000Z\",\"MIDOutputParams\":{\"MIDResp\":\"{\\\"code\\\":\\\"0000\\\",\\\"fullcode\\\":\\\"0\\\",\\\"message\\\":\\\"success\\\",\\\"msisdn\\\":\\\"0965277111\\\",\\\"reqSeq\\\":\\\"01170985781802851328466605\\\",\\\"rspSeq\\\":\\\"3257407\\\",\\\"rspTime\\\":\\\"2024-03-08T00:30:18Z\\\",\\\"srvCode\\\":\\\"364\\\",\\\"result\\\":{}}\",\"VerifyCode\":\"0000\",\"VerifyMsg\":\"success\"}}"}`, //nolint:lll
		},
	}

	for _, test := range tests {
		t.Run(test.ResultCode, func(t *testing.T) {
			// transaction
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/IDPortal/ServerSideTransaction", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/x-www-form-urlencoded; charset=utf-8", r.Header.Get("Content-Type"))
				assert.Equal(t, MockBusinessNo, r.FormValue("BusinessNo"))
				assert.Equal(t, MockAPIVersion, r.FormValue("ApiVersion"))
				assert.Equal(t, MockHashKeyNo, r.FormValue("HashKeyNo"))
				assert.Equal(t, test.ResultCode, r.FormValue("VerifyNo"))
				assert.NotEmpty(t, r.FormValue("IdentifyNo"))

				var inputParams MockInputParams
				assert.NoError(t, json.Unmarshal([]byte(r.FormValue("InputParams")), &inputParams))
				assert.NotNil(t, inputParams)
				assert.Equal(t, test.MemberNo, inputParams.MemberNo)
				assert.Equal(t, "ValidateMSISDNAdvance", inputParams.Action)
				assert.NotNil(t, inputParams.MIDInputParams)
				assert.Equal(t, test.Msisdn, inputParams.MIDInputParams.Msisdn)
				assert.Equal(t, MockClauseVer, inputParams.MIDInputParams.ClauseVer)
				assert.NotEmpty(t, inputParams.MIDInputParams.ClauseTime)

				_, err := w.Write([]byte(test.ResponseJSON))
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

			resp, err := client.ServerSideTransaction(context.Background(), ServerSideTransactionRequest{
				VerifyNo: test.ResultCode,
				MemberNo: test.MemberNo,
				Action:   ValidateMSISDNAdvanceAction,
				MIDInputParams: &MIDInputParams{
					Msisdn:     test.Msisdn,
					Birthday:   nil,
					ClauseVer:  MockClauseVer,
					ClauseTime: time.Now().Format(time.RFC3339),
				},
			})
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.OutputParams)
			assert.NotNil(t, resp.OutputParams.MIDOutputParams)
			assert.NotNil(t, resp.OutputParams.MIDOutputParams.MIDResp)
			assert.Equal(t, test.ResultCode, resp.OutputParams.MIDOutputParams.MIDResp.Code)
		})
	}
}
