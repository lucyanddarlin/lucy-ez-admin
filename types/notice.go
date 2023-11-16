package types

type AddNoticeRequest struct {
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateNoticeRequest struct {
	ID      int64  `json:"id" binding:"required"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  *bool  `json:"status"`
}

type DeleteNoticeRequest struct {
	ID int64 `json:"id" binding:"required"`
}
