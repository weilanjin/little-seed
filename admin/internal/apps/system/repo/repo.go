package repo

import (
	"context"
	"little-seed/kit/etcd"
)

type Data struct {
	EtcdCli *EctdRepo

	logDir string
}

func NewData(etcdCfg etcd.Config) *Data {
	return &Data{
		EtcdCli: NewEctdClient(context.Background(), etcdCfg),
		logDir:  "logs",
	}
}
