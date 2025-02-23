# gin-api-template
这是一个基于Gin框架封装的面向API开发的框架模板，基于Go1.22.5+构建。
# 功能特性
1. 支持Logrus日志
2. 支持请求超时控制
3. 支持配置yaml文件读取
4. 支持Swagger的API文档
5. 支持基于cors的跨域组件
6. 统一的响应格式和参数校验
7. 实现基于GORM的分页构造器组件
8. 双Token（访问令牌和刷新）认证机制
9. 支持MySQL、PostgreSQL、Redis连接管理
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
# 技术栈
|  框架  |   版本   |    用途     |
| :----: | :------: | :---------: |
|  gin   | v1.10.0  |   web框架   |
|  jwt   |  v5.2.1  |  用户认证   |
|  gorm  | v1.25.12 |   ORM框架   |
|  ini   | v1.67.0  | 处理ini文件 |
| crypto | v0.33.0  |  密码加密   |
| logrus |  v1.9.3  |  日志框架   |
|  swag  | v0.23.0  | 编写API文档 |

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

# 分页构造器
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