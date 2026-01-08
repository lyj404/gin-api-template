package dto

type BindResourceRequest struct {
	RoleID  uint `json:"role_id" binding:"required"`
	IsWrite bool `json:"is_write"`
}

type ResourceListResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Method      string `json:"method"`
	Entity      string `json:"entity"`
	Action      string `json:"action"`
	Description string `json:"description"`
}
