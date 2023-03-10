package dto

type CreateCommentRequest struct {
	Username string `json:"username"`
	IsiKomen string `json:"komen" binding:"required"`
	BlogID   uint64 `json:"blog_id"`
}
