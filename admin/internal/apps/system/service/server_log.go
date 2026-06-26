package service

import (
	"context"

	"little-seed/admin/internal/apps/system/repo"
)

type ServerLog struct {
	data *repo.Data
}

type LogSearchReq struct {
	FileName string
	Include  string
	Exclude  string
	StartAt  string
	EndAt    string
}

type LogListResp struct {
	List []repo.LogFile `json:"list"`
}

type LogSearchResp struct {
	List []repo.LogLine `json:"list"`
}

func NewServerLog(data *repo.Data) *ServerLog {
	return &ServerLog{
		data: data,
	}
}

func (s *ServerLog) Search(ctx context.Context, query repo.LogQuery) (*LogSearchResp, error) {
	lines, err := s.data.SearchLogs(query)
	if err != nil {
		return nil, err
	}
	return &LogSearchResp{List: lines}, nil
}

func (s *ServerLog) List(ctx context.Context, req *struct{}) (*LogListResp, error) {
	logs, err := s.data.FindLogList()
	if err != nil {
		return nil, err
	}
	return &LogListResp{List: logs}, nil
}
