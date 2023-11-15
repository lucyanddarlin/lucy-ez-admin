package types

type AddNoticeRequest struct {
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
