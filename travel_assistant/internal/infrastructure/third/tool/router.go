package tool

import (
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	"github.com/Wenrh2004/travel_assistant/pkg/third/amap"
)

type RouterService struct {
	client *amap.Client
}

func NewRouterService(client *amap.Client) *RouterService {
	return &RouterService{
		client: client,
	}
}

func (r *RouterService) GetDrivingRouterTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"Driving Router Tool",
		"search route by driving",
		r.client.GetDrivingRoute,
	)
}

func (r *RouterService) GetWalkingRouterTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"Walking Router Tool",
		"search route by walking",
		r.client.GetWalkingRoute,
	)
}
