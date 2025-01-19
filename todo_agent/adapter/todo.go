package adapter

import (
	"context"
	"fmt"
	"github.com/Wenrh2004/eino-learn-demo/api"
	"github.com/Wenrh2004/eino-learn-demo/service"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"log"
	"strconv"
	"time"
)

type TodoAdapter interface {
	AddTodo(ctx context.Context)
	DeleteTodo(ctx context.Context)
	UpdateTodo(ctx context.Context)
	ListTodo(ctx context.Context)
	QueryAssistant(ctx context.Context)
}

type todoAdapter struct {
	service.TodoService
	msg   []*schema.Message
	agent compose.Runnable[[]*schema.Message, []*schema.Message]
}

func NewTodoAdapter(service service.TodoService, agent compose.Runnable[[]*schema.Message, []*schema.Message]) TodoAdapter {
	messages := []*schema.Message{
		{
			Role:    schema.Assistant,
			Content: "Hello, I'm Eino, your assistant. I am good at schedule planning, transaction management and data sorting, which can help you make efficient work plans and improve work efficiency.",
		},
	}
	return &todoAdapter{
		TodoService: service,
		msg:         messages,
		agent:       agent,
	}
}

func (t *todoAdapter) AddTodo(ctx context.Context) {
	fmt.Print("Please input the content of the todo item: ")
	var content string
	_, err := fmt.Scanln(&content)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	fmt.Print("Please input the description of the todo item: ")
	var description string
	_, err = fmt.Scanln(&description)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	fmt.Print("Please input the remark of the todo item: ")
	var remark string
	_, err = fmt.Scanln(&remark)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	fmt.Print("Please input the started time of the todo item: ")
	var startedAt string
	_, err = fmt.Scanln(&startedAt)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	startedAtTime, err := time.Parse(time.Kitchen, startedAt)
	if err != nil {
		log.Fatalf("time.Parse failed, err = %v", err)
	}
	startedAtUnix := startedAtTime.Unix()
	fmt.Print("Please input the deadline of the todo item: ")
	var deadline string
	_, err = fmt.Scanln(&deadline)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	deadlineTime, err := time.Parse(time.Kitchen, deadline)
	if err != nil {
		log.Fatalf("time.Parse failed, err = %v", err)
	}
	if !startedAtTime.Before(deadlineTime) {
		log.Fatalf("deadline should be after started time")
	}
	deadlineUnix := deadlineTime.Unix()
	if err = t.TodoService.AddTodo(ctx, &api.TodoAddParams{
		Content:     content,
		Description: description,
		Remark:      remark,
		StartedAt:   &startedAtUnix,
		Deadline:    &deadlineUnix,
	}); err != nil {
		fmt.Printf("add todo failed, err = %v, please try again\n", err)
	}
	fmt.Println("add todo success")
}

func (t *todoAdapter) DeleteTodo(ctx context.Context) {
	fmt.Print("Please input the content of the todo item: ")
	var content string
	_, err := fmt.Scanln(&content)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	if err = t.TodoService.DeleteTodo(ctx, content); err != nil {
		fmt.Printf("delete todo failed, err = %v, please try again\n", err)
	}
}

func (t *todoAdapter) UpdateTodo(ctx context.Context) {
	fmt.Print("Please input the id of the todo item: ")
	var id string
	_, err := fmt.Scanln(&id)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	fmt.Print("Please input the content of the todo item: ")
	var content string
	_, err = fmt.Scanln(&content)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	fmt.Print("Please input the description of the todo item: ")
	var description string
	_, err = fmt.Scanln(&description)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	fmt.Print("Please input the remark of the todo item: ")
	var remark string
	_, err = fmt.Scanln(&remark)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	fmt.Print("Please input the started time of the todo item: ")
	var startedAt string
	_, err = fmt.Scanln(&startedAt)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	startedAtTime, err := time.Parse(time.Kitchen, startedAt)
	if err != nil {
		log.Fatalf("time.Parse failed, err = %v", err)
	}
	startedAtUnix := startedAtTime.Unix()
	fmt.Print("Please input the deadline of the todo item: ")
	var deadline string
	_, err = fmt.Scanln(&deadline)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	deadlineTime, err := time.Parse(time.Kitchen, deadline)
	if err != nil {
		log.Fatalf("time.Parse failed, err = %v", err)
	}
	if !startedAtTime.Before(deadlineTime) {
		log.Fatalf("deadline should be after started time")
	}
	deadlineUnix := deadlineTime.Unix()
	fmt.Print("Please input the done status of the todo item: ")
	var done string
	_, err = fmt.Scanln(&done)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	doneBool, err := strconv.ParseBool(done)
	if err != nil {
		log.Fatalf("strconv.ParseBool failed, err = %v", err)
	}
	if err = t.TodoService.UpdateTodo(ctx, &api.TodoUpdateParams{
		Id:          id,
		Content:     content,
		Description: description,
		Remark:      remark,
		StartedAt:   &startedAtUnix,
		Deadline:    &deadlineUnix,
		Done:        doneBool,
	}); err != nil {
		fmt.Printf("update todo failed, err = %v, please try again\n", err)
	}
	fmt.Println("update todo success")
}

func (t *todoAdapter) ListTodo(ctx context.Context) {
	listTodo, err := t.TodoService.ListTodo(ctx)
	if err != nil {
		log.Fatalf("list todo failed, err = %v", err)
	}
	fmt.Println(listTodo)
}

func (t *todoAdapter) QueryAssistant(ctx context.Context) {
	fmt.Print("Please input the query of the assistant: ")
	var query string
	_, err := fmt.Scanln(&query)
	if err != nil {
		log.Fatalf("fmt.Scanln failed, err = %v", err)
	}
	t.msg = append(t.msg, schema.UserMessage(query))
	resp, err := t.agent.Invoke(ctx, t.msg)
	if err != nil {
		fmt.Printf("query assistant failed, err = %v, please try again\n", err)
	}
	for i, message := range resp {
		fmt.Printf("message[%d][%s]: %s\n", i, message.Role, message.Content)
	}
}
