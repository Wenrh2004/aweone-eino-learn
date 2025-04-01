package tool

import (
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	"github.com/Wenrh2004/travel_assistant/pkg/third/amap"
)

type WeatherService struct {
	client *amap.Client
}

func NewWeatherService(client *amap.Client) *WeatherService {
	return &WeatherService{
		client: client,
	}
}

func (w *WeatherService) GetCurrentWeatherTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"Current Weather Tool",
		"search current weather",
		w.client.GetWeatherLive,
	)
}

func (w *WeatherService) GetWeatherForecastTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"Weather Forecast Tool",
		"search weather forecast",
		w.client.GetWeatherForecast,
	)
}
