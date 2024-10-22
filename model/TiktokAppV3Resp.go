package model

// 注意：此go文件还未开发完毕。

type TiktokAppV3Params struct {
	Keyword     string `json:"keyword"`
	Offset      int    `json:"offset"`
	Count       int    `json:"count"`
	SortType    int    `json:"sort_type"`
	PublishTime int    `json:"publish_time"`
}

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
