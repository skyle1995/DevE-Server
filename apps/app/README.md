# App 模块

## 简介

`App` 模块负责系统中应用的创建、查询、更新和删除等核心功能，是整个系统的基础模块之一。该模块允许用户创建和管理自己的应用，设置应用参数，查看应用日志等。

## 功能特点

- 应用管理：创建、查询、更新和删除应用
- 密钥管理：生成和重新生成应用密钥
- 日志管理：查看和清理应用日志
- 支持多种计费模式：时长计费和点数计费
- 提供试用功能配置

## 模块结构

```
app/
├── controller.go          # 控制器，处理HTTP请求
├── model/                 # 数据模型
│   ├── request.go         # 请求模型
│   └── response.go        # 响应模型
├── router.go              # 路由配置
├── service.go             # 业务逻辑服务
└── README.md              # 模块说明文档
```

## API 接口

### 创建应用

- **URL**: `/api/apps`
- **方法**: POST
- **认证**: 需要JWT令牌
- **描述**: 创建新的应用
- **请求示例**:

```json
{
  "name": "测试应用",
  "description": "这是一个测试应用",
  "version": "1.0.0",
  "download_url": "https://example.com/download",
  "trial_time": 7,
  "trial_count": 3
}
```

- **响应示例**:

```json
{
  "code": 200,
  "data": {
    "id": 1,
    "name": "测试应用",
    "description": "这是一个测试应用",
    "app_key": "app_123456789",
    "app_secret": "secret_987654321",
    "status": 1,
    "version": "1.0.0",
    "download_url": "https://example.com/download",
    "trial_time": 7,
    "trial_count": 3,
    "user_id": 1,
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  },
  "message": "应用创建成功"
}
```

### 获取应用列表

- **URL**: `/api/apps`
- **方法**: GET
- **认证**: 需要JWT令牌
- **描述**: 获取当前用户的应用列表
- **查询参数**:
  - `page`: 页码，默认为1
  - `page_size`: 每页记录数，默认为10
- **响应示例**:

```json
{
  "code": 200,
  "data": {
    "list": [...],  // 应用列表
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total": 5
    }
  },
  "message": "获取成功"
}
```

### 获取应用详情

- **URL**: `/api/apps/:id`
- **方法**: GET
- **认证**: 需要JWT令牌
- **描述**: 获取指定应用的详细信息
- **响应示例**:

```json
{
  "code": 200,
  "data": {
    "id": 1,
    "name": "测试应用",
    "description": "这是一个测试应用",
    "app_key": "app_123456789",
    "status": 1,
    "version": "1.0.0",
    "download_url": "https://example.com/download",
    "billing_mode": 0,
    "allow_trial": 1,
    "trial_amount": 7,
    "user_id": 1,
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  },
  "message": "获取成功"
}
```

## 使用说明

1. 用户登录后可以创建自己的应用
2. 创建应用时，系统会自动生成应用密钥（AppKey和AppSecret）
3. 应用密钥仅在创建和重新生成时返回，请妥善保存
4. 用户可以查看、更新和删除自己的应用
5. 应用可以配置试用功能，支持时长试用或点数试用

## 应用状态说明

- **启用 (1)**: 应用正常运行
- **禁用 (0)**: 应用被暂时禁用，无法使用

## 计费模式说明

- **时长计费 (0)**: 按照使用时长计费，试用额度为小时数
- **点数计费 (1)**: 按照消耗点数计费，试用额度为点数

## 开发与扩展

如需扩展应用模块功能，可以考虑以下方向：

1. 添加应用统计功能，如用户数、活跃度等
2. 实现应用版本管理和更新通知
3. 添加应用分类和标签功能
4. 实现应用审核流程
5. 添加应用评分和评论功能