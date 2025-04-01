package amap

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bytedance/sonic"
)

const (
	// BaseURL 高德API基础URL
	BaseURL = "https://restapi.amap.com"

	// API路径
	GeocodePath      = "/v3/geocode/geo"         // 地理编码
	ReGeocodePath    = "/v3/geocode/regeo"       // 逆地理编码
	WeatherPath      = "/v3/weather/weatherInfo" // 天气查询
	POISearchPath    = "/v3/place/text"          // POI搜索
	RouteWalkingPath = "/v3/direction/walking"   // 步行路线规划
	RouteDrivingPath = "/v3/direction/driving"   // 驾车路线规划
)

// Client 高德API客户端
type Client struct {
	Key        string
	HTTPClient *http.Client
}

// NewClient 创建新的高德API客户端
func NewClient(key string, timeout time.Duration) *Client {
	return &Client{
		Key: key,
		HTTPClient: &http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

// Get 发送GET请求到高德API
func (c *Client) Get(path string, params url.Values) ([]byte, error) {
	// 添加API密钥
	params.Add("key", c.Key)

	// 构建完整URL
	fullURL := fmt.Sprintf("%s%s?%s", BaseURL, path, params.Encode())

	// 创建请求
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求高德API失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("高德API返回非200状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}

// BaseResponse 通用响应结构
type BaseResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	InfoCode string `json:"infocode"`
}

// CheckResponse 检查API响应是否成功
func (r *BaseResponse) CheckResponse(data []byte) error {
	if err := sonic.Unmarshal(data, r); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if r.Status != "1" {
		return fmt.Errorf("API调用失败: %s (代码: %s)", r.Info, r.InfoCode)
	}

	return nil
}
