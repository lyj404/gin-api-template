package pagination

// Pagination 定义分页查询的结果结构
// 包含分页信息和查询结果数据
type Pagination struct {
	Page      int         `json:"page"`      // 当前页码
	PageSize  int         `json:"pageSize"`  // 每页数量
	Total     int64       `json:"total"`     // 总记录数
	TotalPage int         `json:"totalPage"` // 总页数
	Data      interface{} `json:"data"`      // 数据
}

// Condition 定义查询条件的结构
// 用于Where和Having查询条件的封装
type Condition struct {
	Query interface{}   // 查询条件表达式
	Args  []interface{} // 查询条件参数
}

const (
	ASC  string = "ASC"  // 升序排列
	DESC string = "DESC" // 降序排列
)
