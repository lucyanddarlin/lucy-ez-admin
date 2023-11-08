package types

type AddTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    int64  `json:"parent_id" binding:"required"`
}

type UpdateTeamRequest struct {
	ID          int64  `json:"id" form:"id" binding:"required"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	ParentID    int64  `json:"parent_id" form:"parent_id"`
}

type DeleteTeamRequest struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}
