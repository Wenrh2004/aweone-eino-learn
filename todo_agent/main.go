package main

import (
	"context"
	"fmt"
	"github.com/Wenrh2004/eino-learn-demo/adapter"
	"github.com/Wenrh2004/eino-learn-demo/repository"
	"github.com/Wenrh2004/eino-learn-demo/service"
	"log"
	"strconv"
)

func main() {
	ctx := context.Background()
	todoList := repository.NewTodoList()
	todoService := service.NewTodoService(todoList)
	todoAdapter := adapter.NewTodoAdapter(todoService)
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
		fmt.Println("5. quit")
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
			return
		default:
			fmt.Println("invalid input")
		}
	}
}
