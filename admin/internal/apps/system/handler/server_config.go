package handler

import (
	"context"
	"little-seed/admin/internal/apps/system/dto"
	"little-seed/admin/internal/apps/system/repo"
	"little-seed/admin/internal/apps/system/service"
)

type ServerConfigApi struct {
	svc *service.Service
}

func NewServerConfigApi(svc *service.Service) *ServerConfigApi {
	return &ServerConfigApi{svc: svc}
}

// Post 创建服务配置。
func (api *ServerConfigApi) Post(ctx context.Context, req *dto.ConfigCreateReq) (*dto.ConfigMutationResp, error) {
	err := api.svc.ServerConfig.Create(ctx, &dto.ConfigCreateReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
		Content:     req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &dto.ConfigMutationResp{Success: true}, nil
}

// Put 更新服务配置。
func (api *ServerConfigApi) Put(ctx context.Context, req *dto.ConfigUpdateReq) (*dto.ConfigMutationResp, error) {
	err := api.svc.ServerConfig.Update(ctx, &dto.ConfigUpdateReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
		Content:     req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &dto.ConfigMutationResp{Success: true}, nil
}

// Delete 删除服务配置。
func (api *ServerConfigApi) Delete(ctx context.Context, req *dto.ConfigDeleteReq) (*dto.ConfigMutationResp, error) {
	err := api.svc.ServerConfig.Delete(ctx, &dto.ConfigDeleteReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
	})
	if err != nil {
		return nil, err
	}
	return &dto.ConfigMutationResp{Success: true}, nil
}

// Get 获取服务配置详情。
func (api *ServerConfigApi) Get(ctx context.Context, req *dto.ConfigGetReq) (*repo.ConfigDetail, error) {
	return api.svc.ServerConfig.Get(ctx, &dto.ConfigGetReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
	})
}

// GetList 获取服务配置列表。
func (api *ServerConfigApi) GetList(ctx context.Context, req *struct{}) (*dto.ConfigListResp, error) {
	return api.svc.ServerConfig.List(ctx)
}
