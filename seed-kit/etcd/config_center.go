package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

func LoadYAML(ctx context.Context, client *clientv3.Client, key string, out any) error {
	resp, err := client.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return fmt.Errorf("etcd config key %s not found", key)
	}
	return yaml.Unmarshal(resp.Kvs[0].Value, out)
}
