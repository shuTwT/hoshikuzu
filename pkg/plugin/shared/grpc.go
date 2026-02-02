package shared

import (
	"context"

	pb "github.com/shuTwT/hoshikuzu/pkg/plugin/proto"
)

type GRPCClient struct{ client pb.PluginStoreServiceClient }

// Info：实现业务接口 PluginStore 的 Info 方法，发送 gRPC 请求给插件端
func (c *GRPCClient) Info(ctx context.Context) (PluginBasicInfo, error) {
	// 1. 发送 gRPC 请求
	resp, err := c.client.Info(ctx, &pb.EmptyMsg{})
	if err != nil {
		return PluginBasicInfo{}, err
	}
	// 2. 转换为业务结构体（宿主侧可直接使用）
	return PluginBasicInfo{
		Name:        resp.Name,
		DisplayName: resp.DisplayName,
		Description: resp.Description,
		Author:      resp.Author,
	}, nil
}

// Version：实现业务接口 PluginStore 的 Version 方法
func (c *GRPCClient) Version(ctx context.Context) (string, error) {
	resp, err := c.client.Version(ctx, &pb.EmptyMsg{})
	if err != nil {
		return "", err
	}
	return resp.Version, nil
}

// Init：实现业务接口 PluginStore 的 Init 方法
func (c *GRPCClient) Init(ctx context.Context, config map[string]string) error {
	_, err := c.client.Init(ctx, &pb.InitMsg{Config: config})
	return err
}

// Destroy：实现业务接口 PluginStore 的 Destroy 方法
func (c *GRPCClient) Destroy(ctx context.Context) error {
	_, err := c.client.Destroy(ctx, &pb.EmptyMsg{})
	return err
}

// Health：实现业务接口 PluginStore 的 Health 方法
func (c *GRPCClient) Health(ctx context.Context) (bool, error) {
	resp, err := c.client.Health(ctx, &pb.EmptyMsg{})
	if err != nil {
		return false, err
	}
	return resp.IsHealthy, nil
}
