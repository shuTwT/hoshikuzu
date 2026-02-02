package shared

import (
	"context"
	"net/rpc"

	pb "github.com/shuTwT/hoshikuzu/pkg/plugin/proto"
)

type RPCClient struct{ client *rpc.Client }

type RPCServer struct {
	// 嵌入 pb 包的未实现结构体（避免手动实现所有方法，兼容gRPC版本）
	pb.UnimplementedPluginStoreServiceServer
	// 持有 PluginStore 业务接口的具体实现（插件侧运行时，会注入实际的插件实例）
	Impl PluginStore
}

// 给 GRPCServer 注入 PluginStore 实现的方法（插件侧使用）
func (s *GRPCServer) SetImpl(impl PluginStore) {
	s.impl = impl
}

// Info：实现 gRPC 服务端的 Info 方法，转发给业务接口 PluginStore 的 Info 方法
func (s *GRPCServer) Info(ctx context.Context, req *pb.EmptyMsg) (*pb.PluginBasicInfoMsg, error) {
	// 1. 调用插件侧的 PluginStore.Info 业务方法
	info, err := s.impl.Info(ctx)
	if err != nil {
		return nil, err
	}
	// 2. 转换为 pb 包的消息结构体（gRPC 传输格式）
	return &pb.PluginBasicInfoMsg{
		Name:        info.Name,
		DisplayName: info.DisplayName,
		Description: info.Description,
		Author:      info.Author,
	}, nil
}

// Version：实现 gRPC 服务端的 Version 方法
func (s *GRPCServer) Version(ctx context.Context, req *pb.EmptyMsg) (*pb.VersionMsg, error) {
	version, err := s.impl.Version(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.VersionMsg{Version: version}, nil
}

// Init：实现 gRPC 服务端的 Init 方法
func (s *GRPCServer) Init(ctx context.Context, req *pb.InitMsg) (*pb.EmptyMsg, error) {
	err := s.impl.Init(ctx, req.Config)
	return &pb.EmptyMsg{}, err
}

// Destroy：实现 gRPC 服务端的 Destroy 方法
func (s *GRPCServer) Destroy(ctx context.Context, req *pb.EmptyMsg) (*pb.EmptyMsg, error) {
	err := s.impl.Destroy(ctx)
	return &pb.EmptyMsg{}, err
}

// Health：实现 gRPC 服务端的 Health 方法
func (s *GRPCServer) Health(ctx context.Context, req *pb.EmptyMsg) (*pb.HealthMsg, error) {
	healthy, err := s.impl.Health(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.HealthMsg{IsHealthy: healthy}, nil
}
