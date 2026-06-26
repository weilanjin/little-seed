package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"little-seed/kit/etcd"
)

type ServiceInstance struct {
	Name     string            `json:"name"`
	Addr     string            `json:"addr"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type ConfigSummary struct {
	ServiceName string `json:"service_name"`
	ConfigName  string `json:"config_name"`
	Key         string `json:"key"`
}

type ConfigDetail struct {
	ServiceName string    `json:"service_name"`
	ConfigName  string    `json:"config_name"`
	Key         string    `json:"key"`
	Content     string    `json:"content"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *Data) FindServiceList(ctx context.Context) ([]ServiceInstance, error) {
	client, err := d.newEtcdClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	resp, err := client.Get(ctx, d.etcdCfg.ServicePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	list := make([]ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var svc ServiceInstance
		if err := json.Unmarshal(kv.Value, &svc); err != nil {
			return nil, err
		}
		list = append(list, svc)
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Name == list[j].Name {
			return list[i].Addr < list[j].Addr
		}
		return list[i].Name < list[j].Name
	})
	return list, nil
}

func (d *Data) FindConfigList(ctx context.Context) ([]ConfigSummary, error) {
	client, err := d.newEtcdClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	resp, err := client.Get(ctx, d.configRoot, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	list := make([]ConfigSummary, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		serviceName, configName := d.parseConfigKey(string(kv.Key))
		list = append(list, ConfigSummary{
			ServiceName: serviceName,
			ConfigName:  configName,
			Key:         string(kv.Key),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].ServiceName == list[j].ServiceName {
			return list[i].ConfigName < list[j].ConfigName
		}
		return list[i].ServiceName < list[j].ServiceName
	})
	return list, nil
}

func (d *Data) GetConfig(ctx context.Context, serviceName, configName string) (*ConfigDetail, error) {
	client, err := d.newEtcdClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	key := d.configKey(serviceName, configName)
	resp, err := client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("config not found")
	}

	return &ConfigDetail{
		ServiceName: serviceName,
		ConfigName:  configName,
		Key:         key,
		Content:     string(resp.Kvs[0].Value),
		UpdatedAt:   time.Unix(0, resp.Kvs[0].ModRevision),
	}, nil
}

func (d *Data) CreateConfig(ctx context.Context, serviceName, configName, content string) error {
	client, err := d.newEtcdClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	key := d.configKey(serviceName, configName)
	txnResp, err := client.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, content)).
		Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return fmt.Errorf("config already exists")
	}
	return nil
}

func (d *Data) UpdateConfig(ctx context.Context, serviceName, configName, content string) error {
	client, err := d.newEtcdClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	key := d.configKey(serviceName, configName)
	txnResp, err := client.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), ">", 0)).
		Then(clientv3.OpPut(key, content)).
		Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return fmt.Errorf("config not found")
	}
	return nil
}

func (d *Data) DeleteConfig(ctx context.Context, serviceName, configName string) error {
	client, err := d.newEtcdClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	resp, err := client.Delete(ctx, d.configKey(serviceName, configName))
	if err != nil {
		return err
	}
	if resp.Deleted == 0 {
		return fmt.Errorf("config not found")
	}
	return nil
}

func (d *Data) newEtcdClient(ctx context.Context) (*clientv3.Client, error) {
	if !d.etcdCfg.Enabled() {
		return nil, fmt.Errorf("etcd endpoints are required")
	}
	return etcd.NewClient(ctx, d.etcdCfg)
}

func (d *Data) configKey(serviceName, configName string) string {
	return path.Join("/", strings.Trim(d.configRoot, "/"), strings.Trim(serviceName, "/"), strings.Trim(configName, "/"))
}

func (d *Data) parseConfigKey(key string) (string, string) {
	rel := strings.TrimPrefix(strings.Trim(key, "/"), strings.Trim(d.configRoot, "/"))
	parts := strings.Split(strings.Trim(rel, "/"), "/")
	if len(parts) < 2 {
		return "", strings.Trim(key, "/")
	}
	return parts[0], strings.Join(parts[1:], "/")
}
