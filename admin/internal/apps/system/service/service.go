package service

import "little-seed/admin/internal/apps/system/repo"

type Service struct {
	ServerHealth *ServerHealth
	ServerConfig *ServerConfig
	ServerLog    *ServerLog
}

func NewService(data *repo.Data) *Service {
	return &Service{
		ServerHealth: NewServerHealth(data),
		ServerConfig: NewServerConfig(data),
		ServerLog:    NewServerLog(data),
	}
}
