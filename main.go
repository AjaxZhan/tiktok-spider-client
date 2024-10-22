package main

import (
	"log"
	"tiktok-spider/client"
	"tiktok-spider/conf"
	"tiktok-spider/model"
)

func crawlYoutube() {
	// 启动爬虫客户端
	params := model.YoutubeParams{
		SearchQuery:  "china travel",
		LanguageCode: "en",
		OrderBy:      "this_month",
		CountryCode:  "us",
	}
	youtubeClient := client.NewYoutubeClient(params)
	youtubeClient.SearchVideoAndStore()
}

func crawlTiktokWeb() {
	params := model.TiktokWebParamsSend{
		Keyword: "China travel",
		Offset:  0,
		Count:   1,
	}
	tiktokWebClient := client.NewTiktokWebClient(params)
	tiktokWebClient.SearchVideoAndStore()
}

func main() {
	// 加载配置
	loadConf()
	// 开启爬虫
	//crawlYoutube()
	crawlTiktokWeb()
}

// 加载配置
func loadConf() {
	err := conf.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}
