package model

// 注意：此go文件还未开发完毕。

type TiktokAppV3Params struct {
	Keyword     string `json:"keyword"`
	Offset      string `json:"offset"`
	Count       string `json:"count"`
	SortType    string `json:"sort_type"`
	PublishTime string `json:"publish_time"`
}

type TiktokAppV3Response struct {
	Code   int               `json:"code"`
	Router string            `json:"router"`
	Params TiktokAppV3Params `json:"params"`
	Data   TiktokAppV3Data   `json:"data"`
}

type TiktokAppV3Data struct {
	cursor int
}
type TiktokData struct {
}

func main() {
	// 示例使用
}
