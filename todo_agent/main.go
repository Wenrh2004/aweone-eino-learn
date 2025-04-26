package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Wenrh2004/todo_agent/adapter"
	"github.com/Wenrh2004/todo_agent/agent"
	"github.com/Wenrh2004/todo_agent/agent/utils"
	"github.com/Wenrh2004/todo_agent/config"
	"github.com/Wenrh2004/todo_agent/repository"
	"github.com/Wenrh2004/todo_agent/service"
)

func main() {
	ctx := context.Background()
	conf := config.NewConfig("config.yaml")

	todoList := repository.NewTodoList()
	todoService := service.NewTodoService(todoList)
	tools := utils.NewTools(ctx, conf, todoService)
	chatModel := utils.NewChatModel(ctx, conf)
	todoAgent := agent.NewAgent(ctx, tools, chatModel)
	todoAdapter := adapter.NewTodoAdapter(todoService, todoAgent)
	NewRobot(ctx, todoAdapter)
}

func NewRobot(ctx context.Context, a adapter.TodoAdapter) {
	fmt.Println("Hello, I'm Eino, your assistant. you can ask me the following questions: ")
	for {
		var userInput string
		fmt.Println("1. add a todo item")
		fmt.Println("2. remove a todo item")
		fmt.Println("3. update a todo item")
		fmt.Println("4. list the todo items")
		fmt.Println("5. query assistant")
		fmt.Println("6. quit")
		fmt.Print("Please input the number of the question: ")
		_, err := fmt.Scanln(&userInput)
		if err != nil {
			log.Fatalf("fmt.Scanln failed, err = %v", err)
		}
		userKey, err := strconv.Atoi(userInput)
		if err != nil {
			log.Fatalf("strconv.Atoi failed, err = %v", err)
		}
		switch userKey {
		case 1:
			a.AddTodo(ctx)
		case 2:
			a.DeleteTodo(ctx)
		case 3:
			a.UpdateTodo(ctx)
		case 4:
			a.ListTodo(ctx)
		case 5:
			a.QueryAssistant(ctx)
		case 6:
			return
		default:
			fmt.Println("invalid input")
		}
	}
}
