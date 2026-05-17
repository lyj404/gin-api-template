# AGENTS.md

## High-Signal Development Commands

### Core
```bash
rtk make build          # Build binary to bin/
rtk make run            # Run from cmd/ (entry point)
rtk make swagger        # Update Swagger docs (swag init -g cmd/main.go)
rtk wire                # Re-generate dependency injection (run in cmd/)
```

### Data & RBAC Setup (Critical Sequence)
```bash
rtk make create-admin   # Create superuser (interactive, auto-seeds resources/menus/dicts)
rtk make seed           # Re-initialize all default data (resources + menus + dicts)
rtk make seed-resources # Only sync API/Entity definitions to DB
rtk make seed-menus     # Only sync default side-menu structure
rtk make seed-dict      # Only sync system dictionaries (status, types, etc.)
```

### Backend Verification
```bash
rtk go test ./...       # Run Go tests
rtk go fmt ./...        # Standard formatting
rtk go vet ./...        # Static analysis
```

### Frontend (web/)
```bash
cd web && rtk npm run dev    # Start dev server
cd web && rtk npm run build  # Production build
```

## Architecture Signal

- **Dependency Injection**: Uses `google/wire`. If you add new services/repositories, you **MUST** update `cmd/wire.go` and run `wire` inside `cmd/`.
- **RBAC Enforcement**: Permissions are checked in `api/middleware/rbac_middleware.go`. Resources are defined in `cmd/rbaccli/main.go` and synced via `make seed-resources`.
- **ID System**: Uses Snowflake IDs (`internal/idgen/snowflake.go`). Database IDs are `uint64`.
- **Domain Boundaries**: 
  - `domain/`: Interfaces and DTOs (Source of truth for contracts).
  - `service/`: Business logic.
  - `repository/`: Data access.
- **Frontend State**: Pinia stores in `web/src/stores/` handle auth, permissions, and themes.

## Important Constraints

- **Migrations**: No dedicated migration tool; GORM AutoMigrate is used in `bootstrap/boot.go`.
- **ID Precision**: Snowflake IDs (uint64) require `json:",string"` tags in DTOs to prevent precision loss in JavaScript.
- **API Documentation**: Always update Swagger comments in `api/handler/` after changing endpoints, then run `make swagger`.
- **Context**: Always propagate `context.Context` from handler -> service -> repository for timeout/traceability.
- **Validation**: Use `binding` tags in DTOs for Gin's built-in validator.
- **RBAC Sync**: Super admin permissions are reset/synced every time `make create-admin` or `make seed-resources` is run.

## Workflow Quirks
- **First Run**: `make create-admin` is mandatory to access the UI.
- **Environment**: `.env` overrides `config/config.yml`. Sensitive keys like `ACCESS_TOKEN_SECRET` belong in `.env`.
- **Audit Logging**: Automatic for RBAC changes, manual for business logic via `AuditLogService`.
- **Tree Structures**: Org Units and Menus use parent-child IDs. Recursive queries are handled in services.

### Testing & Linting
```bash
# Run all tests
go test ./...

# Format & Lint
go fmt ./...
golint ./...
go vet ./...
```

## Architecture Verities

### Dependency Injection (Wire)
- Entry point: `cmd/wire.go`
- Pattern: Handlers → Services → Repositories.
- **Agent Rule**: If you add a new service or repository, you **must** add it to `providerSet` in `cmd/wire.go` and run `wire ./cmd` to update `wire_gen.go`.

### RBAC & Data Model
- **Snowflake IDs**: Entities use `uint64` (Snowflake). Frontend receives them as strings to avoid precision loss.
- **Multi-Tenancy (OrgUnit)**: Users are bound to a specific `OrgUnit`. Permissions are scoped to the user's unit and its descendants.
- **Audit Logs**: Any state-changing operation (POST/PUT/DELETE) should be recorded in `AuditLog`.

### Project Layers
- `api/handler`: JSON binding & response (uses `domain/dto`)
- `service`: Business logic & Permission checks
- `repository`: GORM queries (uses `database.WithContext(ctx)`)
- `domain/`: Interfaces and shared types (GORM entities, DTOs)

## Operational Gotchas
- **Database**: Uses GORM. Ensure `global.G_DB` is initialized.
- **Redis**: Required for Rate Limiting and Permission Caching if enabled in `config.yml`.
- **Swagger**: Annotations must be on Handler methods. Run `make swagger` to update.
- **Environment**: `.env` overrides `config.yml`. Use `ADMIN_EMAIL` and `ADMIN_PASSWORD` with `make create-admin` for automation.

## Coding Style
- Indentation: **Tabs**
- Comments: **Chinese** for logic/business explanations, **English** for standard API docs (Swagger).
- Error Handling: Return errors to the caller. Use `result.ErrorResponse()` only in handlers.
- Context: Always propagate `ctx` from Handler down to Repository for timeout and trace tracking.

### Testing
No test framework is currently configured. When adding tests, use Go's built-in `testing` package.

```bash
# Run all tests
go test ./...
 
# Run a specific test file
go test ./path/to/package -run TestFunctionName
 
# Run tests with coverage
go test ./... -cover
 
# Run tests with verbose output
go test ./... -v
```

### Other Commands
```bash
# Format code
go fmt ./...
 
# Run linter (golint)
golint ./...
 
# Check for issues
go vet ./...
```
 
## RBAC 权限管理
 
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

### 初始化系统管理员
使用 CLI 命令创建系统管理员：

```bash
make create-admin
```

按照提示输入管理员邮箱和密码。

### 权限服务
- **PermissionService**：权限检查服务
  - `CheckPermission(userID, resource, method)` - API级别权限检查
  - `CheckEntityPermission(userID, entityType, entityID, action)` - 实体级别权限检查
  - `GetUserPermissions(userID)` - 获取用户权限
  - `GetUserOrgScope(userID)` - 获取用户组织范围
  - `ClearUserCache(userID)` - 清除用户缓存

### 权限中间件
- **RBACMiddleware**：权限控制中间件
  - `CheckPermission(resource)` - API级别权限检查
  - `CheckEntityPermission(entityType, action)` - 实体级别权限检查

### 组织树可见性
- **上级可见下级**：父节点组织的用户可以看到所有子组织的数据
- **下级不可见上级**：子节点组织的用户不能看到父组织的数据
- **同级不可见**：同一级别的组织用户不能互相看到对方的数据

### 审计日志
支持多种查询方式：

```bash
# 按操作者查询
GET /audit-logs?operator_id=1

# 按目标查询
GET /audit-logs/target?target_type=role&target_id=1

# 按时间范围查询
GET /audit-logs/time?start_time=2024-01-01&end_time=2024-01-31
```

### 权限缓存
当启用 Redis 时，用户权限会被缓存，提升性能：

```yaml
redis:
  Enabled: true
```

权限变更时自动清除缓存。

### 通配符权限
支持通配符匹配，如：

- `/users/*` - 匹配所有用户相关API
- `entity:*` - 匹配所有实体
- `entity:user:*` - 匹配所有用户实体操作

## Code Style Guidelines


### File & Package Organization
- Use lowercase, single-word package names (e.g., `handler`, `service`)
- Organize by layer (api, domain, service, repository, internal, pkg, util)
- Entry point: `cmd/main.go`

### Imports
- Third-party imports first, then standard library, separated with blank lines
- Group imports logically with blank lines between packages

### Naming Conventions
- Functions/Methods: camelCase (`GetUserByEmail`, `HashPassword`)
- Variables: camelCase (`contextTimeOut`, `authHeader`)
- Constants: UPPER_SNAKE_CASE (use config instead where possible)
- Types/Interfaces: PascalCase (`User`, `LoginService`)
- File names: snake_case (`user_service.go`)
- Interface methods: PascalCase (`Create`, `GetByEmail`)

### Formatting
- Use tabs for indentation
- Use Chinese comments for documentation
- One blank line between function definitions

### Type Conventions
- Use generics for flexible types (`ResponseResult[T any]`)
- All entities embed `global.G_MODEL` (ID, CreatedAt, UpdatedAt, DeletedAt)
- DTOs: Add "Request" or "Response" suffix (`LoginRequest`)
- Entities: Singular form (`User`, not `Users`)

### Error Handling
- Always return errors, never ignore them
- Use context with timeout in service/repository layers
- Return GORM errors directly from repository
- Use `result.ErrorResponse()` for HTTP error responses
- Return zero values on error (e.g., `entity.User{}`)

### Struct Conventions
- Use `json` tags for API response/request structs
- Use `gorm` tags for database entities
- Use `binding` tags for validation (`required`, `email`, `min=6`)
- Use `omitempty` for optional JSON fields
- Use `json:"-"` for sensitive fields

### Handler Patterns
- Add Swagger comments to all handler methods
- Use `result.SimpleSuccessResponse()`, `result.SuccessResponse()`, or `result.ErrorResponse()`

### Service Layer
- Define interfaces in `domain/services.go`, implement in `service/`
- Pass context as first parameter
- Create context with timeout for all repository calls
- Inject repository through constructor

### Repository Layer
- Define interfaces in `domain/repositories.go`, implement in `repository/`
- Use `database.WithContext(ctx)` for all queries
- Return `gorm.DB.Error` directly

### Configuration
- Use `config/config.yml` and access via package-level vars like `config.CfgServer.HttpPort`

### Domain-Driven Design
- `domain/` contains interfaces, entities, DTOs, result types
- High-level modules depend on abstractions (interfaces), not implementations
- Clear separation: handlers (HTTP) → services (logic) → repositories (data)

### Global Variables
- `global.G_DB` for GORM, `global.G_REDIS` for Redis, initialized in `bootstrap/Boot()`

### Comments & Documentation
- Add Swagger annotations to all handler methods
- Use Chinese comments for business logic explanations
- Add field comments for Swagger generation

### Middleware
- Location: `api/middleware/`
- Return `gin.HandlerFunc` from factory functions
- Use `result.ErrorResponse()` and `c.Abort()` for errors, `c.Next()` to continue

### Constants & Magic Values
- Use config YAML instead of hardcoded values
- Access timeouts via `config.CfgTimeout.ContextTimeout`
- Use `config.CfgToken` for JWT secrets

### Best Practices
- Propagate context from handlers → services → repositories
- Use `defer` for cleanup (e.g., `defer cancel()`)
- Use pointers for entity structs in services (`*entity.User`)
- Return zero values on error (e.g., `entity.User{}`)
- English for HTTP error messages, Chinese for internal comments

### Project Structure
```
cmd/              - Application entry point
api/              - HTTP layer (handlers, middleware, routes)
domain/           - Domain layer (entities, DTOs, interfaces, result types)
service/          - Business logic layer
repository/       - Data access layer
internal/         - Private packages (tokenutil, redisutil)
pkg/              - Shared libraries (database, redis, logger, captcha)
util/             - Utility functions
config/           - Configuration management
global/           - Global variables and shared types
bootstrap/        - Application initialization
```
