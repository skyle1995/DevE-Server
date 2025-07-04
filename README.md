# DevE 网络验证系统

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-v1.9.1-brightgreen.svg)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 项目介绍

DevE 是一个基于 Golang 开发的完整网络验证系统，主要包含用户后台核心模块，支持多应用管理，每个应用拥有独立的卡密和设备管理。系统实现了应用管理、卡密管理、系统设置等核心功能，提供安全、高效的验证服务。

## 功能特点

### 用户后台

- **仪表盘**：展示卡密信息、系统状态等统计数据
- **应用管理**：创建与配置应用、管理应用状态和密钥
- **卡密管理**：生成与批量生成卡密、管理卡密状态
- **用户设置**：登录策略、权限模板、密码策略设置
- **日志管理**：查询与导出登录日志和操作日志
- **系统设置**：基础参数配置、安全策略设置、备份与恢复

### 权限管理

系统支持三种用户角色：

- **系统管理员**：可以管理所有用户和数据，修改系统全局配置
- **普通会员**：只能管理自己创建的应用、卡密和设备
- **VIP会员**：拥有更高的应用创建上限和API调用频率

## 技术栈

- **后端开发语言**：Golang
- **前端开发技术**：layui vue admin框架（Vue）
- **后端架构**：MVC（采用cobra+gin框架）
- **数据库**：MySql或SQLite（通过配置可选）
- **配置管理**：github.com/spf13/viper
- **日志记录**：github.com/sirupsen/logrus
- **静态资源**：使用embed打包编译

## 项目结构

```
.
├── apps/               # 应用目录（按功能模块划分）
├── cmd/                # 命令行入口
├── config/             # 配置文件
├── data/               # 数据目录
├── database/           # 数据库模块
├── middleware/         # 中间件
├── public/             # 前端资源
├── server/             # 服务端
├── utils/              # 工具包
├── go.mod              # Go模块定义
├── go.sum              # 依赖校验
├── main.go             # 程序入口
└── README.md           # 项目说明
```

## 安装与使用

### 环境要求

- Go 1.23+
- MySQL 5.7+ 或 SQLite 3

### 从源码构建

1. 克隆仓库

```bash
git clone https://github.com/skyle1995/DevE-Server.git
cd DevE-Server
```

2. 安装依赖

```bash
go mod tidy
```

3. 编译项目

```bash
go build -o deve-server
```

4. 运行服务器

```bash
# 使用默认配置
./deve-server start

# 指定配置文件
./deve-server start --config=config.yaml

# 指定端口
./deve-server start --port=8080

# 启用调试模式
./deve-server start --debug
```

### 配置说明

系统默认使用 `config.yaml` 作为配置文件，可以通过命令行参数 `--config` 指定其他配置文件。配置文件包含以下主要部分：

- **服务器配置**：端口、运行模式、日志级别等
- **数据库配置**：数据库类型、连接参数等
- **安全配置**：JWT密钥、密码加密强度等

## 开发指南

### 添加新模块

1. 在 `apps` 目录下创建新的模块目录
2. 实现 controller.go、service.go 和 router.go
3. 在 server/router.go 中注册新模块的路由

### 数据库迁移

系统使用 GORM 进行数据库操作，数据库迁移在 `database/migration.go` 中定义。

## 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。

## 贡献指南

欢迎提交问题和功能请求。如果您想贡献代码，请遵循以下步骤：

1. Fork 仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request