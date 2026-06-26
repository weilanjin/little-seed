package service

import (
	"context"

	"little-seed/admin/internal/apps/system/repo"
)

type ServerConfig struct {
	data *repo.Data
}

type ConfigListResp struct {
	List []repo.ConfigSummary `json:"list"`
}

type ConfigCreateReq struct {
	ServiceName string
	ConfigName  string
	Content     string
}

type ConfigUpdateReq struct {
	ServiceName string
	ConfigName  string
	Content     string
}

type ConfigDeleteReq struct {
	ServiceName string
	ConfigName  string
}

type ConfigGetReq struct {
	ServiceName string
	ConfigName  string
}

func NewServerConfig(data *repo.Data) *ServerConfig {
	return &ServerConfig{
		data: data,
	}
}

func (s *ServerConfig) Create(ctx context.Context, req ConfigCreateReq) error {
	return s.data.CreateConfig(ctx, req.ServiceName, req.ConfigName, req.Content)
}

func (s *ServerConfig) Update(ctx context.Context, req ConfigUpdateReq) error {
	return s.data.UpdateConfig(ctx, req.ServiceName, req.ConfigName, req.Content)
}

func (s *ServerConfig) Delete(ctx context.Context, req ConfigDeleteReq) error {
	return s.data.DeleteConfig(ctx, req.ServiceName, req.ConfigName)
}

func (s *ServerConfig) Get(ctx context.Context, req ConfigGetReq) (*repo.ConfigDetail, error) {
	return s.data.GetConfig(ctx, req.ServiceName, req.ConfigName)
}

func (s *ServerConfig) List(ctx context.Context) (*ConfigListResp, error) {
	list, err := s.data.FindConfigList(ctx)
	if err != nil {
		return nil, err
	}
	return &ConfigListResp{List: list}, nil
}
