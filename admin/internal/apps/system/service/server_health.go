package service

import (
	"context"

	"little-seed/admin/internal/apps/system/repo"
)

type ServerHealth struct {
	data *repo.Data
}

type ServiceListResp struct {
	List []repo.ServiceInstance `json:"list"`
}

func NewServerHealth(data *repo.Data) *ServerHealth {
	return &ServerHealth{
		data: data,
	}
}

func (s *ServerHealth) List(ctx context.Context) (*ServiceListResp, error) {
	services, err := s.data.EtcdCli.FindServiceList(ctx)
	if err != nil {
		return nil, err
	}
	return &ServiceListResp{List: services}, nil
}
