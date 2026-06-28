package common

import _ "embed"

// EtcdConfigYAML 是公共 etcd 配置文件内容。
//
//go:embed config/etcd.yaml
var EtcdConfigYAML []byte
