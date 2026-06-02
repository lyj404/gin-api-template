package dto

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int    `form:"page" json:"page"`         // 当前页码，默认1
	PageSize int    `form:"page_size" json:"page_size"` // 每页数量，默认10
	Keyword  string `form:"keyword" json:"keyword"`     // 搜索关键词
	OrderBy  string `form:"order_by" json:"order_by"`   // 排序字段
	Sort     string `form:"sort" json:"sort"`           // 排序方式：asc/desc
}

// SetDefaults 设置分页默认值
func (p *PaginationRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Page      int         `json:"page"`       // 当前页码
	PageSize  int         `json:"page_size"`  // 每页数量
	Total     int64       `json:"total"`      // 总记录数
	TotalPage int         `json:"total_page"` // 总页数
	Data      interface{} `json:"data"`       // 数据列表
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(page, pageSize int, total int64, data interface{}) *PaginationResponse {
	totalPage := int((total + int64(pageSize) - 1) / int64(pageSize))
	return &PaginationResponse{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
		Data:      data,
	}
}