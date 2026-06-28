package etcd

import (
	"little-seed/common"
	"little-seed/kit/types"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Endpoints     []string       `yaml:"endpoints"`
	Username      string         `yaml:"username"`
	Password      string         `yaml:"password"`
	DialTimeout   types.Duration `yaml:"dialTimeout"`
	ConfigPrefix  string         `yaml:"configPrefix"`
	ConfigNames   []string       `yaml:"configNames"`
	ServicePrefix string         `yaml:"servicePrefix"`
	RegisterTTL   int64          `yaml:"registerTTL"`
}

// LoadCommonConfig 加载 common 模块中的公共 etcd 配置。
func LoadCommonConfig() (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(common.EtcdConfigYAML, &cfg); err != nil {
		return Config{}, err
	}
	return cfg.withDefaults(), nil
}

func (c Config) Enabled() bool {
	return len(c.withDefaults().Endpoints) > 0
}

func (c Config) ttl() int64 {
	return c.withDefaults().RegisterTTL
}

func (c Config) withDefaults() Config {
	if c.DialTimeout.TimeDuration() <= 0 {
		c.DialTimeout = "5s"
	}
	if c.RegisterTTL <= 0 {
		c.RegisterTTL = 10
	}
	return c
}
