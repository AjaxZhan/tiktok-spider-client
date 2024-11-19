package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"tiktok-spider/client"
	"tiktok-spider/conf"
	"tiktok-spider/model"
	"tiktok-spider/repo"
	"time"
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

// 根据用户标签爬取
func crawlByUser() {
	secIds := [...]string{"MS4wLjABAAAA7iezvfSQ9ZQhe5oR4OhXA7bPSYQdNdWGUGdxJ9_Hec4uKxzZCzO_5rO-ykvRox7H",
		"MS4wLjABAAAAR9oFTAhfp2mBEevTOaVJK_HjzHs9t6GeTrbKklb6olLBOhupaj06A4mWYetb4JI4",
		"MS4wLjABAAAA-poLKkss7IfCj_R4w5aNsnCUnEMZ3YOV7ADSJq41plduH5TL3_Pc-Xfr_5H2TQ2L",
		"MS4wLjABAAAAoFWBcQQR7EhotYxTg-J6_BFHxgmJQPalgY4nN4nFIVfCus_RAG847z81hcCGjLiA",
		"MS4wLjABAAAAQR4rDH3IY7iexL_0pLNsZNXQtFoGX8tV8I3lQ2ek6kq-7PJskpM70OWnwAcnUI2P",
		"MS4wLjABAAAAX2flms30fOk_6fDtPTFWU7K0TtuXVtrq26wjvpW2SZaXQq2K1JHUQfHj4FNKKIil",
		"MS4wLjABAAAArSzT1mfGV3zJ5fbFFImiyw36RGvW05UGdEUilFeah8QQQ2Mu2DIi_G9r_1xfTRey",
		"MS4wLjABAAAAgd1ysR7DO8wQMfRj9ii6aA_RK-s03A0Jvwg-aq2ScZaWDn_nZeUHMw55FJZT2sgo",
		"MS4wLjABAAAABmCVX6gsimTUTmicvKBwL-Fvh6A4nF00nm9ACkwYTRZauIUnVdohHoUQVNzcshRC"}
	v3Client := client.NewTiktokV3Client(model.TiktokAppV3Params{}, model.TiktokTagParams{})
	v3Client.SearchVideoByUserSecIdsAndStoreBatch(secIds[:])
}

func testRand() {
	secIds := [...]string{"MS4wLjABAAAA7iezvfSQ9ZQhe5oR4OhXA7bPSYQdNdWGUGdxJ9_Hec4uKxzZCzO_5rO-ykvRox7H",
		"MS4wLjABAAAAR9oFTAhfp2mBEevTOaVJK_HjzHs9t6GeTrbKklb6olLBOhupaj06A4mWYetb4JI4",
		"MS4wLjABAAAA-poLKkss7IfCj_R4w5aNsnCUnEMZ3YOV7ADSJq41plduH5TL3_Pc-Xfr_5H2TQ2L",
		"MS4wLjABAAAAoFWBcQQR7EhotYxTg-J6_BFHxgmJQPalgY4nN4nFIVfCus_RAG847z81hcCGjLiA",
		"MS4wLjABAAAAQR4rDH3IY7iexL_0pLNsZNXQtFoGX8tV8I3lQ2ek6kq-7PJskpM70OWnwAcnUI2P",
		"MS4wLjABAAAAX2flms30fOk_6fDtPTFWU7K0TtuXVtrq26wjvpW2SZaXQq2K1JHUQfHj4FNKKIil",
		"MS4wLjABAAAArSzT1mfGV3zJ5fbFFImiyw36RGvW05UGdEUilFeah8QQQ2Mu2DIi_G9r_1xfTRey",
		"MS4wLjABAAAAgd1ysR7DO8wQMfRj9ii6aA_RK-s03A0Jvwg-aq2ScZaWDn_nZeUHMw55FJZT2sgo",
		"MS4wLjABAAAABmCVX6gsimTUTmicvKBwL-Fvh6A4nF00nm9ACkwYTRZauIUnVdohHoUQVNzcshRC"}
	var wt sync.WaitGroup
	for _, v := range secIds[:] {
		wt.Add(1)
		go func() {
			defer wt.Done()
			fmt.Println(v, strconv.FormatInt(time.Now().UnixMicro()+rand.Int63(), 10))
		}()
	}
	wt.Wait()
}

// 将TikTok-App爬取到的JSON视频数据整合到csv文件示例（需自行编写repo中的response结构体）
func saveToCSV() {
	repository := repo.NewTiktokAppV3Repository()

	// 匹配文件（需改成的json输出目录）
	dir := "./tiktok_app_v3_user"

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
		// 过滤出符合条件的文件（第一个数字 1-10000，第二个数字范围 1-100）
		if firstNum >= 1 && firstNum <= 10000 && secondNum >= 1 && secondNum <= 100 {
			fmt.Println("正在处理文件:", filename)
			// 执行你要对文件的处理逻辑
			_ = repository.DoneOne(file, "./tiktok_user.csv")
		}
	}
}

// 根据id将json保存到不同的csv
func saveDiffCSV() {
	repository := repo.NewTiktokAppV3Repository()

	// 匹配文件（需改成的 json 输出目录）
	dir := "./tiktok_app_v3_user"

	// 查找所有 .json 文件
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		log.Fatal(err)
	}

	// 定义正则表达式，匹配文件名格式：<id>_<number>-<number>.json
	re := regexp.MustCompile(`^(\d+)_\d+-\d+\.json`)

	// 创建一个 map 来存储文件的分组
	groupedFiles := make(map[string][]string)

	// 遍历匹配的文件
	for _, file := range files {
		// 从文件路径中提取文件名
		filename := filepath.Base(file)

		// 使用正则表达式提取 ID
		matches := re.FindStringSubmatch(filename)
		if len(matches) != 2 {
			// 如果文件名不符合格式，跳过
			continue
		}

		// 提取 ID
		id := matches[1]

		// 将文件分组到对应的 ID 下
		groupedFiles[id] = append(groupedFiles[id], file)
	}

	// 使用 WaitGroup 来管理 goroutine
	var wg sync.WaitGroup

	// 遍历分组后的文件
	for id, files := range groupedFiles {
		wg.Add(1)

		// 启动一个 goroutine 处理每个分组
		go func(id string, files []string) {
			defer wg.Done()

			csvFilePath := "./" + id + ".csv"

			// 处理分组内的文件
			for _, file := range files {
				log.Printf("正在处理文件: %s\n", file)
				err := repository.DoneOne(file, csvFilePath) // 替换为实际处理逻辑
				if err != nil {
					log.Printf("处理文件 %s 时出错: %v\n", file, err)
					continue
				}
			}
		}(id, files)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
}

// 过滤csv
func filterCSVByDesc(inputFilePath string) error {
	// 打开输入的 CSV 文件
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("无法打开文件 %s: %v", inputFilePath, err)
	}
	defer inputFile.Close()

	// 创建 CSV Reader
	reader := csv.NewReader(inputFile)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("读取 CSV 文件失败: %v", err)
	}

	// 检查是否有数据
	if len(records) == 0 {
		return fmt.Errorf("CSV 文件为空")
	}

	// 找到表头中 "Desc" 的索引
	header := records[0]
	descIndex := -1
	for i, col := range header {
		if strings.EqualFold(col, "Desc") {
			descIndex = i
			break
		}
	}
	if descIndex == -1 {
		return fmt.Errorf("CSV 文件中未找到表头 'Desc'")
	}

	// 保留表头
	filteredRecords := [][]string{header}

	// 遍历每一行数据，检查 "Desc" 列是否包含 "travel"
	for _, record := range records[1:] {
		if len(record) > descIndex && strings.Contains(strings.ToLower(record[descIndex]), "travel") {
			filteredRecords = append(filteredRecords, record)
		}
	}

	// 构造输出文件路径
	outputFilePath := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath)) + "_filter.csv"

	// 创建输出文件
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("无法创建文件 %s: %v", outputFilePath, err)
	}
	defer outputFile.Close()

	// 创建 CSV Writer
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// 写入过滤后的数据
	if err := writer.WriteAll(filteredRecords); err != nil {
		return fmt.Errorf("写入 CSV 文件失败: %v", err)
	}

	fmt.Printf("过滤后的数据已保存到 %s\n", outputFilePath)
	return nil
}

func filterTravel() error {
	folderPath := "./csv"
	// 检查文件夹是否存在
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("文件夹 %s 不存在", folderPath)
	}

	// 遍历文件夹下的所有文件
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("遍历文件 %s 时出错: %v", path, err)
		}

		// 检查是否为 CSV 文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".csv") {
			fmt.Printf("正在处理文件: %s\n", path)

			// 调用过滤函数
			if err := filterCSVByDesc(path); err != nil {
				fmt.Printf("处理文件 %s 时出错: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历文件夹时出错: %v", err)
	}

	fmt.Println("所有 CSV 文件已处理完毕。")
	return nil
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

	// 根据作者信息爬虫
	//crawlByUser()

	//数据保存到CSV
	//saveToCSV()
	//saveDiffCSV()

	// 过滤csv
	err := filterTravel()
	if err != nil {
		fmt.Println(err)
	}

	//testRand()

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
