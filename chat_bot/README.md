# 豆包大模型聊天服务

这是一个使用Eino框架调用火山引擎豆包大模型的简单聊天服务。

## 功能特点

- 支持普通对话模式和流式对话模式
- 显示Token使用情况统计
- 自动管理对话历史
- 简单的命令行交互界面

## 环境要求

- Go 1.18或更高版本
- 火山引擎账号和API密钥
- 豆包大模型的访问权限

## 安装步骤

1. 克隆代码库

```bash
git clone https://github.com/Wenrh2004/aweone-eino-learn.git
cd aweone-eino-learn/chat-bot
```

2. 安装依赖

```bash
go mod tidy
```

## 使用方法

1. 设置环境变量

```bash
export ARK_API_KEY="你的火山引擎API密钥"
export ARK_MODEL_ID="你的豆包大模型端点ID"  # 可选，默认使用代码中设置的ID
```

2. 运行程序

```bash
go run main.go
```

3. 交互命令

- 输入任意文本与模型对话
- 输入 `stream` 切换到流式模式
- 输入 `normal` 切换到普通模式
- 输入 `quit` 退出程序

## 代码说明

- 使用Eino框架的ChatModel组件
- 通过ARK实现连接到火山引擎豆包大模型
- 支持两种对话模式：
  - 普通模式：等待完整回复后显示
  - 流式模式：实时显示模型生成的内容

## 注意事项

- 请确保正确设置API密钥和模型ID
- 流式模式下可以获得更好的交互体验
- 程序会自动管理对话历史，保留最近的10条消息