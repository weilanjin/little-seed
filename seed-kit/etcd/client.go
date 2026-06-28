package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient(ctx context.Context, cfg Config) (*clientv3.Client, error) {
	cfg = cfg.withDefaults()

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DialTimeout: cfg.DialTimeout.TimeDuration(),
	})
	if err != nil {
		return nil, fmt.Errorf("create etcd client: %w", err)
	}

	if err := client.Sync(ctx); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("sync etcd client: %w", err)
	}
	return client, nil
}
