# 快速开始

本指南将帮助您快速搭建和运行 Hoshikuzu 系统。

## 环境要求

### 后端环境

- Go 1.21 或更高版本
- MySQL 5.7 或更高版本 / PostgreSQL 12 或更高版本
- Redis 5.0 或更高版本（可选）

### 前端环境

- Node.js 18 或更高版本
- npm 或 yarn 或 pnpm

## 安装步骤

### 1. 克隆项目

```bash
git clone https://github.com/your-org/hoshikuzu.git
cd hoshikuzu
```

### 2. 后端安装

#### 2.1 安装依赖

```bash
go mod download
```

#### 2.2 配置数据库

在项目根目录创建配置文件 `config.yaml`：

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  name: hoshikuzu
  username: root
  password: your_password

server:
  port: 13000
  mode: debug
```

#### 2.3 初始化数据库

```bash
go run cmd/migrate/main.go
```

#### 2.4 生成 Ent 代码

```bash
go generate ./ent
```

#### 2.5 启动后端服务

```bash
go run cmd/server/main.go
```

后端服务将在 `http://localhost:13000` 启动。

### 3. 前端安装

#### 3.1 进入前端目录

```bash
cd ui
```

#### 3.2 安装依赖

```bash
npm install
```

#### 3.3 配置 API 地址

在 `ui/src/config/index.ts` 中配置后端 API 地址：

```typescript
export const BASE_URL = 'http://localhost:13000/api'
```

#### 3.4 启动前端服务

```bash
npm run dev
```

前端服务将在 `http://localhost:5732` 启动。

## 验证安装

### 访问管理后台

打开浏览器访问 `http://localhost:5732`，您将看到 Hoshikuzu 的管理后台界面。

### 测试 API 接口

使用以下命令测试后端 API：

```bash
curl http://localhost:13000/api/v1/health
```

预期返回：

```json
{
  "status": "ok",
  "message": "Hoshikuzu is running"
}
```

## 下一步

- 阅读 [配置文档](config.md) 了解系统配置
- 查看 [开发指南](dev/integration.md) 学习如何接入前台项目
- 探索各个功能模块的详细文档

## 常见问题

### 端口被占用

如果默认端口被占用，可以在配置文件中修改端口号：

**后端端口**：在 `config.yaml` 中修改 `server.port`

**前端端口**：在 `ui/vite.config.ts` 中修改 `server.port`

### 数据库连接失败

请检查数据库配置是否正确，确保数据库服务已启动，并且用户名密码正确。

### 依赖安装失败

尝试使用国内镜像源：

```bash
# Go 依赖
go env -w GOPROXY=https://goproxy.cn,direct

# npm 依赖
npm config set registry https://registry.npmmirror.com
```

## 获取帮助

如果您在安装过程中遇到问题，可以：

- 查看 [常见问题](faq.md) 文档
- 在 GitHub 上提交 Issue
- 联系技术支持团队
