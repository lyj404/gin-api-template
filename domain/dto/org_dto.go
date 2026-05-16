package dto

type CreateOrgUnitRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint64  `json:"parent_id"`
}

type UpdateOrgUnitRequest struct {
	Name     string `json:"name"`
	ParentID *uint64  `json:"parent_id"`
}

type OrgUnitResponse struct {
	ID       uint64   `json:"id"`
	Name     string `json:"name"`
	ParentID *uint64  `json:"parent_id"`
	Path     string `json:"path"`
	Level    int    `json:"level"`
}
