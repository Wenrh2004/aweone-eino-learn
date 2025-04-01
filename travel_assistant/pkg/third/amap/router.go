package amap

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bytedance/sonic"
)

type RouterRequest struct {
	Origin      Point `json:"origin" jsonschema_description:"origin coordinate"`           // 起点坐标
	Destination Point `json:"destination" jsonschema_description:"destination coordinate"` // 终点坐标
}

type Point struct {
	Lon float64 `json:"lon" jsonschema_description:"lon"` // 经度
	Lat float64 `json:"lat" jsonschema_description:"lat"` // 纬度
}

// RouteResponse 路线规划响应
type RouteResponse struct {
	BaseResponse
	Route Route `json:"route"`
}

// Route 路线信息
type Route struct {
	Origin      string `json:"origin" jsonschema_description:"origin coordinate"`                   // 起点坐标
	Destination string `json:"destination" jsonschema_description:"destination coordinate"`         // 终点坐标
	Paths       []Path `json:"paths" jsonschema_description:"the list to the route search results"` // 路径规划信息列表
}

// Path 路径信息
type Path struct {
	Distance string `json:"distance" jsonschema_description:"path distance"`  // 路径距离
	Duration string `json:"duration" jsonschema_description:"estimated time"` // 预计耗时
	Steps    []Step `json:"steps" jsonschema_description:"step info"`         // 路径详情
}

// Step 路径详情
type Step struct {
	Instruction string `json:"instruction"` // 行走指示
	Road        string `json:"road"`        // 道路名称
	Distance    string `json:"distance"`    // 此段距离
	Duration    string `json:"duration"`    // 此段耗时
	Polyline    string `json:"polyline"`    // 此段坐标点串
}

// GetWalkingRoute 获取步行路线
func (c *Client) GetWalkingRoute(ctx context.Context, request *RouterRequest) (*RouteResponse, error) {
	params := url.Values{}
	params.Add("origin", fmt.Sprintf("%.6f,%.6f", request.Origin.Lon, request.Origin.Lat))
	params.Add("destination", fmt.Sprintf("%.6f,%.6f", request.Destination.Lon, request.Destination.Lat))

	data, err := c.Get(RouteWalkingPath, params)
	if err != nil {
		return nil, err
	}

	var resp RouteResponse
	if err := sonic.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetDrivingRoute 获取驾车路线
func (c *Client) GetDrivingRoute(ctx context.Context, request *RouterRequest) (*RouteResponse, error) {
	params := url.Values{}
	params.Add("origin", fmt.Sprintf("%.6f,%.6f", request.Origin.Lon, request.Origin.Lat))
	params.Add("destination", fmt.Sprintf("%.6f,%.6f", request.Destination.Lon, request.Destination.Lat))
	params.Add("extensions", "all") // 返回详细信息

	data, err := c.Get(RouteDrivingPath, params)
	if err != nil {
		return nil, err
	}

	var resp RouteResponse
	if err := sonic.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
