package amap

import (
	"context"
	"net/url"

	"github.com/bytedance/sonic"
)

type WeatherRequest struct {
	City string `json:"city" jsonschema_description:"the city for search"` // 城市
}

// WeatherResponse 天气查询响应
type WeatherResponse struct {
	BaseResponse
	Lives     []WeatherLive     `json:"lives,omitempty"`     // 实时天气
	Forecasts []WeatherForecast `json:"forecasts,omitempty"` // 天气预报
}

// WeatherLive 实时天气信息
type WeatherLive struct {
	Province      string `json:"province"`      // 省份
	City          string `json:"city"`          // 城市
	Adcode        string `json:"adcode"`        // 区域编码
	Weather       string `json:"weather"`       // 天气现象
	Temperature   string `json:"temperature"`   // 实时温度
	WindDirection string `json:"winddirection"` // 风向
	WindPower     string `json:"windpower"`     // 风力
	Humidity      string `json:"humidity"`      // 湿度
	ReportTime    string `json:"reporttime"`    // 数据发布时间
}

// WeatherForecast 天气预报信息
type WeatherForecast struct {
	City       string        `json:"city"`       // 城市
	Adcode     string        `json:"adcode"`     // 区域编码
	Province   string        `json:"province"`   // 省份
	Reporttime string        `json:"reporttime"` // 数据发布时间
	Casts      []WeatherCast `json:"casts"`      // 预报数据
}

// WeatherCast 天气预报数据
type WeatherCast struct {
	Date         string `json:"date"`         // 日期
	Week         string `json:"week"`         // 星期
	DayWeather   string `json:"dayweather"`   // 白天天气现象
	NightWeather string `json:"nightweather"` // 晚上天气现象
	DayTemp      string `json:"daytemp"`      // 白天温度
	NightTemp    string `json:"nighttemp"`    // 晚上温度
	DayWind      string `json:"daywind"`      // 白天风向
	NightWind    string `json:"nightwind"`    // 晚上风向
	DayPower     string `json:"daypower"`     // 白天风力
	NightPower   string `json:"nightpower"`   // 晚上风力
}

// GetWeatherLive 获取实时天气
func (c *Client) GetWeatherLive(ctx context.Context, request *WeatherRequest) (*WeatherResponse, error) {
	params := url.Values{}
	params.Add("city", request.City)
	params.Add("extensions", "base") // 获取实时天气

	data, err := c.Get(WeatherPath, params)
	if err != nil {
		return nil, err
	}

	var resp WeatherResponse
	if err := sonic.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetWeatherForecast 获取天气预报
func (c *Client) GetWeatherForecast(ctx context.Context, request *WeatherRequest) (*WeatherResponse, error) {
	params := url.Values{}
	params.Add("city", request.City)
	params.Add("extensions", "all") // 获取天气预报

	data, err := c.Get(WeatherPath, params)
	if err != nil {
		return nil, err
	}

	var resp WeatherResponse
	if err := sonic.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
