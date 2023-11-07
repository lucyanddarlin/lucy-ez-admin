package types

type AddTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    int64  `json:"parent_id" binding:"required"`
}
