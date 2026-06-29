package handler

import (
	"context"
	"fmt"
	"little-seed/admin/internal/apps/system/dto"
	"time"

	"little-seed/admin/internal/apps/system/repo"
	"little-seed/admin/internal/apps/system/service"
)

type ServerLogApi struct {
	svc *service.Service
}

func NewServerLogApi(svc *service.Service) *ServerLogApi {
	return &ServerLogApi{svc: svc}
}

// GetSearch 查询服务日志内容。
func (api *ServerLogApi) GetSearch(ctx context.Context, req *dto.LogSearchReq) (*service.LogSearchResp, error) {
	startAt, err := parseTime(req.StartAt)
	if err != nil {
		return nil, fmt.Errorf("start_at invalid: %w", err)
	}
	endAt, err := parseTime(req.EndAt)
	if err != nil {
		return nil, fmt.Errorf("end_at invalid: %w", err)
	}
	if !startAt.IsZero() && !endAt.IsZero() && startAt.After(endAt) {
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	return api.svc.ServerLog.Search(ctx, repo.LogQuery{
		FileName: req.FileName,
		Include:  req.Include,
		Exclude:  req.Exclude,
		StartAt:  startAt,
		EndAt:    endAt,
	})
}

// GetList 获取服务日志文件列表。
func (api *ServerLogApi) GetList(ctx context.Context, req *struct{}) (*service.LogListResp, error) {
	return api.svc.ServerLog.List(ctx, req)
}

func parseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
		t, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("use RFC3339 or 2006-01-02 15:04:05")
}
