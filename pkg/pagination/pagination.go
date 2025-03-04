package pagination

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
	if model != nil {
		p.model = model
	}
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
	if query != "" {
		p.preloads = append(p.preloads, query)
	}
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
	if query != "" {
		p.joins = append(p.joins, query)
	}
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
//	builder.OrderBy("created_at", pagination.DESC)
//	builder.OrderBy("age", pagination.ASC)
//
// 返回:
//   - *PaginationBuilder 分页构造器实例(链式调用)
func (p *PaginationBuilder) OrderBy(field string, order string) *PaginationBuilder {
	p.orderFields = append(p.orderFields, field+" "+string(order))
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
