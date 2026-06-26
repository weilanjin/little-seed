package system

import (
	"net/http"

	"little-seed/admin/internal/apps/system/handler"
	"little-seed/admin/internal/apps/system/repo"
	"little-seed/admin/internal/apps/system/service"
	"little-seed/kit/core/hs"
	"little-seed/kit/etcd"
)

func Register(mux *http.ServeMux, etcdCfg etcd.Config) {
	data := repo.NewData(etcdCfg)
	svc := service.NewService(data)
	h := handler.NewHandler(svc)

	group := hs.NewGroup("/api/system", mux)
	hs.RegisterService(group, "/server/health", h.ServerHealth)
	hs.RegisterService(group, "/server/config", h.ServerConfig)
	hs.RegisterService(group, "/server/log", h.ServerLog)
}
