package main

import (
	"context"
	"log"

	"little-seed/kit/conf"
	"little-seed/kit/core"
	"little-seed/kit/etcd"
)

type Config struct {
	Name   string       `yaml:"name"`
	Server ServerConfig `yaml:"server"`
	Etcd   etcd.Config  `yaml:"etcd"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

func main() {
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if err := loadRemoteConfig(context.Background(), &cfg); err != nil {
		log.Fatalf("failed to load remote config: %v", err)
	}

	server := core.NewServer[Config]()
	server.Init(func(app *Config) error {
		*app = cfg
		return nil
	})
	server.Add(func(app *Config) (core.Service, error) {
		return etcd.NewRegistryService(app.Etcd, etcd.Service{
			Name: app.Name,
			Addr: app.Server.Addr,
		}), nil
	})

	server.Add(func(app *Config) (core.Service, error) {
		return core.NewGRPCServer(app.Server.Addr), nil
	})

	if err := server.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	if err := conf.LoadYAML(path, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func loadRemoteConfig(ctx context.Context, cfg *Config) error {
	if !cfg.Etcd.Enabled() || cfg.Etcd.ConfigKey == "" {
		return nil
	}

	client, err := etcd.NewClient(ctx, cfg.Etcd)
	if err != nil {
		return err
	}
	defer client.Close()

	return etcd.LoadYAML(ctx, client, cfg.Etcd.ConfigKey, cfg)
}
