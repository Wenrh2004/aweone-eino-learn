package adapter

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/devops"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sse"
	"go.uber.org/zap"

	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/domain/agent"
	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/log"
)

type LLMHandler struct {
	agent.ChatService
	logger *log.Logger
}

func NewLLMHandler(logger *log.Logger, domain agent.ChatService) *LLMHandler {
	return &LLMHandler{
		ChatService: domain,
		logger:      logger,
	}
}

type ChatRequest struct {
	Message string    `json:"message"`
	History []Message `json:"history,omitempty"`
}

type Message struct {
	Role      string            `json:"role"`
	Content   string            `json:"content"`
	ToolCalls []schema.ToolCall `json:"tool_calls"`
}

func (l *LLMHandler) Chat(ctx context.Context, c *app.RequestContext) {
	var req ChatRequest
	if err := devops.Init(ctx); err != nil {
		return
	}
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	history := make([]*schema.Message, 0, len(req.History))
	for _, message := range req.History {
		switch message.Role {
		case "user":
			history = append(history, schema.UserMessage(message.Content))
		case "assistant":
			history = append(history, schema.AssistantMessage(message.Content, message.ToolCalls))
		case "system":
			history = append(history, schema.SystemMessage(message.Content))
		default:
			c.JSON(http.StatusBadRequest, "invalid role")
		}
	}

	resp, err := l.ChatService.Query(ctx, req.Message, history)
	if err != nil {
		l.logger.Error("failed to query", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 客户端可以通过 Last-Event-ID 告知服务器收到的最后一个事件
	lastEventID := sse.GetLastEventID(c)
	hlog.CtxInfof(ctx, "last event ID: %s", lastEventID)

	// 在第一次渲染调用之前必须先行设置状态代码和响应头文件
	c.SetStatusCode(http.StatusOK)
	s := sse.NewStream(c)
	for {
		message, err := resp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err := s.Publish(&sse.Event{
					Event: "end",
					Data:  []byte("end flag"),
				})
				if err != nil {
					return
				}
				// finish
				break
			}
			// error
			l.logger.WithContext(ctx).Error("failed to recv: %v\n", zap.Error(err))
			return
		}
		msg, err := sonic.Marshal(message)
		if err != nil {
			return
		}
		event := &sse.Event{
			Event: "message",
			Data:  msg,
		}
		err = s.Publish(event)
		if err != nil {
			return
		}
	}
}
