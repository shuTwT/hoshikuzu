# 配置

Hoshikuzu 的配置文件采用 YAML 格式，位于项目根目录的 `config.yaml` 文件中。

## 配置文件结构

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

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

log:
  level: info
  format: json
  output: stdout

jwt:
  secret: your_jwt_secret
  expire: 86400

upload:
  path: ./uploads
  max_size: 10485760
  allowed_types:
    - image/jpeg
    - image/png
    - image/gif

payment:
  alipay:
    app_id: your_alipay_app_id
    private_key: your_alipay_private_key
    public_key: your_alipay_public_key
  wechat:
    app_id: your_wechat_app_id
    mch_id: your_wechat_mch_id
    api_key: your_wechat_api_key
  epay:
    url: your_epay_url
    pid: your_epay_pid
    key: your_epay_key

ai:
  api_key: your_ai_api_key
  model: gpt-4
  max_tokens: 2000
```

## 配置项说明

### 数据库配置 (database)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| driver | string | 数据库驱动类型 (mysql/postgres) | mysql |
| host | string | 数据库主机地址 | localhost |
| port | int | 数据库端口 | 3306 |
| name | string | 数据库名称 | gobee |
| username | string | 数据库用户名 | root |
| password | string | 数据库密码 | - |

### 服务器配置 (server)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| port | int | 服务器监听端口 | 13000 |
| mode | string | 运行模式 (debug/release) | debug |

### Redis 配置 (redis)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| host | string | Redis 主机地址 | localhost |
| port | int | Redis 端口 | 6379 |
| password | string | Redis 密码 | - |
| db | int | Redis 数据库编号 | 0 |

### 日志配置 (log)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| level | string | 日志级别 (debug/info/warn/error) | info |
| format | string | 日志格式 (json/text) | json |
| output | string | 日志输出 (stdout/file) | stdout |

### JWT 配置 (jwt)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| secret | string | JWT 签名密钥 | - |
| expire | int | JWT 过期时间（秒） | 86400 |

### 文件上传配置 (upload)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| path | string | 文件上传路径 | ./uploads |
| max_size | int | 最大文件大小（字节） | 10485760 |
| allowed_types | array | 允许的文件类型 | - |

### 支付配置 (payment)

#### 支付宝配置 (alipay)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| app_id | string | 支付宝应用 ID | - |
| private_key | string | 支付宝应用私钥 | - |
| public_key | string | 支付宝公钥 | - |

#### 微信配置 (wechat)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| app_id | string | 微信应用 ID | - |
| mch_id | string | 微信商户号 | - |
| api_key | string | 微信 API 密钥 | - |

#### 易支付配置 (epay)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| url | string | 易支付接口地址 | - |
| pid | string | 易支付商户 ID | - |
| key | string | 易支付密钥 | - |

### AI 配置 (ai)

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| api_key | string | AI 服务 API 密钥 | - |
| model | string | AI 模型名称 | gpt-4 |
| max_tokens | int | 最大 token 数 | 2000 |

## 环境变量

除了配置文件，Gobee 也支持通过环境变量进行配置：

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=gobee
export DB_USERNAME=root
export DB_PASSWORD=your_password
export SERVER_PORT=13000
export JWT_SECRET=your_jwt_secret
```

## 配置最佳实践

### 生产环境配置

1. **使用强密码**：确保数据库、JWT、支付等敏感配置使用强密码
2. **修改默认端口**：在生产环境中修改默认端口
3. **启用 HTTPS**：在生产环境中使用 HTTPS
4. **日志级别**：生产环境建议使用 `info` 或 `warn` 级别
5. **文件上传**：限制文件上传大小和类型

### 开发环境配置

1. **使用调试模式**：设置 `server.mode` 为 `debug`
2. **详细日志**：设置 `log.level` 为 `debug`
3. **本地数据库**：使用本地数据库进行开发

### 安全建议

1. **不要提交配置文件**：将 `config.yaml` 添加到 `.gitignore`
2. **使用环境变量**：敏感信息通过环境变量配置
3. **定期更换密钥**：定期更换 JWT 密钥和 API 密钥
4. **备份配置**：定期备份配置文件

## 配置验证

启动服务时，Gobee 会自动验证配置文件的正确性。如果配置有误，服务将无法启动并显示错误信息。

## 相关文档

- [快速开始](quickstart.md) - 快速搭建 Gobee 系统
- [开发指南](dev/integration.md) - 前台项目接入指南
