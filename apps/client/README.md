# Client 模块

## 简介

`Client` 模块提供了面向客户端应用的API接口，是连接客户端应用与服务器的桥梁。该模块主要负责处理卡密激活、设备验证、卡密换绑和解绑等功能，确保只有合法的客户端和有效的卡密才能访问系统资源。

## 功能特点

- 应用验证：验证客户端应用的合法性
- 卡密激活：激活卡密并绑定到设备
- 设备验证：验证设备是否有权限使用应用
- 卡密换绑：将卡密绑定到新设备
- 卡密解绑：解除卡密与设备的绑定关系
- 安全机制：时间戳验证、应用密钥验证、设备绑定和换绑限制

## 模块结构

```
client/
├── controller.go          # 控制器，处理HTTP请求
├── model/                 # 数据模型
│   ├── request.go         # 请求模型
│   └── response.go        # 响应模型
├── router.go              # 路由配置
├── service.go             # 业务逻辑服务
└── README.md              # 模块说明文档
```

## API 接口

### 验证应用

- **URL**: `/api/client/verify-app`
- **方法**: POST
- **认证**: 需要ClientAuthMiddleware
- **描述**: 验证客户端应用的合法性
- **请求示例**:

```json
{
  "app_key": "APP_KEY_123",
  "timestamp": 1609459200
}
```

- **响应示例**:

```json
{
  "code": 200,
  "data": {
    "success": true,
    "app_name": "测试应用",
    "version": "1.0.0"
  },
  "message": "验证成功"
}
```

### 激活卡密

- **URL**: `/api/client/activate`
- **方法**: POST
- **认证**: 需要ClientAuthMiddleware
- **描述**: 激活卡密并绑定到设备
- **请求示例**:

```json
{
  "card_no": "TEST123456",
  "card_key": "ABCDEF123456",
  "device_id": "DEVICE_UUID_123",
  "device_info": {
    "os": "Windows",
    "version": "10",
    "cpu": "Intel i7",
    "mac": "00:11:22:33:44:55"
  },
  "app_key": "APP_KEY_123",
  "timestamp": 1609459200
}
```

- **响应示例**:

```json
{
  "code": 200,
  "data": {
    "card_no": "TEST123456",
    "status": 1,
    "activated": true,
    "expire_time": "2023-02-01T00:00:00Z",
    "device_binding": true,
    "binding_info": "DEVICE_UUID_123",
    "bind_count": 1,
    "max_bind_count": 3,
    "can_rebind": true,
    "can_unbind": true,
    "message": "卡密激活成功"
  },
  "message": "激活成功"
}
```

### 验证设备

- **URL**: `/api/client/verify`
- **方法**: POST
- **认证**: 需要ClientAuthMiddleware
- **描述**: 验证设备是否有权限使用应用
- **请求示例**:

```json
{
  "card_no": "TEST123456",
  "device_id": "DEVICE_UUID_123",
  "device_info": {
    "os": "Windows",
    "version": "10"
  },
  "app_key": "APP_KEY_123",
  "timestamp": 1609459200
}
```

- **响应示例**:

```json
{
  "code": 200,
  "data": {
    "success": true,
    "card_no": "TEST123456",
    "expire_time": "2023-02-01T00:00:00Z",
    "remain_days": 30,
    "message": "设备验证通过"
  },
  "message": "验证成功"
}
```

## 使用说明

1. 客户端应用启动时，调用`verify-app`接口验证应用的合法性
2. 用户输入卡密后，调用`activate`接口激活卡密并绑定设备
3. 应用每次启动或定期调用`verify`接口验证设备权限
4. 用户需要更换设备时，调用`rebind`接口进行换绑
5. 用户需要解除绑定时，调用`unbind`接口进行解绑

## 客户端认证流程

1. **应用验证**：
   - 客户端应用启动时，发送AppKey和时间戳
   - 服务器验证AppKey的有效性和时间戳的合理性
   - 验证通过后，客户端可以继续使用

2. **卡密激活**：
   - 用户输入卡号和卡密
   - 客户端收集设备信息，生成设备ID
   - 发送卡密信息、设备ID和设备信息到服务器
   - 服务器验证卡密，并将其绑定到设备

3. **设备验证**：
   - 客户端发送卡号、设备ID和设备信息
   - 服务器验证设备是否与卡密绑定
   - 验证卡密是否有效（未过期、未禁用）

## 开发与扩展

如需扩展客户端接口模块功能，可以考虑以下方向：

1. 添加设备指纹识别，提高设备识别的准确性
2. 实现卡密自动续期功能
3. 添加客户端行为分析和异常检测
4. 实现多设备同时在线限制
5. 添加客户端远程控制功能