package dto

import (
	"fmt"
	"little-seed/admin/internal/apps/system/repo"
	"strings"
)

type ConfigListResp struct {
	List []repo.ConfigSummary `json:"list"`
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

type ConfigCreateReq struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
	Content     string `json:"content"`
}

func (req *ConfigCreateReq) Validate() error {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return err
	}
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}

type ConfigUpdateReq struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
	Content     string `json:"content"`
}

func (req *ConfigUpdateReq) Validate() error {
	if err := checkConfigKey(req.ServiceName, req.ConfigName); err != nil {
		return err
	}
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}

type ConfigDeleteReq struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
}

func (req *ConfigDeleteReq) Validate() error {
	return checkConfigKey(req.ServiceName, req.ConfigName)
}

type ConfigGetReq struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
}

func (req *ConfigGetReq) Validate() error {
	return checkConfigKey(req.ServiceName, req.ConfigName)
}

type ConfigMutationResp struct {
	Success bool `json:"success"`
}
