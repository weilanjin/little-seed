package handler

import (
	"context"

	"little-seed/admin/internal/apps/system/service"
)

type ServerHealthApi struct {
	svc *service.Service
}

func NewServerHealthApi(svc *service.Service) *ServerHealthApi {
	return &ServerHealthApi{svc: svc}
}

// 获取注册服务列表
func (api *ServerHealthApi) GetList(ctx context.Context, req *struct{}) (*service.ServiceListResp, error) {
	return api.svc.ServerHealth.List(ctx)
}
