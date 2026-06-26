package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient(ctx context.Context, cfg Config) (*clientv3.Client, error) {
	dialTimeout := cfg.DialTimeout.TimeDuration()
	if dialTimeout <= 0 {
		dialTimeout = 5 * time.Second
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	if err := client.Sync(ctx); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}
