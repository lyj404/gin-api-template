package dto

type CreateOrgUnitRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateOrgUnitRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type OrgUnitResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
	Path     string `json:"path"`
	Level    int    `json:"level"`
}
