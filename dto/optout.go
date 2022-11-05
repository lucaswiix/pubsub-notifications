package dto

type OptOut struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}
