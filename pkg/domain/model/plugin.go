package model

import (
	"mime/multipart"
	"time"
)

type PluginConfig struct {
	Name             string   `yaml:"name" validate:"required"`
	Key              string   `yaml:"key" validate:"required"`
	Version          string   `yaml:"version" validate:"required"`
	Description      string   `yaml:"description,omitempty"`
	ProtocolVersion  string   `yaml:"protocol_version,omitempty"`
	MagicCookieKey   string   `yaml:"magic_cookie_key,omitempty"`
	MagicCookieValue string   `yaml:"magic_cookie_value" validate:"required"`
	Config           string   `yaml:"config,omitempty"`
	Dependencies     []string `yaml:"dependencies,omitempty"`
	AutoStart        bool     `yaml:"auto_start,omitempty"`
}

type CreatePluginReq struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
}

type PluginResp struct {
	ID               int        `json:"id"`
	CreatedAt        LocalTime  `json:"created_at"`
	UpdatedAt        LocalTime  `json:"updated_at"`
	Key              string     `json:"key"`
	Name             string     `json:"name"`
	Version          string     `json:"version"`
	Description      string     `json:"description"`
	BinPath          string     `json:"bin_path"`
	ProtocolVersion  string     `json:"protocol_version"`
	MagicCookieKey   string     `json:"magic_cookie_key"`
	MagicCookieValue string     `json:"magic_cookie_value"`
	Dependencies     []string   `json:"dependencies"`
	Config           string     `json:"config"`
	Enabled          bool       `json:"enabled"`
	AutoStart        bool       `json:"auto_start"`
	Status           string     `json:"status"`
	LastError        string     `json:"last_error"`
	LastStartedAt    *LocalTime `json:"last_started_at"`
	LastStoppedAt    *LocalTime `json:"last_stopped_at"`
}

// PluginRegisterReq 插件注册请求结构体
type PluginRegisterReq struct {
	Name        string            `json:"name" validate:"required"`
	Version     string            `json:"version" validate:"required"`
	GrpcAddress string            `json:"grpc_address" validate:"required"`
	Status      string            `json:"status" validate:"required"`
	StartTime   *time.Time        `json:"start_time"`
	Metadata    map[string]string `json:"metadata"`
}

// PluginHeartbeatReq 插件心跳请求结构体
type PluginHeartbeatReq struct {
	Name   string `json:"name" validate:"required"`
	Status string `json:"status,omitempty"`
}

// PluginHeartbeatResp 插件心跳响应结构体
type PluginHeartbeatResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
