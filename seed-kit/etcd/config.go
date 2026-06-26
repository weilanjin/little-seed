package etcd

import "little-seed/kit/types"

type Config struct {
	Endpoints     []string       `yaml:"endpoints"`
	Username      string         `yaml:"username"`
	Password      string         `yaml:"password"`
	DialTimeout   types.Duration `yaml:"dialTimeout"`
	ConfigKey     string         `yaml:"configKey"`
	ServicePrefix string         `yaml:"servicePrefix"`
	RegisterTTL   int64          `yaml:"registerTTL"`
}

func (c Config) Enabled() bool {
	return len(c.Endpoints) > 0
}

func (c Config) ttl() int64 {
	if c.RegisterTTL > 0 {
		return c.RegisterTTL
	}
	return 10
}
