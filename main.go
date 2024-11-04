package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"tiktok-spider/client"
	"tiktok-spider/conf"
	"tiktok-spider/model"
	"tiktok-spider/repo"
)

// 爬取YouTube视频示例
func crawlYoutube() {
	// 准备爬虫参数
	params := model.YoutubeParams{
		SearchQuery:  "china travel",
		LanguageCode: "en",
		OrderBy:      "this_month",
		CountryCode:  "us",
	}
	// 准备客户端
	youtubeClient := client.NewYoutubeClient(params)
	youtubeClient.SearchVideoAndStore()
}

// 基于TikTok的Web段接口爬虫示例（视频搜索接口）
func crawlTiktokWeb() {
	params := model.TiktokWebParamsSend{
		Keyword: "China travel",
		Offset:  0,
		Count:   1,
	}
	tiktokWebClient := client.NewTiktokWebClient(params)
	tiktokWebClient.SearchVideoAndStore()
}

// 基于TikTok的App接口爬虫示例（视频搜索接口｜标签搜索接口）
func crawlTiktokAppV3() {
	// 直接进行视频搜索（似乎视频量有限）
	params := model.TiktokAppV3Params{
		Keyword:     "china travel",
		Offset:      660,
		Count:       20,
		SortType:    0,
		PublishTime: 0,
	}
	// 根据标签搜索视频（视频量最大为5050）
	tagParams := model.TiktokTagParams{
		//ChId:   "7884", // #travel
		ChId:   "74640468", // #chinatravel
		Cursor: 0,
		Count:  10,
	}
	v3Client := client.NewTiktokV3Client(params, tagParams)
	//v3Client.SearchVideoAndStore()
	v3Client.SearchVideoByTagAndStore()
}

// 将TikTok-App爬取到的JSON视频数据整合到csv文件示例（需自行编写repo中的response结构体）
func saveToCSV() {
	repository := repo.NewTiktokAppV3Repository()

	// 匹配文件（需改成的json输出目录）
	dir := "./tiktok_app_v3"

	// 查找所有 .json 文件
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		log.Fatal(err)
	}

	// 定义正则表达式，匹配文件名格式：<number>_<number>.json
	re := regexp.MustCompile(`\d+_(\d+)-(\d+)\.json`)

	// 遍历匹配的文件
	for _, file := range files {
		// 从文件路径中提取文件名
		filename := filepath.Base(file)
		//fmt.Println(filename)
		// 使用正则表达式提取数字
		matches := re.FindStringSubmatch(filename)
		if len(matches) != 3 {
			// 如果文件名不符合格式，跳过
			continue
		}
		// 提取第一个数字
		firstNum, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Println("Error converting first number:", err)
			continue
		}

		// 提取第二个数字
		secondNum, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Println("Error converting second number:", err)
			continue
		}
		// 过滤出符合条件的文件（第一个数字 1-1000，第二个数字范围 1-10）
		if firstNum >= 1 && firstNum <= 1000 && secondNum >= 1 && secondNum <= 10 {
			fmt.Println("正在处理文件:", filename)
			// 执行你要对文件的处理逻辑
			_ = repository.DoneOne(file)
		}
	}
}

// UpdateAuthorData 额外补充的方法实示例，使用免费接口更新标签接口没有的作者信息
func UpdateAuthorData() {
	r := repo.NewTiktokAppV3Repository()
	err := r.UpdateAuthorData2()
	if err != nil {
		fmt.Println("更新作者信息时发生错误！", err.Error())
	}
}

func main() {
	// 加载配置
	loadConf()

	// 开启爬虫
	//crawlYoutube()
	//crawlTiktokWeb()
	//crawlTiktokAppV3()

	//数据保存到CSV
	//saveToCSV()

	// 更新作者信息
	//UpdateAuthorData()
}

// 加载配置
func loadConf() {
	err := conf.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}
