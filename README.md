# gin-api-template
这是一个基于Gin框架封装的面向API开发的框架模板，基于Go1.22.5+构建。
# ⚙️ 功能特性
1. 实现完整日志记录功能
2. 支持请求超时控制
3. 支持配置yaml文件读取
4. 支持Swagger的API文档
5. 支持基于cors的跨域组件
6. 统一的响应格式和参数校验
7. 实现基于GORM的分页构造器组件
8. 双Token（访问令牌和刷新）认证机制
9. 支持MySQL、PostgreSQL、Redis连接管理
10. 请求追踪（Trace ID）
11. 统一错误处理
12. 限流（内存/Redis）
13. RBAC权限管理系统
    - 基于角色的访问控制（RBAC）
    - 支持API路径和业务实体权限
    - 树形组织结构管理
    - 读写权限分离
    - 上级可见下级数据
    - 通配符权限匹配
    - 审计日志记录
    - Redis缓存支持（可选）
# 📂 目录结构
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
├── global			# 存放全局变量
├── internal
│   ├── redisutil   # Redis工具类
│   └── tokenutil   # token工具类
├── logs            # 存放日志目录
├── pkg				# 存放项目的可重用包和库代码
│   ├── lib			# 包含数据库、Redis、日志
│   │   └── logger  # 日志服务
│   └── pagination  # 基于gorm的分页构建器组件
├── repo            # 服务层
├── service         # 持久层
└── util            # 工具类存放目录

```
# 🧰 技术栈
|  框架  |   版本   |    用途     |
| :----: | :------: | :---------: |
|  gin   | v1.10.0  |   web框架   |
|  jwt   |  v5.2.1  |  用户认证   |
|  gorm  | v1.25.12 |   ORM框架   |
|  ini   | v1.67.0  | 处理ini文件 |
| crypto | v0.33.0  |  密码加密   |
| zap |  v1.9.3  |  日志框架   |
|  swag  | v0.23.0  | 编写API文档 |

# 🚀 运行项目
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
> 想要执行`make`命令需要安装`GNU Make`工具

## 配置文件

项目使用 `config/config.yml` 作为主配置文件，同时支持通过 `.env` 文件覆盖配置。

### .env 文件

复制 `.env.example` 为 `.env` 并填入实际值，可以覆盖 `config.yml` 中的配置：

```bash
# 数据库配置
DB_PASSWORD=your_actual_password

# Token 密钥
ACCESS_TOKEN_SECRET=your_secret_key
REFRESH_TOKEN_SECRET=your_refresh_secret
```

优先级：**环境变量 / .env > config.yml**

这样可以将敏感信息（密码、密钥）放在本地 `.env` 中，避免提交到版本控制。

### create-admin 环境变量

创建管理员时可使用环境变量自动填充：

```bash
ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD=your_password go run cmd/rbaccli/main.go create-admin
```

# 🔐 RBAC 权限管理

## 初始化系统管理员

首次运行项目前，创建系统管理员：

```bash
make create-admin
```

按照提示输入管理员邮箱和密码。

## 权限模型

### 数据模型
- **Resource（资源）**：API路径或业务实体，支持通配符（如 `/users/*`）
- **Role（角色）**：角色定义，可绑定资源和组织范围
- **OrgUnit（组织）**：树形组织结构，支持任意层级
- **UserRole（用户角色）**：用户绑定角色，指定生效组织
- **RoleResource（角色资源）**：角色绑定资源，默认读权限，写权限需单独配置
- **RoleOrgScope（角色组织范围）**：角色可访问的组织节点，支持包含子级
- **AuditLog（审计日志）**：记录所有权限变更操作

### 权限检查
- **API级别**：基于HTTP方法和URL路径的权限验证
- **实体级别**：基于业务实体的读写权限验证
- **组织范围**：用户只能看到其所在组织及子组织的数据

## 使用示例

### 创建组织节点

```bash
POST /org-units
{
  "name": "技术部",
  "parent_id": 1
}
```

### 创建角色并绑定权限

```bash
# 创建角色
POST /roles
{
  "name": "技术管理员",
  "description": "技术部管理员角色"
}

# 绑定资源（默认读权限）
POST /roles/1/resources
{
  "role_id": 1,
  "resource_id": 5,
  "is_write": false
}

# 绑定写权限
POST /roles/1/resources
{
  "role_id": 1,
  "resource_id": 5,
  "is_write": true
}
```

### 分配角色给用户

```bash
POST /users/1/roles
{
  "role_id": 1,
  "org_unit_id": 2
}
```

### 获取用户权限

```bash
GET /user/permissions
```

返回用户的所有权限和组织范围，供前端控制UI显示。

## 中间件使用

### API级别权限检查

```go
router.GET("/users", 
    middleware.JwtAuthMiddleware(),
    rbacMiddleware.CheckPermission("/users"),
    userHandler.ListUsers)
```

### 实体级别权限检查

```go
router.DELETE("/users/:id",
    middleware.JwtAuthMiddleware(),
    rbacMiddleware.CheckEntityPermission("user", "delete"),
    userHandler.DeleteUser)
```

## 组织树可见性

- **上级可见下级**：父节点组织的用户可以看到所有子组织的数据
- **下级不可见上级**：子节点组织的用户不能看到父组织的数据
- **同级不可见**：同一级别的组织用户不能互相看到对方的数据

## 审计日志

支持多种查询方式：

```bash
# 按操作者查询
GET /audit-logs?operator_id=1

# 按目标查询
GET /audit-logs/target?target_type=role&target_id=1

# 按时间范围查询
GET /audit-logs/time?start_time=2024-01-01&end_time=2024-01-31
```

## 权限缓存

当启用 Redis 时，用户权限会被缓存，提升性能：

```yaml
redis:
  Enabled: true
```

权限变更时自动清除缓存。

## Makefile 命令
项目提供了一下make命令用于简化操作：
```
make build      # 构建项目，生成可执行文件
make run            # 运行项目
make clean          # 清理构建文件
make clean-logs     # 清理日志文件
make swagger        # 用户生成Swagger文档
make create-admin   # 创建系统管理员（首次运行）
```
> 执行`make`命令默认执行`make run`

# 📄 分页构造器
## 使用方法
```go
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Age       uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Products  []Product      `gorm:"foreignKey:UserID"`
}
type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"size:100;not null"`
	Description string  `gorm:"size:255"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	UserID      uint    `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```
>基于`User`和`Product`表结构进行分页查询
```go
// 单表分页
pagination, err := pagination.NewPaginationBuilder(db).
		Model(&User{}).
		SetPage(1).
		SetPageSize(10).
		Build(&User)

// 条件查询的单表分页
pagination, err := pagination.NewPaginationBuilder(db).
		Model(&User{}).
		SetPage(1).
		SetPageSize(10).
		Where("age > ?", 25).
		Build(&User)

// 多个查询条件加排序
pagination, err := pagination.NewPaginationBuilder(db).
		Model(&User{}).
		SetPage(1).
		SetPageSize(10).
		Where("age BETWEEN ? AND ?", 20, 30).
		Where("name LIKE ?", "%a%").
		OrderBy("age ASC").
		Build(&User)

// 关联数据
pagination, err := pagination.NewPaginationBuilder(db).
		Model(&User{}).
		SetPage(1).
		SetPageSize(2).
		Preload("Products").
		Build(&User)

// 多表
pagination, err := pagination.NewPaginationBuilder(db).
		Model(&User{}).
		SetPage(1).
		SetPageSize(10).
		Join("JOIN product ON user.id = product.user_id").
		Where("product.price > ?", 200).
		GroupBy("user.id"). // To avoid duplicate users
		Preload("Products").
		Build(&User)
```
## 分页构造器的方法
1. SetPage：设置当前查询的页码
2. SetPageSize：设置每页数量
3. Model：设置模型（数据库实体对应的结构体）
4. Preload 添加预加载关系（关联关系）
5. Join：添加连接查询
6. Select：设置查询字段
7. Where：添加查询条件
8. OrderBy：添加排序
9. GroupBy：添加分组
10. Having：添加HAVING条件
11. Build：构建分页

# 🌐 前端项目 (web/)

基于 Vue 3 的后台管理系统前端，使用 Vite 构建。

## 技术栈

| 框架 | 版本 | 用途 |
| :--- | :--- | :--- |
| Vue 3 | ^3.5 | 渐进式框架 |
| TypeScript | ^5.0 | 类型安全 |
| Naive UI | ^2.40 | UI 组件库 |
| Pinia | ^3.0 | 状态管理 |
| Vue Router | ^4.0 | 路由管理 |
| Axios | ^1.7 | HTTP 客户端 |
| UnoCSS | ^0.65 | 原子化 CSS |

## 项目结构

```
web/
├── src/
│   ├── api/              # API 接口封装
│   ├── types/           # TypeScript 类型定义
│   ├── stores/          # Pinia 状态管理
│   │   ├── auth.ts      # 认证状态 (Token)
│   │   └── permission.ts # 权限状态 (菜单/权限)
│   ├── router/          # 路由配置
│   ├── utils/           # 工具函数
│   │   ├── http.ts     # Axios 封装 + Token自动刷新
│   │   └── auth.ts     # JWT 解码
│   ├── views/           # 页面组件
│   │   ├── login/       # 登录页
│   │   ├── dashboard/    # 仪表盘
│   │   ├── roles/      # 角色管理
│   │   ├── menus/      # 菜单管理
│   │   ├── orgs/       # 组织管理
│   │   ├── resources/  # 资源管理
│   │   ├── audit-logs/ # 审计日志
│   │   ├── common/     # 通用页面
│   │   └── error/     # 错误页面
│   ├── layouts/         # 布局组件
│   │   └── DefaultLayout.vue # 默认布局 (侧边栏 + 顶栏)
│   ├── directives/     # 自定义指令
│   │   └── permission.ts # v-permission 按钮权限指令
│   └── styles/         # 样式文件
└── vite.config.ts      # Vite 配置
```

## 前端运行

```bash
# 安装依赖
cd web
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
```

## 主要功能

- **登录认证**：邮箱/密码/验证码登录，JWT Token 管理，支持 Token 自动刷新
- **动态菜单**：根据后端返回的菜单树动态生成侧边栏
- **权限控制**：按钮级权限控制（`v-permission` 指令），基于 RBAC 模型
- **页面管理**：角色管理、菜单管理、组织管理、资源管理、审计日志
- **响应式**：支持侧边栏折叠

## 认证流程

1. 用户输入邮箱、密码和验证码登录
2. 后端返回 `accessToken` 和 `refreshToken`
3. `accessToken` 存储于 localStorage，用于 API 请求认证
4. Token 过期时，Axios 拦截器自动使用 `refreshToken` 刷新
5. 刷新失败时，自动跳转登录页

## 权限指令

```vue
<n-button v-permission="['/users', 'POST']">新增用户</n-button>
```

当用户不具备指定资源权限时，按钮自动隐藏。
