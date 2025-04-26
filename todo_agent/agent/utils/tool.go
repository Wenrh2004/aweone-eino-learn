package utils

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/spf13/viper"

	"github.com/Wenrh2004/todo_agent/api"
	"github.com/Wenrh2004/todo_agent/service"
)

func NewTools(ctx context.Context, conf *viper.Viper, todoService service.TodoService) []tool.BaseTool {
	updateTool, err := utils.InferTool("update_todo", "Update a todo item, eg: content,deadline...", UpdateTodoFunc)
	if err != nil {
		panic(err)
	}

	listTool := &ListTodoTool{
		TodoService: todoService,
	}

	searchTool, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{
		DDGConfig: &ddgsearch.Config{
			Proxy: conf.GetString("ddg.proxy"),
			// Timeout: conf.GetDuration("ddg.timeout"),
			Cache: conf.GetBool("ddg.cache"),
		},
	})
	if err != nil {
		panic(err)
	}

	return []tool.BaseTool{
		getAddTodoTool(), // 使用 NewTool 方式
		updateTool,       // 使用 InferTool 方式
		listTool,         // 使用结构体实现方式, 此处未实现底层逻辑
		searchTool,
	}
}

// 获取添加 todo 工具
// 使用 utils.NewTool 创建工具
func getAddTodoTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "add_todo",
		Desc: "Add a todo item",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Desc:     "The content of the todo item",
				Type:     schema.String,
				Required: true,
			},
			"description": {
				Desc:     "The description of the todo item",
				Type:     schema.String,
				Required: true,
			},
			"remark": {
				Desc: "The remark of the todo item, to save some additional information, eg: search information, repository information",
				Type: schema.String,
			},
			"started_at": {
				Desc: "The started time of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
			"deadline": {
				Desc: "The deadline of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
		}),
	}

	return utils.NewTool(info, AddTodoFunc)
}

// ListTodoTool
// 获取列出 todo 工具
type ListTodoTool struct {
	service.TodoService
}

func (lt *ListTodoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "list_todo",
		Desc: "List all todo items",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"finished": {
				Desc:     "filter todo items if finished",
				Type:     schema.Boolean,
				Required: false,
			},
		}),
	}, nil
}

type TodoUpdateParams struct {
	Id        string  `json:"id" jsonschema:"description=id of the todo"`
	Content   *string `json:"content,omitempty" jsonschema:"description=content of the todo"`
	StartedAt *int64  `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
	Deadline  *int64  `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
	Done      *bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

func (lt *ListTodoTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	log.Printf("invoke utils list_todo: %s", argumentsInJSON)

	// 具体的调用逻辑
	listTodo, err := lt.TodoService.ListTodo(ctx)
	if err != nil {
		return "", err
	}
	return listTodo, nil
}

func AddTodoFunc(ctx context.Context, params *api.TodoAddParams) (string, error) {
	log.Printf("invoke utils add_todo: %+v", params)
	// 具体的调用逻辑
	if err := service.GetTodoService().AddTodo(ctx, params); err != nil {
		return `"msg": "add todo failed, please try again"`, err
	}
	return `{"msg": "add todo success"}`, nil
}

func UpdateTodoFunc(ctx context.Context, params *api.TodoUpdateParams) (string, error) {
	log.Printf("invoke utils update_todo: %+v", params)

	// 具体的调用逻辑
	if err := service.GetTodoService().UpdateTodo(ctx, params); err != nil {
		return `"msg": "add todo failed, please try again"`, err
	}

	return `{"msg": "update todo success"}`, nil
}
