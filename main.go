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

func crawlTiktokAppV3() {
	params := model.TiktokAppV3Params{
		Keyword:     "china travel",
		Offset:      660,
		Count:       20,
		SortType:    0,
		PublishTime: 0,
	}
	tagParams := model.TiktokTagParams{
		ChId:   "7884", // #travel
		Cursor: 3340,
		Count:  10,
	}
	v3Client := client.NewTiktokV3Client(params, tagParams)
	//v3Client.SearchVideoAndStore()
	v3Client.SearchVideoByTagAndStore()
}

func saveToCSV() {
	repository := repo.NewTiktokAppV3Repository()
	// 先处理tag系列的，第一批，编号744
	//err := repository.DoneOne("./tiktok_app_v3/172961426355374_1-10.json")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 匹配文件
	dir := "./tiktok_app_v3_2"

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
		// 过滤出符合条件的文件（第一个数字 1-66，第二个数字范围 1-10）
		if firstNum >= 1 && firstNum <= 200 && secondNum >= 1 && secondNum <= 10 {
			fmt.Println("正在处理文件:", filename)
			// 执行你要对文件的处理逻辑
			_ = repository.DoneOne(file)
		}
	}
}

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
	UpdateAuthorData()
}

// 加载配置
func loadConf() {
	err := conf.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}
