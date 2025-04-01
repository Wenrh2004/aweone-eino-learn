package tool

import (
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	"github.com/Wenrh2004/travel_assistant/pkg/third/amap"
)

type POIService struct {
	client *amap.Client
}

func NewPOIService(client *amap.Client) *POIService {
	return &POIService{
		client: client,
	}
}

func (p *POIService) GetPOISearchTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"POI Search Tool",
		"search point of interest",
		p.client.SearchPOI,
	)
}
