package service

import (
	"context"
	"github.com/Wenrh2004/eino-learn-demo/api"
	"github.com/Wenrh2004/eino-learn-demo/repository"
	"github.com/bytedance/sonic"
)

var todo *todoService

type TodoService interface {
	AddTodo(ctx context.Context, param *api.TodoAddParams) error
	DeleteTodo(ctx context.Context, content string) error
	UpdateTodo(ctx context.Context, param *api.TodoUpdateParams) error
	ListTodo(ctx context.Context) (string, error)
}

type todoService struct {
	repository.TodoList
}

func NewTodoService(todoList repository.TodoList) TodoService {
	todo = &todoService{
		TodoList: todoList,
	}
	return todo
}

func GetTodoService() TodoService {
	return todo
}

func (t *todoService) AddTodo(ctx context.Context, param *api.TodoAddParams) error {
	model := param.ConvertToModel()
	if err := t.TodoList.AddTodo(ctx, model); err != nil {
		return err
	}
	return nil
}

func (t *todoService) DeleteTodo(ctx context.Context, content string) error {
	id, err := t.TodoList.GetTodoByContent(ctx, content)
	if err != nil {
		return err
	}
	if err = t.TodoList.DeleteTodo(ctx, id); err != nil {
		return err
	}
	return nil
}

func (t *todoService) UpdateTodo(ctx context.Context, param *api.TodoUpdateParams) error {
	model := param.ConvertToModel()
	if err := t.TodoList.UpdateTodo(ctx, model); err != nil {
		return err
	}
	return nil
}

func (t *todoService) ListTodo(ctx context.Context) (string, error) {
	todoList, err := t.TodoList.ListTodo(ctx)
	info, err := sonic.Marshal(todoList)
	if err != nil {
		return "", err
	}
	return string(info), nil
}
