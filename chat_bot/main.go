package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 从环境变量获取API密钥
	apiKey := os.Getenv("ARK_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置环境变量ARK_API_KEY")
	}

	// 从环境变量获取模型ID
	modelID := os.Getenv("ARK_MODEL_ID")
	if modelID == "" {
		log.Println("未设置ARK_MODEL_ID环境变量，使用默认模型ID")
		modelID = "ep-xxx" // 替换为你的豆包大模型端点ID
	}

	// 设置超时时间
	timeout := 30 * time.Second

	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  apiKey,
		Region:  "cn-beijing", // 火山引擎区域
		Model:   modelID,
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatalf("初始化模型失败: %v", err)
	}

	fmt.Println("欢迎使用豆包大模型聊天服务！")
	fmt.Println("输入'quit'退出，输入'stream'切换到流式模式，输入'normal'切换到普通模式")

	// 默认使用普通模式
	streamMode := false

	// 创建系统消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一个由火山引擎提供的豆包大模型助手，请尽可能地回答用户的问题。"),
	}

	// 开始交互循环
	for {
		// 获取用户输入
		fmt.Print("\n用户: ")
		var input string
		fmt.Scanln(&input)

		// 检查是否退出
		if input == "quit" {
			break
		}

		// 检查是否切换模式
		if input == "stream" {
			streamMode = true
			fmt.Println("已切换到流式模式")
			continue
		} else if input == "normal" {
			streamMode = false
			fmt.Println("已切换到普通模式")
			continue
		}

		// 添加用户消息
		messages = append(messages, schema.UserMessage(input))

		// 根据模式选择不同的调用方式
		if streamMode {
			// 流式模式
			fmt.Print("\n助手: ")
			reader, err := model.Stream(ctx, messages)
			if err != nil {
				log.Printf("流式生成失败: %v", err)
				continue
			}
			defer reader.Close()

			// 处理流式内容
			var fullContent string
			for {
				chunk, err := reader.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("接收流数据失败: %v", err)
					break
				}
				fmt.Print(chunk.Content)
				fullContent += chunk.Content
			}
			// 添加助手回复到消息历史
			messages = append(messages, schema.AssistantMessage(fullContent, nil))
		} else {
			// 普通模式
			response, err := model.Generate(ctx, messages)
			if err != nil {
				log.Printf("生成回复失败: %v", err)
				continue
			}

			// 输出回复
			fmt.Printf("\n助手: %s\n", response.Content)

			// 添加助手回复到消息历史
			messages = append(messages, response)

			// 输出Token使用情况（如果有）
			if usage := response.ResponseMeta.Usage; usage != nil {
				fmt.Printf("Token使用情况 - 提示: %d, 生成: %d, 总计: %d\n",
					usage.PromptTokens, usage.CompletionTokens, usage.TotalTokens)
			}
		}

		// 限制消息历史长度，保留最近的10条消息
		if len(messages) > 11 { // 系统消息 + 10条对话
			messages = append([]*schema.Message{messages[0]}, messages[len(messages)-10:]...)
		}
	}

	fmt.Println("感谢使用豆包大模型聊天服务，再见！")
}
