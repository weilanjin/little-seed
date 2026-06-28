package test

import (
	"bytes"
	"encoding/json"
	"little-seed/kit/core/hs/response"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

// TestServerHealth 测试服务健康列表接口。
func TestServerHealth(t *testing.T) {

	// 获取服务列表
	t.Run("List", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/server/health/list", nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		response.Assert(t, w)
	})
}

// TestServerLog 测试服务日志查询接口。
func TestServerLog(t *testing.T) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("failed to create log dir: %v", err)
	}

	logName := "system_test.log"
	logPath := filepath.Join(logDir, logName)
	content := "2026-06-28 10:00:00 api started\n2026-06-28 10:01:00 api failed\n"
	if err := os.WriteFile(logPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(logPath); err != nil && !os.IsNotExist(err) {
			t.Logf("failed to remove log file: %v", err)
		}
	})

	t.Run("List", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/server/log/list", nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		response.Assert(t, w)
	})

	t.Run("Search", func(t *testing.T) {
		query := url.Values{}
		query.Set("file_name", logName)
		query.Set("include", "failed")

		req := httptest.NewRequest(http.MethodGet, "/api/system/server/log/search?"+query.Encode(), nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		response.Assert(t, w)
	})
}

// TestServerConfig 测试服务配置管理接口。
func TestServerConfig(t *testing.T) {

	t.Run("List", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/server/config/list", nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		response.Assert(t, w)
	})

	t.Run("Create", func(t *testing.T) {
		reqBody := map[string]string{
			"service_name": "admin",
			"config_name":  "system_test.yaml",
			"content":      "name: admin",
		}
		req := newJSONRequest(t, http.MethodPost, "/api/system/server/config", reqBody)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		response.Assert(t, w)
	})

	t.Run("Update", func(t *testing.T) {
		reqBody := map[string]string{
			"service_name": "admin",
			"config_name":  "system_test.yaml",
			"content":      "name: admin-test",
		}
		req := newJSONRequest(t, http.MethodPut, "/api/system/server/config", reqBody)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		response.Assert(t, w)
	})

	t.Run("Get", func(t *testing.T) {
		query := url.Values{}
		query.Set("service_name", "admin")
		query.Set("config_name", "system_test.yaml")

		req := httptest.NewRequest(http.MethodGet, "/api/system/server/config?"+query.Encode(), nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		response.Assert(t, w)
	})

	t.Run("Delete", func(t *testing.T) {
		query := url.Values{}
		query.Set("service_name", "admin")
		query.Set("config_name", "system_test.yaml")

		req := httptest.NewRequest(http.MethodDelete, "/api/system/server/config?"+query.Encode(), nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		response.Assert(t, w)
	})
}

func newJSONRequest(t *testing.T, method, target string, body any) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Fatalf("failed to encode request body: %v", err)
	}

	req := httptest.NewRequest(method, target, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}
