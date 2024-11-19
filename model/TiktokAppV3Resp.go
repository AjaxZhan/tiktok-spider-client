package model

// TiktokAppV3Params APP搜索视频参数
type TiktokAppV3Params struct {
	Keyword     string `json:"keyword"`
	Offset      int    `json:"offset"`
	Count       int    `json:"count"`
	SortType    int    `json:"sort_type"`
	PublishTime int    `json:"publish_time"`
}

// TiktokAppV3Response APP应答
type TiktokAppV3Response struct {
	Code   int    `json:"code"`
	Router string `json:"router"`
	Data   data   `json:"data"`
}
type data struct {
	Cursor     int           `json:"cursor"`
	StatusCode int           `json:"status_code"`
	Data       []interface{} `json:"data"`
}

// TiktokTagParams 标签搜索请求
type TiktokTagParams struct {
	ChId   string `json:"ch_id"`
	Cursor int    `json:"cursor"`
	Count  int    `json:"count"`
}

// TiktokTagResponse 标签搜索应答
type TiktokTagResponse struct {
	Code   int    `json:"code"`
	Router string `json:"router"`
	Data   struct {
		AwemeList []interface{} `json:"aweme_list"`
		Cursor    int           `json:"cursor"`
	} `json:"data"`
}

// TiktokSearchVideoByUserResponse 标签搜索应答
type TiktokSearchVideoByUserResponse struct {
	Code   int    `json:"code"`
	Router string `json:"router"`
	Data   struct {
		AwemeList []interface{} `json:"aweme_list"`
		MaxCursor int           `json:"max_cursor"`
		HasMore   int           `json:"has_more"`
	} `json:"data"`
}
