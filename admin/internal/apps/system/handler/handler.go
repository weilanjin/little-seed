package handler

import (
	"little-seed/admin/internal/apps/system/service"
)

type Handler struct {
	ServerHealth *ServerHealthApi
	ServerConfig *ServerConfigApi
	ServerLog    *ServerLogApi
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		ServerHealth: NewServerHealthApi(svc),
		ServerConfig: NewServerConfigApi(svc),
		ServerLog:    NewServerLogApi(svc),
	}
}
