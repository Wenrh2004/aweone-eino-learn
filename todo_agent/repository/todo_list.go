package repository

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"

	"github.com/Wenrh2004/todo_agent/model"
)

// TodoList TODO 需求抽象
type TodoList interface {
	AddTodo(ctx context.Context, params *model.TodoListParams) error      // 新增 TODO
	DeleteTodo(ctx context.Context, id string) error                      // 删除 TODO
	UpdateTodo(ctx context.Context, params *model.TodoListParams) error   // 更新 TODO
	ListTodo(ctx context.Context) ([]*model.TodoListParams, error)        // 展示 TODO
	GetTodoByContent(ctx context.Context, content string) (string, error) // 通过内容获取 TODO
}

// todoList TODO 列表
type todoList struct {
	uid   *uuid.UUID
	Todos []*model.TodoListParams `json:"todos,omitempty"`
}

func NewTodoList() TodoList {
	uidGenerator, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return &todoList{
		uid:   &uidGenerator,
		Todos: make([]*model.TodoListParams, 0),
	}
}

// AddTodo 新增 TODO
func (t *todoList) AddTodo(ctx context.Context, params *model.TodoListParams) error {
	if len(t.Todos) != 0 {
		for _, todo := range t.Todos {
			if todo.Content == params.Content {
				return errors.New("todo already exists")
			}
		}
	}
	params.Id = t.uid.String()
	t.Todos = append(t.Todos, params)
	return nil
}

// DeleteTodo 删除 TODO
func (t *todoList) DeleteTodo(ctx context.Context, id string) error {
	if len(t.Todos) == 0 {
		return errors.New("the todo list is null")
	}

	for i, todo := range t.Todos {
		if todo.Id == id {
			t.Todos = append(t.Todos[:i], t.Todos[i+1:]...)
			break
		}
	}
	return nil
}

// UpdateTodo 更新 TODO
func (t *todoList) UpdateTodo(ctx context.Context, params *model.TodoListParams) error {
	if len(t.Todos) == 0 {
		return errors.New("the todo list is null")
	}
	for _, todo := range t.Todos {
		if todo.Id == params.Id {
			todo.Content = params.Content
			todo.StartedAt = params.StartedAt
			todo.Deadline = params.Deadline
			todo.Done = params.Done
			return nil
		}
	}
	return errors.New("todo not found")
}

// ListTodo 展示 TODO
func (t *todoList) ListTodo(ctx context.Context) ([]*model.TodoListParams, error) {
	// 按照开始时间排序
	sort.Slice(t.Todos, func(i, j int) bool {
		return *t.Todos[i].StartedAt < *t.Todos[j].StartedAt
	})

	return t.Todos, nil
}

// GetTodoByContent 通过内容获取 TODO
func (t *todoList) GetTodoByContent(ctx context.Context, content string) (string, error) {
	if len(t.Todos) == 0 {
		return "", errors.New("the todo list is null")
	}
	for _, todo := range t.Todos {
		if todo.Content == content {
			return todo.Id, nil
		}
	}
	return "", errors.New("todo not found")
}
