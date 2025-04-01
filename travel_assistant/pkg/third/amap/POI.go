package amap

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bytedance/sonic"
)

type POISearchRequest struct {
	Keyword string `json:"keyword" jsonschema_description:"Search keywords for points of interest"`
	City    string `json:"city" jsonschema_description:"City name"`

	Page     int `json:"page" jsonschema_description:"Page number"`
	PageSize int `json:"pageSize" jsonschema_description:"Number of items per page"`
}

type Suggestion struct {
	Keywords []interface{} `json:"keywords,omitempty"`
	Cities   []interface{} `json:"cities,omitempty"`
}

type Photo struct {
	Title interface{} `json:"title,omitempty"`
	Url   string      `json:"url,omitempty"`
}

type BizExt struct {
	Cost         interface{} `json:"cost,omitempty"`
	Rating       interface{} `json:"rating,omitempty"`
	MealOrdering string      `json:"meal_ordering,omitempty"`
}

type IndoorData struct {
	Cmsid     interface{}   `json:"cmsid,omitempty"`
	Truefloor []interface{} `json:"truefloor,omitempty"`
	Cpid      interface{}   `json:"cpid,omitempty"`
	Floor     []interface{} `json:"floor,omitempty"`
}

type POI struct {
	Parent       interface{}   `json:"parent,omitempty"`
	Distance     []interface{} `json:"distance,omitempty"`
	Pcode        string        `json:"pcode,omitempty"`
	Importance   []interface{} `json:"importance,omitempty"`
	BizExt       BizExt        `json:"biz_ext,omitempty"`
	Recommend    string        `json:"recommend,omitempty"`
	Type         string        `json:"type,omitempty"`
	Photos       []Photo       `json:"photos,omitempty"`
	DiscountNum  string        `json:"discount_num,omitempty"`
	Gridcode     string        `json:"gridcode,omitempty"`
	Typecode     string        `json:"typecode,omitempty"`
	Shopinfo     string        `json:"shopinfo,omitempty"`
	Poiweight    []interface{} `json:"poiweight,omitempty"`
	Citycode     string        `json:"citycode,omitempty"`
	Adname       string        `json:"adname,omitempty"`
	Children     []interface{} `json:"children,omitempty"`
	Alias        interface{}   `json:"alias,omitempty"`
	Tel          interface{}   `json:"tel,omitempty"`
	Id           string        `json:"id,omitempty"`
	Tag          interface{}   `json:"tag"`
	Event        []interface{} `json:"event"`
	EntrLocation interface{}   `json:"entr_location"`
	IndoorMap    string        `json:"indoor_map"`
	Email        []interface{} `json:"email"`
	Timestamp    string        `json:"timestamp"`
	Website      interface{}   `json:"website"`
	Address      string        `json:"address"`
	Adcode       string        `json:"adcode"`
	Pname        string        `json:"pname"`
	BizType      interface{}   `json:"biz_type"`
	Cityname     string        `json:"cityname"`
	Postcode     []interface{} `json:"postcode"`
	Match        string        `json:"match"`
	BusinessArea interface{}   `json:"business_area"`
	IndoorData   IndoorData    `json:"indoor_data"`
	Childtype    interface{}   `json:"childtype"`
	ExitLocation []interface{} `json:"exit_location"`
	Name         string        `json:"name"`
	Location     string        `json:"location"`
	Shopid       []interface{} `json:"shopid"`
	NaviPoiid    interface{}   `json:"navi_poiid"`
	GroupbuyNum  string        `json:"groupbuy_num"`
}

// POISearchResponse POI搜索响应
type POISearchResponse struct {
	*BaseResponse
	Suggestion Suggestion `json:"suggestion"`
	Count      string     `json:"count"`
	Pois       []POI      `json:"pois"`
}

// SearchPOI 搜索兴趣点
func (c *Client) SearchPOI(ctx context.Context, request *POISearchRequest) (*POISearchResponse, error) {
	params := url.Values{}
	params.Add("keywords", request.Keyword)
	params.Add("city", request.City)
	params.Add("offset", strconv.Itoa(request.PageSize))
	params.Add("page", strconv.Itoa(request.Page))
	params.Add("extensions", "all") // 返回详细信息

	data, err := c.Get(POISearchPath, params)
	if err != nil {
		return nil, err
	}

	var resp POISearchResponse
	if err := sonic.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SearchNearbyPOI 搜索附近兴趣点
func (c *Client) SearchNearbyPOI(keyword string, longitude, latitude float64, radius int) (*POISearchResponse, error) {
	params := url.Values{}
	params.Add("keywords", keyword)
	params.Add("location", fmt.Sprintf("%.6f,%.6f", longitude, latitude))
	params.Add("radius", strconv.Itoa(radius))
	params.Add("extensions", "all") // 返回详细信息

	data, err := c.Get(POISearchPath, params)
	if err != nil {
		return nil, err
	}

	var resp POISearchResponse
	if err := sonic.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
