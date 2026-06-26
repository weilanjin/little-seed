package etcd

import (
	"context"
	"encoding/json"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceEvent struct {
	Type    mvccpb.Event_EventType
	Service Service
}

func FindServices(ctx context.Context, client *clientv3.Client, prefix, name string) ([]Service, error) {
	resp, err := client.Get(ctx, servicePrefix(prefix, name), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	services := make([]Service, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var service Service
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

func WatchServices(ctx context.Context, client *clientv3.Client, prefix, name string) <-chan ServiceEvent {
	out := make(chan ServiceEvent)
	watchCh := client.Watch(ctx, servicePrefix(prefix, name), clientv3.WithPrefix(), clientv3.WithPrevKV())

	go func() {
		defer close(out)
		for resp := range watchCh {
			for _, event := range resp.Events {
				kv := event.Kv
				if event.Type == clientv3.EventTypeDelete && event.PrevKv != nil {
					kv = event.PrevKv
				}

				var service Service
				if err := json.Unmarshal(kv.Value, &service); err != nil {
					continue
				}
				out <- ServiceEvent{
					Type:    event.Type,
					Service: service,
				}
			}
		}
	}()
	return out
}

func servicePrefix(prefix, name string) string {
	return serviceKey(prefix, name, "")
}
