package handler

import (
	"context"
	"fmt"
	"little-seed/admin/internal/apps/system/dto"
	"strings"

	"little-seed/admin/internal/apps/system/repo"
	"little-seed/admin/internal/apps/system/service"
)

type ConfigCreateReq struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
	Content     string `json:"content"`
}

type ConfigUpdateReq struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
	Content     string `json:"content"`
}

type ConfigDeleteReq struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
}

type ConfigGetReq struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
}

type ConfigMutationResp struct {
	Success bool `json:"success"`
}

type ServerConfigApi struct {
	svc *service.Service
}

func NewServerConfigApi(svc *service.Service) *ServerConfigApi {
	return &ServerConfigApi{svc: svc}
}

// Post 创建服务配置。
func (api *ServerConfigApi) Post(ctx context.Context, req *ConfigCreateReq) (*ConfigMutationResp, error) {
	return api.create(ctx, req)
}

// Put 更新服务配置。
func (api *ServerConfigApi) Put(ctx context.Context, req *ConfigUpdateReq) (*ConfigMutationResp, error) {
	return api.update(ctx, req)
}

// Delete 删除服务配置。
func (api *ServerConfigApi) Delete(ctx context.Context, req *ConfigDeleteReq) (*ConfigMutationResp, error) {
	return api.delete(ctx, req)
}

// Get 获取服务配置详情。
func (api *ServerConfigApi) Get(ctx context.Context, req *ConfigGetReq) (*repo.ConfigDetail, error) {
	return api.get(ctx, req)
}

// GetList 获取服务配置列表。
func (api *ServerConfigApi) GetList(ctx context.Context, req *struct{}) (*dto.ConfigListResp, error) {
	return api.list(ctx, req)
}

func (api *ServerConfigApi) create(ctx context.Context, req *ConfigCreateReq) (*ConfigMutationResp, error) {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("content is required")
	}

	err := api.svc.ServerConfig.Create(ctx, &dto.ConfigCreateReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
		Content:     req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigMutationResp{Success: true}, nil
}

func (api *ServerConfigApi) update(ctx context.Context, req *ConfigUpdateReq) (*ConfigMutationResp, error) {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("content is required")
	}

	err := api.svc.ServerConfig.Update(ctx, &dto.ConfigUpdateReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
		Content:     req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigMutationResp{Success: true}, nil
}

func (api *ServerConfigApi) delete(ctx context.Context, req *ConfigDeleteReq) (*ConfigMutationResp, error) {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return nil, err
	}

	err := api.svc.ServerConfig.Delete(ctx, &dto.ConfigDeleteReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigMutationResp{Success: true}, nil
}

func (api *ServerConfigApi) get(ctx context.Context, req *ConfigGetReq) (*repo.ConfigDetail, error) {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return nil, err
	}

	return api.svc.ServerConfig.Get(ctx, &dto.ConfigGetReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
	})
}

func (api *ServerConfigApi) list(ctx context.Context, req *struct{}) (*dto.ConfigListResp, error) {
	return api.svc.ServerConfig.List(ctx)
}

func checkConfigKey(serviceName, configName string) error {
	if strings.TrimSpace(serviceName) == "" {
		return fmt.Errorf("service_name is required")
	}
	if strings.TrimSpace(configName) == "" {
		return fmt.Errorf("config_name is required")
	}
	return nil
}
