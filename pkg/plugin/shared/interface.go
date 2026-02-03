package shared

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	pb "github.com/shuTwT/hoshikuzu/pkg/plugin/proto"
	"google.golang.org/grpc"
)

// PluginBasicInfo 核心数据结构体（数据契约），宿主和插件共享
// 注意：字段名、字段类型、字段含义必须完全一致，可选字段用指针或`omitempty`标识（JSON序列化时）
type PluginBasicInfo struct {
	Name        string `json:"name" protobuf:"bytes,1,opt,name=name"`                         // 插件唯一标识
	DisplayName string `json:"display_name" protobuf:"bytes,2,opt,name=display_name"`         // 展示名称
	Description string `json:"description,omitempty" protobuf:"bytes,3,opt,name=description"` // 可选描述
	Author      string `json:"author" protobuf:"bytes,4,opt,name=author"`                     // 插件作者
}

// PluginStore 核心接口（行为契约），宿主和插件必须共享完全一致的定义
// 注意：方法签名（参数、返回值）不能有任何差异，包括context传递、错误返回
type PluginStore interface {
	// Info 获取插件基础信息
	Info(ctx context.Context) (PluginBasicInfo, error)

	// Version 获取插件当前版本
	Version(ctx context.Context) (string, error)

	// Init 插件初始化（宿主加载时注入配置）
	Init(ctx context.Context, config map[string]string) error

	// Destroy 插件销毁（宿主关闭时释放资源）
	Destroy(ctx context.Context) error

	// Health 插件健康检查
	Health(ctx context.Context) (bool, error)

	// GetStaticResource 读取静态资源
	GetStaticResource(ctx context.Context, path string) ([]byte, string, error)
}

// --------------------------
// 第一步：定义框架要求的 Plugin 结构体（载体）
// --------------------------
// Plugin 实现 hashicorp/go-plugin 框架的 Plugin 接口
// 作用：包装我们的业务接口 PluginStore，让框架能识别并完成RPC通信
type StorePlugin struct {
	impl PluginStore
}

type StoreGRPCPlugin struct {
	plugin.Plugin
	impl PluginStore
}

// --------------------------
// 1. 实现 gRPC 协议专属方法：GRPCPlugin()
// 作用：告诉 go-plugin 框架，这是一个 gRPC 插件，使用 gRPC 协议通信
// （这是你缺失的关键方法，必须添加）
// --------------------------
func (p *StoreGRPCPlugin) GRPCPlugin() bool {
	return true
}

// --------------------------
// 2. 实现 gRPC 协议的 Server() 方法（不变，之前的实现正确）
// --------------------------
func (p *StorePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GRPCServer{}, nil
}

// --------------------------
// 3. 修正 Client() 方法注释（明确是 gRPC 协议的 Client 方法）
// （你的实现本身正确，只是框架没识别到 gRPC 协议）
// --------------------------
func (p *StorePlugin) Client(broker *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// --------------------------
// 实现 GRPCPlugin 接口的 GRPCServer 方法（插件侧运行）
// 作用：将 PluginStore 业务接口的实现注册到 gRPC 服务端
// --------------------------
func (p *StoreGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	// 1. 创建 gRPC 服务端实例（封装业务实现）
	grpcServer := &GRPCServer{}
	// 2. 注册到 gRPC 服务端（使用 pb 包生成的注册方法）
	pb.RegisterPluginStoreServiceServer(s, grpcServer)
	// 3. 无错误返回 nil
	return nil
}

// --------------------------
// 实现 GRPCPlugin 接口的 GRPCClient 方法（宿主侧运行）
// 作用：创建 gRPC 客户端，封装为 PluginStore 业务接口返回给宿主
// --------------------------
func (p *StoreGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	// 1. 创建 pb 包生成的 gRPC 客户端
	grpcClient := pb.NewPluginStoreServiceClient(conn)
	// 2. 封装为业务接口 GRPCClient，返回给宿主
	return &GRPCClient{client: grpcClient}, nil
}

// ProtocolVersion：框架要求，返回通信协议版本（需与握手配置一致）
// func (p *Plugin) ProtocolVersion() int {
// 	return consts.ProtocolVersion // 公共包中定义的常量（之前的 consts.go 中）
// }

// --------------------------
// 第三步：实现 gRPC 服务端（插件侧，接收宿主请求并转发给 PluginStore 实现）
// --------------------------
// GRPCServer 实现 pb 包中生成的 PluginStoreServiceServer 接口（gRPC服务端）
// 作用：将 gRPC 请求映射到我们的业务接口 PluginStore 的具体实现
type GRPCServer struct {
	// 嵌入 pb 包的未实现结构体（避免手动实现所有方法，兼容gRPC版本）
	pb.UnimplementedPluginStoreServiceServer
	// 持有 PluginStore 业务接口的具体实现（插件侧运行时，会注入实际的插件实例）
	impl PluginStore
}

// GetStaticResource：实现 gRPC 服务端的 GetStaticResource 方法，处理宿主的静态资源读取请求
func (s *GRPCServer) GetStaticResource(ctx context.Context, req *pb.StaticResourceRequestMsg) (*pb.StaticResourceResponseMsg, error) {
	// 调用插件实现的 GetStaticResource 方法
	content, contentType, err := s.impl.GetStaticResource(ctx, req.Path)
	if err != nil {
		return nil, err
	}
	// 转换为 gRPC 响应
	return &pb.StaticResourceResponseMsg{
		Content:     content,
		ContentType: contentType,
	}, nil
}

// 确保 GRPCClient 完全实现 PluginStore 接口（编译期校验）
var _ PluginStore = &GRPCClient{}
