package repo

import "little-seed/kit/etcd"

type Data struct {
	etcdCfg    etcd.Config
	configRoot string
	logDir     string
}

func NewData(etcdCfg etcd.Config) *Data {
	configRoot := etcdCfg.ConfigKey
	if configRoot == "" {
		configRoot = "/little-seed/configs"
	}

	return &Data{
		etcdCfg:    etcdCfg,
		configRoot: configRoot,
		logDir:     "logs",
	}
}
