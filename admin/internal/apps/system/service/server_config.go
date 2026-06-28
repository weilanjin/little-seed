package service

import (
	"context"
	"little-seed/admin/internal/apps/system/dto"

	"little-seed/admin/internal/apps/system/repo"
)

type ServerConfig struct {
	data *repo.Data
}

func NewServerConfig(data *repo.Data) *ServerConfig {
	return &ServerConfig{
		data: data,
	}
}

func (s *ServerConfig) Create(ctx context.Context, req *dto.ConfigCreateReq) error {
	return s.data.EtcdCli.CreateConfig(ctx, req.ServiceName, req.ConfigName, req.Content)
}

func (s *ServerConfig) Update(ctx context.Context, req *dto.ConfigUpdateReq) error {
	return s.data.EtcdCli.UpdateConfig(ctx, req.ServiceName, req.ConfigName, req.Content)
}

func (s *ServerConfig) Delete(ctx context.Context, req *dto.ConfigDeleteReq) error {
	return s.data.EtcdCli.DeleteConfig(ctx, req.ServiceName, req.ConfigName)
}

func (s *ServerConfig) Get(ctx context.Context, req *dto.ConfigGetReq) (*repo.ConfigDetail, error) {
	return s.data.EtcdCli.GetConfig(ctx, req.ServiceName, req.ConfigName)
}

func (s *ServerConfig) List(ctx context.Context) (*dto.ConfigListResp, error) {
	list, err := s.data.EtcdCli.FindConfigList(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.ConfigListResp{List: list}, nil
}
