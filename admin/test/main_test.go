package test

import (
	"little-seed/admin/internal/apps/system"
	"little-seed/kit/etcd"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	system.Register(mux, etcd.Config{Endpoints: []string{"http://localhost:2379"}})
	http.DefaultServeMux = mux

	code := m.Run()
	os.Exit(code)
}
