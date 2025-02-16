# gin-api-template
这是一个基于Gin框架封装的面向API开发的框架模板，基于Go1.22.5+构建。
# 功能特性
1. 支持Logrus日志
2. 支持请求超时控制
3. 支持配置文件读取
4. 支持Swagger的API文档
5. 统一的响应格式和参数校验
6. 双Token（访问令牌和刷新）认证机制
7. 支持MySQL、PostgreSQL、Redis连接管理
# 目录结构
```
.
├── api
│   ├── handler     # 请求处理器
│   ├── middleware  # 中间件
│   └── route       # 路由配置
├── bootstrap       # 应用启动初始化
├── cmd             # 主程序入口
├── config          # 配置文件
├── docs            # swagger文档存放目录
├── domain
│   ├── dto         # 请求处理器的请求和响应结构体
│   ├── entity      # 数据库实体对应的结构体
│   └── result      # 通用的响应结构体
├── internal
│   ├── redisutil   # Redis工具类
│   └── tokenutil   # token工具类
├── logs            # 存放日志目录
├── repo            # 服务层
├── service         # 持久层
└── util            # 工具类存放目录

```
# 技术栈
|  框架  |   版本   |    用途     |
| :----: | :------: | :---------: |
|  gin   | v1.10.0  |   web框架   |
|  jwt   |  v5.2.1  |  用户认证   |
|  gorm  | v1.25.12 |   ORM框架   |
|  ini   | v1.67.0  | 处理ini文件 |
| crypto | v0.33.0  |  密码加密   |
| logrus |  v1.9.3  |  日志框架   |

# 运行项目
## 安装项目依赖
```
go mod download
go mod tidy
```
> 上诉两条命令均可用来安装项目依赖，但`tidy`命令会修改`go.mod`文件，并且将会移除未使用的项目依赖

## 启动项目
```
# 使用go启动项目
go run cmd/main.go
# 使用make启动项目
make run
```
# Makefile 命令
项目提供了一下make命令用于简化操作：
```
make build      # 构建项目，生成可执行文件
make run            # 运行项目
make clean          # 清理构建文件
make clean-logs     # 清理日志文件
make swagger        # 用户生成Swagger文档
```
> 执行`make`命令默认执行`make run`