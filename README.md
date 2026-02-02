# 🐝 Hoshikuzu - 现代化内容管理系统

> 🌟 一个基于 Go 语言开发的超现代化 CMS 系统，让你的内容管理变得轻松愉快！

## 📋 项目概述

Hoshikuzu 是一个基于 **Go 1.25+** 和 **Fiber v2.52.9 框架** 开发的现代化内容管理系统（CMS），专为开发者打造的高效、灵活、易扩展的内容管理解决方案。支持内容模型管理、用户系统、支付集成、AI助手等丰富功能，前后端分离架构让你的开发体验飞起来！🚀

### ✨ 核心亮点

- 🏗️ **现代化架构**: 基于 Go 1.25+ 和 Fiber v2.52.9 高性能 Web 框架
- 💾 **强大ORM支持**: Ent ORM v0.14.5 支持 SQLite、MySQL、PostgreSQL
- 🎨 **前后端分离**: 后端提供 RESTful API，前端使用 Vue 3 + TypeScript + Naive UI
- 🔐 **完整用户系统**: OAuth2 认证，支持邮箱/手机号登录
- 🤖 **AI智能助手**: 集成 OpenAI，提供智能对话和内容生成
- 📝 **灵活内容管理**: 文章、文档、知识库、朋友圈等多类型内容管理
- 💰 **商城系统**: 商品、会员、优惠券、支付订单完整电商功能
- 📁 **文件管理**: 支持本地和 AWS S3 云存储
- 🖼️ **相册功能**: 图片管理和相册分类
- 💬 **评论系统**: 内容互动和评论管理
- 🔌 **插件系统**: 支持插件扩展，灵活定制功能

## 🚀 安装指南

### 📦 环境要求

在开始之前，请确保你的开发环境满足以下要求：

| 依赖项 | 最低版本 | 推荐版本 | 备注 |
|--------|----------|----------|------|
| Go | 1.25 | 1.25+ | 后端运行环境 |
| Node.js | 18.0.0 | 18+ | 前端开发环境 |
| SQLite | 3.x | 最新版 | 默认数据库 |
| Git | 任意版本 | 最新版 | 代码管理 |

### 🔧 快速安装

#### 1. 克隆项目代码

```bash
git clone https://github.com/shuTwT/hoshikuzu.git
cd hoshikuzu
```

#### 2. 安装后端依赖

```bash
# 下载 Go 依赖包
go mod download

# 验证依赖完整性
go mod verify
```

#### 3. 安装前端依赖

```bash
# 进入前端目录
cd ui

# 安装 npm 依赖
pnpm install

# 返回项目根目录
cd ..
```

#### 4. 数据库初始化

```bash
# 生成 Ent ORM 代码
go generate ./ent

# 数据库迁移（首次运行会自动创建表结构）
go run main.go
```

## 📖 使用说明

### 🏃‍♂️ 开发模式运行

#### 后端开发（推荐热重载）

```bash
# Windows 系统
air -c .air.windows.toml

# macOS/Linux 系统
air -c .air.mac.toml

# 或者使用 Makefile
make run
```

#### 前端开发

```bash
# 进入前端目录
cd ui

# 启动开发服务器
pnpm run dev

# 前端开发服务器运行在 http://localhost:5173
```

### 🌐 访问应用

当项目成功启动后，你可以通过以下地址访问：

- 🏠 **前端页面**: http://localhost:5379
- 🔧 **管理后台**: http://localhost:5379/console
- 📚 **API 文档**: http://localhost:13000/swagger/index.html

### ⚙️ 环境配置

#### 基础环境变量

```bash
# 数据库连接配置
DATABASE_URL=file:ent?mode=memory&cache=shared&_fk=1

# 服务端口（默认：3000）
PORT=3000

# 运行环境（dev/prod）
STAGE=dev

# JWT 密钥（生产环境务必修改）
SECRET=your-super-secret-jwt-key-here
```

#### 数据库配置示例

| 数据库类型 | 连接字符串示例 |
|------------|----------------|
| SQLite | `file:ent?mode=memory&cache=shared&_fk=1` |
| MySQL | `mysql://user:password@localhost:3306/hoshikuzu` |
| PostgreSQL | `postgres://user:password@localhost:5432/hoshikuzu` |

### 🎯 功能模块详解

#### 📝 内容管理模块 (content)
- ✅ **文章管理**: 文章发布、编辑、分类、标签管理
- ✅ **文档库**: 文档分类管理和文档详情管理
- ✅ **知识库**: 知识库内容管理
- ✅ **朋友圈**: 朋友圈动态发布和管理
- ✅ **友情链接**: 友情链接管理和应用管理
- ✅ **相册管理**: 相册分类和图片管理
- ✅ **评论系统**: 内容评论管理

#### 💰 商城管理模块 (mall)
- ✅ **商品管理**: 商品信息管理和商品分类
- ✅ **会员管理**: 会员信息和会员等级管理
- ✅ **优惠券管理**: 优惠券创建和使用记录管理
- ✅ **支付订单**: 订单管理和支付状态跟踪
- ✅ **钱包管理**: 用户钱包余额管理

#### 🔧 基础设施模块 (infra)
- ✅ **文件管理**: 文件上传下载和文件管理
- ✅ **存储策略**: 本地存储和云存储策略配置
- ✅ **定时任务**: 定时任务调度和管理
- ✅ **数据迁移**: 数据库迁移和备份
- ✅ **访问日志**: 用户访问日志记录和统计
- ✅ **插件管理**: 插件安装、卸载和配置
- ✅ **许可证管理**: 系统许可证管理

#### 👥 系统管理模块 (system)
- ✅ **用户管理**: 用户注册、登录、权限管理
- ✅ **角色管理**: 角色定义和权限分配
- ✅ **系统设置**: 系统参数配置
- ✅ **通知管理**: 系统通知推送和管理
- ✅ **API接口**: API接口管理和权限控制
- ✅ **OAuth2认证**: OAuth2授权和令牌管理
- ✅ **系统初始化**: 系统初始化向导

#### 🤖 AI助手模块 (ai)
- ✅ **智能对话**: 集成OpenAI提供智能对话功能
- ✅ **内容生成**: AI辅助内容生成和优化

### 🏗️ 项目架构

#### 后端架构
```
hoshikuzu/
├── cmd/                    # 命令行工具
│   ├── server/            # 服务器启动入口
│   └── plugindemo/        # 插件示例
├── internal/              # 内部代码
│   ├── handlers/          # HTTP处理器
│   │   ├── content/       # 内容管理模块
│   │   ├── mall/          # 商城管理模块
│   │   ├── infra/         # 基础设施模块
│   │   └── system/        # 系统管理模块
│   ├── services/          # 业务逻辑层
│   │   ├── content/
│   │   ├── mall/
│   │   ├── infra/
│   │   └── system/
│   └── router/            # 路由配置
├── pkg/                   # 公共包
│   ├── config/            # 配置管理
│   ├── domain/            # 领域模型
│   ├── utils/             # 工具函数
│   ├── cache/             # 缓存
│   └── plugin/            # 插件接口
└── ent/                   # Ent ORM 生成代码
```

#### 前端架构
```
ui/
├── src/
│   ├── views/             # 页面组件
│   │   ├── content/       # 内容管理页面
│   │   ├── mall/          # 商城管理页面
│   │   ├── infra/         # 基础设施页面
│   │   ├── system/        # 系统管理页面
│   │   └── ai/            # AI助手页面
│   ├── api/               # API接口
│   │   ├── content/
│   │   ├── mall/
│   │   ├── infra/
│   │   └── system/
│   ├── router/            # 路由配置
│   │   └── modules/       # 路由模块
│   └── components/        # 公共组件
└── package.json
```

### 📚 技术栈

#### 后端技术栈
| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.25+ | 编程语言 |
| Fiber | v2.52.9 | Web框架 |
| Ent | v0.14.5 | ORM框架 |
| JWT | v5.3.0 | 身份认证 |
| Viper | v1.21.0 | 配置管理 |
| Redis | v9.17.2 | 缓存 |
| OpenAI | v1.41.2 | AI集成 |

#### 前端技术栈
| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.x | 前端框架 |
| TypeScript | 5.x | 类型系统 |
| Vite | 5.x | 构建工具 |
| Naive UI | 2.x | UI组件库 |
| Pinia | 2.x | 状态管理 |
| Vue Router | 4.x | 路由管理 |

## 📄 许可证信息

本项目采用 **MIT 许可证** 开源协议，这意味着你可以：

- ✅ **自由使用**: 在个人和商业项目中免费使用
- ✅ **修改源码**: 根据需求修改和定制代码
- ✅ **分发软件**: 重新分发原始或修改后的版本
- ✅ **私人使用**: 在私有项目中使用

### 📋 MIT 许可证摘要

```
MIT License

Copyright (c) 2024 Hoshikuzu Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

完整许可证文本请查看 [LICENSE](LICENSE) 文件。

## 🤝 参与贡献

我们欢迎所有形式的贡献！如果你想为 Hoshikuzu 项目做出贡献，请查看我们的 [贡献指南](CONTRIBUTING.md) 了解详情。

## 📋 行为准则

我们致力于为所有人营造一个开放和友好的环境。请阅读我们的 [行为准则](CODE_OF_CONDUCT.md) 并遵守它。

## 🙏 致谢

Hoshikuzu 项目的成功离不开以下优秀开源项目的支持，在此表示衷心感谢！

| 项目名称 | 用途 | 许可证 |
|----------|------|--------|
| [Fiber](https://github.com/gofiber/fiber) | Web 框架 | MIT |
| [Ent](https://entgo.io/) | ORM 框架 | Apache 2.0 |
| [Vue.js](https://vuejs.org/) | 前端框架 | MIT |
| [Naive UI](https://www.naiveui.com/) | UI 组件库 | MIT |
| [Vite](https://vitejs.dev/) | 构建工具 | MIT |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT 认证 | MIT |
| [Viper](https://github.com/spf13/viper) | 配置管理 | MIT |
| [go-openai](https://github.com/sashabaranov/go-openai) | OpenAI SDK | MIT |
| [AWS SDK for Go](https://github.com/aws/aws-sdk-go-v2) | AWS 集成 | Apache 2.0 |
| [go-redis](https://github.com/redis/go-redis) | Redis 客户端 | BSD-2-Clause |
| [gocron](https://github.com/go-co-op/gocron) | 定时任务 | MIT |
| [go-plugin](https://github.com/hashicorp/go-plugin) | 插件系统 | MPL-2.0 |
| [Goldmark](https://github.com/yuin/goldmark) | Markdown 解析 | MIT |
| [Bluemonday](https://github.com/microcosm-cc/bluemonday) | HTML 清理 | BSD-3-Clause |

## 📞 支持与联系

如果你在使用 Hoshikuzu 过程中遇到问题，或者有任何建议，欢迎通过以下方式联系我们：

- 🐛 **问题反馈**: [提交 Issue](https://github.com/shuTwT/hoshikuzu/issues)
- 💬 **讨论交流**: [加入 Discussions](https://github.com/shuTwT/hoshikuzu/discussions)
- 📧 **邮件联系**: admin@hoshikuzu.moe
- 🌟 **给项目点星**: 如果项目对你有帮助，欢迎点亮 Star！

---

<div align="center">

**🎉 让 Hoshikuzu 成为你内容管理的最佳伙伴！**

<!-- Made with ❤️ by the Hoshikuzu Community -->

</div>