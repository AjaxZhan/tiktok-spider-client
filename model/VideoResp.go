package model

type VideoResp struct {
	Code   int32  `json:"code"`
	Router string `json:"router"`
	params `json:"params"`
	data   `json:"data"`
}

type params struct {
	Keyword  string `json:"keyword"`
	Count    string `json:"count"`
	Offset   string `json:"offset"`
	SearchId string `json:"search_id"`
}

type data struct {
	StatusCode int32  `json:"status_code"`
	ItemList   []Item `json:"item_list"`
	HasMore    int32  `json:"has_more"`
	Cursor     int32  `json:"cursor"`
	extra      `json:"extra"`
	LogPb      struct {
		ImprId string `json:"impr_id"`
	} `json:"log_pb"`
	Backtrace string `json:"backtrace"`
}

type extra struct {
	Now             int64         `json:"now"`
	Logid           string        `json:"logid"`
	FatalItemIds    []interface{} `json:"fatal_item_ids"`
	SearchRequestId string        `json:"search_request_id"`
	ApiDebugInfo    string        `json:"api_debug_info"`
}
