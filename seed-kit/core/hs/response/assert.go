package response

import (
	"encoding/json"
	"little-seed/kit/core/hs/response/codes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 用于测试
func Assert(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Code != http.StatusOK {
		t.Fatalf("%s: expected status 200, got %d, body: %s", t.Name(), w.Code, w.Body.String())
	}

	resp := new(Resp)
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("%s: failed to unmarshal response: %v", t.Name(), err)
	}

	byt, _ := json.MarshalIndent(resp, "", "  ")

	// 检查返回的数据格式
	if codes.Code(resp.Code) != codes.OK {
		t.Fatalf("%s: expected code 0, got:\n %s", t.Name(), string(byt))
	}
	t.Logf("✅ %s passed, response:\n %s", t.Name(), string(byt))
}
