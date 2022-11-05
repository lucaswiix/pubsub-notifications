package dto

type NotifyDTO struct {
	ID            string `json:"id,omitempty"`
	Message       string `json:"message" binding:"required,min=5"`
	Title         string `json:"title" binding:"required,min=5"`
	Image         string `json:"image" binding:"required,endswith=.png"`
	Type          string `json:"type" binding:"required,Enum=web"`
	ToUserID      string `json:"to_user_id,omitempty" binding:"uuid"`
	SchedulerDate string `json:"scheduler_datetime,omitempty" binding:"IsAfterNow" time_format:"2006-01-02 15:04:05" `
	Status        string `json:"status,omitempty"`
}

type CleanNotify struct {
	Message string `json:"message"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	Type    string `json:"type"`
	UserID  string `json:"to_user_id"`
}

const (
	Failed   = "failed"
	Sent     = "sent"
	Normal   = "normal"
	IsOptOut = "opt-out"
)
