package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Service struct {
	Name     string            `json:"name"`
	Addr     string            `json:"addr"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type RegistryService struct {
	cfg     Config
	service Service

	client  *clientv3.Client
	leaseID clientv3.LeaseID
	cancel  context.CancelFunc
}

func NewRegistryService(cfg Config, service Service) *RegistryService {
	return &RegistryService{
		cfg:     cfg.withDefaults(),
		service: service,
	}
}

func (s *RegistryService) Start(ctx context.Context) error {
	if !s.cfg.Enabled() {
		return nil
	}

	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	client, err := NewClient(runCtx, s.cfg)
	if err != nil {
		cancel()
		return err
	}
	s.client = client

	lease, err := client.Grant(runCtx, s.cfg.ttl())
	if err != nil {
		_ = client.Close()
		cancel()
		return err
	}
	s.leaseID = lease.ID

	value, err := json.Marshal(s.service)
	if err != nil {
		_ = client.Close()
		cancel()
		return err
	}

	key := serviceKey(s.cfg.ServicePrefix, s.service.Name, s.service.Addr)
	if _, err := client.Put(runCtx, key, string(value), clientv3.WithLease(lease.ID)); err != nil {
		_ = client.Close()
		cancel()
		return err
	}

	keepAliveCh, err := client.KeepAlive(runCtx, lease.ID)
	if err != nil {
		_ = client.Close()
		cancel()
		return err
	}

	go s.keepAlive(runCtx, keepAliveCh)
	slog.Info("service registered", "name", s.service.Name, "addr", s.service.Addr, "key", key)
	return nil
}

func (s *RegistryService) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	if s.client == nil {
		return nil
	}
	if s.leaseID != 0 {
		if _, err := s.client.Revoke(ctx, s.leaseID); err != nil {
			_ = s.client.Close()
			return err
		}
	}
	return s.client.Close()
}

func (s *RegistryService) keepAlive(ctx context.Context, keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp, ok := <-keepAliveCh:
			if !ok {
				slog.Warn("service keepalive stopped", "name", s.service.Name, "addr", s.service.Addr)
				return
			}
			if resp == nil {
				continue
			}
		}
	}
}

func serviceKey(prefix, name, addr string) string {
	prefix = "/" + strings.Trim(prefix, "/")
	name = strings.Trim(name, "/")
	addr = strings.Trim(addr, "/")
	if addr == "" {
		return fmt.Sprintf("%s/%s", prefix, name)
	}
	return fmt.Sprintf("%s/%s/%s", prefix, name, addr)
}
