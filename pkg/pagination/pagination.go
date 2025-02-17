package pagination

import (
	"math"

	"gorm.io/gorm"
)

// Pagination 定义分页查询的结果结构
// 包含分页信息和查询结果数据
type Pagination struct {
	Page      int         `json:"page`      // 当前页码
	PageSize  int         `json:"pageSize`  // 每页数量
	Total     int64       `json:"total`     // 总记录数
	TotalPage int         `json:"totalPage` // 总页数
	Data      interface{} `json:"data`      // 数据
}

// PaginationBuilder 分页查询构造器
// 使用构建器模式，支持链式调用来构建复杂的分页查询
type PaginationBuilder struct {
	db         *gorm.DB    // 数据库连接实例
	page       int         // 当前页
	pageSize   int         // 每页数量
	model      interface{} // 模型
	preloads   []string    // 添加预加载
	joins      []string    // 添加连接查询
	selects    []string    // 添加查询字段
	conditions []Condition // 添加查询条件
	orders     []string    // 添加排序字段
	groups     []string    // 添加分组字段
	havings    []Condition // 添加HAVING条件
}

// Condition 定义查询条件的结构
// 用于Where和Having查询条件的封装
type Condition struct {
	Query interface{}   // 查询条件表达式
	Args  []interface{} // 查询条件参数
}

// NewPaginationBuilder 创建分页构造器的新实例
// 参数:
//   - db: *gorm.DB GORM数据库连接实例
//
// 返回:
//   - *PaginationBuilder 新的分页构造器实例
func NewPaginationBuilder(db *gorm.DB) *PaginationBuilder {
	return &PaginationBuilder{
		db:       db,
		page:     1,
		pageSize: 10,
	}
}

// SetPage 设置要查询的页码
// 参数:
//   - page: int 页码，从1开始
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) SetPage(page int) *PaginationBuilder {
	if page > 0 {
		p.page = page
	}
	return p
}

// SetPageSize 设置每页显示的记录数
// 参数:
//   - pageSize: int 每页记录数
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) SetPageSize(pageSize int) *PaginationBuilder {
	if pageSize > 0 {
		p.pageSize = pageSize
	}
	return p
}

// Model 设置要查询的数据模型
// 参数:
//   - model: interface{} 模型结构体指针
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) Model(model interface{}) *PaginationBuilder {
	p.model = model
	return p
}

// Preload 添加预加载关系
// 参数:
//   - query: string 要预加载的关系名称
//
// 示例:
//
//	builder.Preload("User")  // 预加载User关系
//	builder.Preload("Orders") // 预加载Orders关系
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) Preload(query string) *PaginationBuilder {
	p.preloads = append(p.preloads, query)
	return p
}

// Join 添加表连接查询
// 参数:
//   - query: string JOIN查询语句
//
// 示例:
//
//	builder.Join("LEFT JOIN orders ON users.id = orders.user_id")
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) Join(query string) *PaginationBuilder {
	p.joins = append(p.joins, query)
	return p
}

// Select 设置要查询的字段
// 参数:
//   - query: string 要查询的字段列表
//
// 示例:
//
//	builder.Select("id, name, email")
//	builder.Select("COUNT(*) as count")
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) Select(query string) *PaginationBuilder {
	p.selects = append(p.selects, query)
	return p
}

// Where 添加查询条件
// 参数:
//   - query: interface{} 查询条件表达式
//   - args: ...interface{} 查询条件参数
//
// 示例:
//
//	builder.Where("age > ?", 18)
//	builder.Where("status = ?", "active")
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) Where(query interface{}, args ...interface{}) *PaginationBuilder {
	p.conditions = append(p.conditions, Condition{
		Query: query,
		Args:  args,
	})
	return p
}

// OrderBy 添加排序规则
// 参数:
//   - order: string 排序表达式
//
// 示例:
//
//	builder.OrderBy("created_at DESC")
//	builder.OrderBy("age ASC")
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) OrderBy(order string) *PaginationBuilder {
	p.orders = append(p.orders, order)
	return p
}

// GroupBy 添加分组字段
// 参数:
//   - group: string 分组字段
//
// 示例:
//
//	builder.GroupBy("department_id")
//	builder.GroupBy("role, status")
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) GroupBy(group string) *PaginationBuilder {
	p.groups = append(p.groups, group)
	return p
}

// Having 添加HAVING条件
// 参数:
//   - query: interface{} HAVING条件表达式
//   - args: ...interface{} HAVING条件参数
//
// 示例:
//
//	builder.Having("COUNT(*) > ?", 5)
//	builder.Having("SUM(amount) > ?", 1000)
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) Having(query interface{}, args ...interface{}) *PaginationBuilder {
	p.havings = append(p.havings, Condition{
		Query: query,
		Args:  args,
	})
	return p
}

// Build 执行分页查询并构建结果
// 参数:
//   - result: interface{} 查询结果的接收对象指针
//
// 返回:
//   - *Pagination 分页查询结果
//   - error 查询过程中的错误
//
// 示例:
//
//	var users []User
//	pagination, err := builder.Build(&users)
func (p *PaginationBuilder) Build(result interface{}) (*Pagination, error) {
	var total int64
	query := p.db.Model(p.model)

	// 查询条件
	for _, condition := range p.conditions {
		query = query.Where(condition.Query, condition.Args...)
	}

	// 连接查询
	for _, join := range p.joins {
		query = query.Joins(join)
	}

	// 查询字段
	if len(p.selects) > 0 {
		// 将[]string 类型转换为 []interface{} 类型
		interfaces := make([]interface{}, len(p.selects[1:]))
		for i, v := range p.selects[1:] {
			interfaces[i] = v
		}
		query = query.Select(p.selects[0], interfaces...)
	}

	// 添加分组
	for _, group := range p.groups {
		query = query.Group(group)
	}

	// 添加HAVING条件
	for _, having := range p.havings {
		query = query.Having(having.Query, having.Args...)
	}

	// 克隆查询以计算总记录数
	countQuery := query
	if len(p.groups) > 0 {
		// 如果有分组，则计算分组后的记录数
		countQuery = countQuery.Session(&gorm.Session{})
		if err := countQuery.Count(&total).Error; err != nil {
			return nil, err
		}
	} else {
		// 无分组时的普通计数
		if err := countQuery.Count(&total).Error; err != nil {
			return nil, err
		}
	}

	// 添加排序
	for _, order := range p.orders {
		query = query.Order(order)
	}

	// 分页查询
	offset := (p.page - 1) * p.pageSize
	query = query.Offset(offset).Limit(p.pageSize)

	// 预加载
	for _, preload := range p.preloads {
		query = query.Preload(preload)
	}

	// 执行查询
	if err := query.Find(result).Error; err != nil {
		return nil, err
	}

	// 构建分页响应
	totalPage := int(math.Ceil(float64(total) / float64(p.pageSize)))
	return &Pagination{
		Page:      p.page,
		PageSize:  p.pageSize,
		Total:     total,
		TotalPage: totalPage,
		Data:      result,
	}, nil
}
