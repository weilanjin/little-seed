package handler

import (
	"context"
	"fmt"
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

func (api *ServerConfigApi) Create(ctx context.Context, req *ConfigCreateReq) (*ConfigMutationResp, error) {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("content is required")
	}

	err := api.svc.ServerConfig.Create(ctx, service.ConfigCreateReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
		Content:     req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigMutationResp{Success: true}, nil
}

func (api *ServerConfigApi) Update(ctx context.Context, req *ConfigUpdateReq) (*ConfigMutationResp, error) {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("content is required")
	}

	err := api.svc.ServerConfig.Update(ctx, service.ConfigUpdateReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
		Content:     req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigMutationResp{Success: true}, nil
}

func (api *ServerConfigApi) Delete(ctx context.Context, req *ConfigDeleteReq) (*ConfigMutationResp, error) {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return nil, err
	}

	err := api.svc.ServerConfig.Delete(ctx, service.ConfigDeleteReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigMutationResp{Success: true}, nil
}

func (api *ServerConfigApi) Get(ctx context.Context, req *ConfigGetReq) (*repo.ConfigDetail, error) {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return nil, err
	}

	return api.svc.ServerConfig.Get(ctx, service.ConfigGetReq{
		ServiceName: req.ServiceName,
		ConfigName:  req.ConfigName,
	})
}

func (api *ServerConfigApi) List(ctx context.Context, req *struct{}) (*service.ConfigListResp, error) {
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
