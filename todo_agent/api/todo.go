package api

import "github.com/Wenrh2004/eino-learn-demo/model"

const (
	StatusTodo = false
)

// TodoAddParams TODO 添加参数
type TodoAddParams struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
	StartedAt   *int64 `json:"started_at,omitempty"`
	Deadline    *int64 `json:"deadline,omitempty"`
}

type TodoUpdateParams struct {
	Id          string `json:"id" jsonschema:"description=id of the todo"`
	Content     string `json:"content,omitempty" jsonschema:"description=content of the todo"`
	Description string `json:"description,omitempty" jsonschema:"description=description of the todo"`
	Remark      string `json:"remark,omitempty" jsonschema:"description=remark of the todo"`
	StartedAt   *int64 `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
	Deadline    *int64 `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
	Done        bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

type TodoListParams []*model.TodoListParams

func (t *TodoAddParams) ConvertToModel() *model.TodoListParams {
	return &model.TodoListParams{
		Content:     t.Content,
		Description: t.Description,
		Remark:      t.Remark,
		StartedAt:   t.StartedAt,
		Deadline:    t.Deadline,
		Done:        StatusTodo,
	}
}

func (t *TodoUpdateParams) ConvertToModel() *model.TodoListParams {
	return &model.TodoListParams{
		Id:          t.Id,
		Content:     t.Content,
		Description: t.Description,
		Remark:      t.Remark,
		StartedAt:   t.StartedAt,
		Deadline:    t.Deadline,
		Done:        t.Done,
	}
}
