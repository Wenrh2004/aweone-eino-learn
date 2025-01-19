package model

// TodoListParams TODO 信息
type TodoListParams struct {
	Id          string `json:"id"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
	StartedAt   *int64 `json:"started_at,omitempty"`
	Deadline    *int64 `json:"deadline,omitempty"`
	Done        bool   `json:"done"`
}
