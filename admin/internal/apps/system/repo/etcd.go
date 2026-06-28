package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
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

type EctdRepo struct {
	cli *clientv3.Client
}

func NewEctdClient(ctx context.Context, etcdCfg etcd.Config) *EctdRepo {
	client, err := etcd.NewClient(ctx, etcdCfg)
	if err != nil {
		log.Fatalf("failed to create etcd client: %v", err)
	}
	return &EctdRepo{
		cli: client,
	}
}

func (repo *EctdRepo) FindServiceList(ctx context.Context) ([]ServiceInstance, error) {
	resp, err := repo.cli.Get(ctx, "service", clientv3.WithPrefix())
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

func (repo *EctdRepo) FindConfigList(ctx context.Context) ([]ConfigSummary, error) {
	resp, err := repo.cli.Get(ctx, "config", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	list := make([]ConfigSummary, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		list = append(list, ConfigSummary{
			ServiceName: "",
			ConfigName:  "",
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

func (repo *EctdRepo) GetConfig(ctx context.Context, serviceName, configName string) (*ConfigDetail, error) {
	//key := repo.configKey(serviceName, configName)
	resp, err := repo.cli.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("config not found")
	}

	return &ConfigDetail{
		ServiceName: serviceName,
		ConfigName:  configName,
		//Key:         key,
		Content:   string(resp.Kvs[0].Value),
		UpdatedAt: time.Unix(0, resp.Kvs[0].ModRevision),
	}, nil
}

func (repo *EctdRepo) CreateConfig(ctx context.Context, serviceName, configName, content string) error {
	//key := repo.configKey(serviceName, configName)
	txnResp, err := repo.cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(""), "=", 0)).
		Then(clientv3.OpPut("", content)).
		Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return fmt.Errorf("config already exists")
	}
	return nil
}

func (repo *EctdRepo) UpdateConfig(ctx context.Context, serviceName, configName, content string) error {
	//key := repo.configKey(serviceName, configName)
	txnResp, err := repo.cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(""), ">", 0)).
		Then(clientv3.OpPut("", content)).
		Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return fmt.Errorf("config not found")
	}
	return nil
}

func (repo *EctdRepo) DeleteConfig(ctx context.Context, serviceName, configName string) error {
	resp, err := repo.cli.Delete(ctx, "")
	if err != nil {
		return err
	}
	if resp.Deleted == 0 {
		return fmt.Errorf("config not found")
	}
	return nil
}
