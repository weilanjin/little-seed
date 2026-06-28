package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"little-seed/kit/conf"
	"little-seed/kit/core"
	"little-seed/kit/core/hs"
	"little-seed/kit/etcd"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Etcd   etcd.Config  `yaml:"-"`
}

type ServerConfig struct {
	Name        string   `yaml:"name"`
	Addr        string   `yaml:"addr"`
	ConfigNames []string `yaml:"configNames"`
}

func main() {
	server := core.NewServer[Config]()
	server.Init(func(app *Config) error {
		cfg, err := conf.LoadYAML[Config]("config.yaml")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		if err := loadRemoteConfig(context.Background(), cfg); err != nil {
			return fmt.Errorf("failed to load remote config: %w", err)
		}
		*app = *cfg
		return nil
	})
	server.Add(func(app *Config) (core.Service, error) {
		return etcd.NewRegistryService(app.Etcd, etcd.Service{
			Name: app.Server.Name,
			Addr: app.Server.Addr,
		}), nil
	})

	server.Add(func(app *Config) (core.Service, error) {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})

		engine := hs.New(app.Server.Addr)
		engine.SetHandler(mux)
		return engine, nil
	})

	if err := server.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func loadRemoteConfig(ctx context.Context, cfg *Config) error {
	etcdCfg, err := etcd.LoadCommonConfig()
	if err != nil {
		return err
	}
	etcdCfg.ConfigNames = cfg.Server.ConfigNames
	cfg.Etcd = etcdCfg

	if !cfg.Etcd.Enabled() {
		return nil
	}

	serviceName := cfg.Server.Name
	if serviceName == "" {
		return fmt.Errorf("server.name is required")
	}

	client, err := etcd.NewClient(ctx, cfg.Etcd)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := etcd.LoadServiceYAML(ctx, client, cfg.Etcd, serviceName, cfg); err != nil {
		return err
	}
	cfg.Server.Name = serviceName
	cfg.Server.ConfigNames = etcdCfg.ConfigNames
	cfg.Etcd = etcdCfg
	return nil
}
