package etcd

import (
	"context"
	"fmt"
	"path"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

// LoadYAML 从 etcd 指定 key 加载 YAML 配置。
func LoadYAML(ctx context.Context, client *clientv3.Client, key string, out any) error {
	resp, err := client.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return fmt.Errorf("etcd config key %s not found", key)
	}
	return yaml.Unmarshal(resp.Kvs[0].Value, out)
}

// LoadServiceYAML 按公共配置到主配置的顺序加载并合并服务配置。
func LoadServiceYAML(ctx context.Context, client *clientv3.Client, cfg Config, serviceName string, out any) error {
	cfg = cfg.withDefaults()
	names := serviceConfigNames(cfg.ConfigNames, serviceName)
	if len(names) == 0 {
		return nil
	}

	merged := &yaml.Node{}
	for i := len(names) - 1; i >= 0; i-- {
		key := configKey(cfg.ConfigPrefix, names[i])
		node, err := getYAMLNode(ctx, client, key)
		if err != nil {
			return err
		}
		merged = mergeYAMLNode(merged, node)
	}
	return merged.Decode(out)
}

func serviceConfigNames(names []string, serviceName string) []string {
	mainName := serviceName + ".yaml"
	out := make([]string, 0, len(names)+1)
	out = append(out, mainName)
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" || name == mainName {
			continue
		}
		out = append(out, name)
	}
	return out
}

func getYAMLNode(ctx context.Context, client *clientv3.Client, key string) (*yaml.Node, error) {
	resp, err := client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("etcd config key %s not found", key)
	}

	var node yaml.Node
	if err := yaml.Unmarshal(resp.Kvs[0].Value, &node); err != nil {
		return nil, err
	}
	return &node, nil
}

func configKey(prefix, configName string) string {
	prefix = "/" + strings.Trim(prefix, "/")
	configName = strings.Trim(configName, "/")
	return path.Join(prefix, configName)
}

func mergeYAMLNode(base, override *yaml.Node) *yaml.Node {
	base = unwrapDocumentNode(base)
	override = unwrapDocumentNode(override)
	if override == nil || override.Kind == 0 {
		return base
	}
	if base == nil || base.Kind == 0 || base.Kind != yaml.MappingNode || override.Kind != yaml.MappingNode {
		return override
	}

	result := cloneYAMLNode(base)
	for i := 0; i < len(override.Content); i += 2 {
		key := override.Content[i]
		value := override.Content[i+1]
		idx := mappingValueIndex(result, key.Value)
		if idx < 0 {
			result.Content = append(result.Content, cloneYAMLNode(key), cloneYAMLNode(value))
			continue
		}
		result.Content[idx] = mergeYAMLNode(result.Content[idx], value)
	}
	return result
}

func unwrapDocumentNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return node.Content[0]
	}
	return node
}

func mappingValueIndex(node *yaml.Node, key string) int {
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return i + 1
		}
	}
	return -1
}

func cloneYAMLNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	clone := *node
	clone.Content = make([]*yaml.Node, len(node.Content))
	for i, child := range node.Content {
		clone.Content[i] = cloneYAMLNode(child)
	}
	return &clone
}
