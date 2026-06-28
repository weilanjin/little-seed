package dto

import "little-seed/admin/internal/apps/system/repo"

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
