package model

type YoutubeResponse struct {
	Code   int           `json:"code"`
	Router string        `json:"router"`
	Params YoutubeParams `json:"params"`
	Data   YoutubeData   `json:"data"`
}

type YoutubeParams struct {
	SearchQuery       string `json:"search_query"`
	LanguageCode      string `json:"language_code"`
	OrderBy           string `json:"order_by"`
	CountryCode       string `json:"country_code"`
	ContinuationToken string `json:"continuation_token"`
}

type YoutubeData struct {
	NumberOfVideos    int            `json:"number_of_videos"`
	Query             string         `json:"query"`
	Country           string         `json:"country"`
	Lang              string         `json:"lang"`
	Timezone          string         `json:"timezone"`
	ContinuationToken string         `json:"continuation_token"`
	Videos            []YoutubeVideo `json:"videos"`
}
type YoutubeVideo struct {
	VideoID       string      `json:"video_id"`
	Title         string      `json:"title"`
	Author        string      `json:"author"`
	NumberOfViews int         `json:"number_of_views"`
	VideoLength   string      `json:"video_length"`
	Description   string      `json:"description"`
	IsLiveContent interface{} `json:"is_live_content"` // 由于可能为空，使用 interface{}
	PublishedTime string      `json:"published_time"`
	ChannelID     string      `json:"channel_id"`
	Category      interface{} `json:"category"` // 由于可能为空，使用 interface{}
	Type          string      `json:"type"`
	Keywords      []string    `json:"keywords"`
	Thumbnails    []Thumbnail `json:"thumbnails"`
}

var YoutubeVideoCSVHeaders = []string{
	"VideoID",
	"Title",
	"Author",
	"NumberOfViews",
	"VideoLength",
	"Description",
	"IsLiveContent",
	"PublishedTime",
	"ChannelID",
	"Category",
	"Type",
	"Keywords",
	"Thumbnails",
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
