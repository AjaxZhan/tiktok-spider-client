package repo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

// 此文件将爬取到的数据持久化数据库或者csv文件中

// TiktokAuthorWebData 爬取作者相关信息
type tiktokAuthorWebData struct {
	Code int `json:"code"`
	Data struct {
		UserInfo struct {
			Stats struct {
				FollowerCount int `json:"followerCount"`
				HeartCount    int `json:"heartCount"`
				VideoCount    int `json:"videoCount"`
				DiggCount     int `json:"diggCount"`
			} `json:"stats"`
		} `json:"userInfo"`
		User struct {
			Verified bool `json:"verified"`
		} `json:"user"`
	} `json:"data"`
}

// TiktokAppV3WebData 从TikTok爬到的数据
type TiktokAppV3WebData struct {
	Author struct {
		Nickname  string `json:"nickname"`
		Region    string `json:"region"`
		Uid       string `json:"uid"`
		Signature string `json:"signature"`
		SecUid    string `json:"sec_uid"`
	} `json:"author"`
	Statistics struct {
		CollectCount  int `json:"collect_count"`
		CommentCount  int `json:"comment_count"`
		DiggCount     int `json:"digg_count"`
		DownloadCount int `json:"download_count"`
		PlayCount     int `json:"play_count"`
		ShareCount    int `json:"share_count"`
	}
	Music struct {
		Album   string `json:"album"`
		PlayUrl struct {
			Uri string `json:"uri"`
		} `json:"play_url"`
		Author string `json:"author"`
	} `json:"music"`
	AwemeId          string `json:"aweme_id"`
	Desc             string `json:"desc"`
	ShareUrl         string `json:"share_url"`
	ContentDescExtra []struct {
		HashtagName string `json:"hashtag_name"`
	} `json:"content_desc_extra"`
	CreateTime int `json:"create_time"`
	Video      struct {
		Duration int `json:"duration"`
	} `json:"video"`
}

type author struct {
	AuthorName       string // author.nickname
	Region           string // author.region
	AuthorId         string // author.uid
	IsVerify         bool
	AuthorSignature  int // author.signature
	AuthorFollower   int
	AuthorHeartCount int
	AuthorVideoCount int
	AuthorDiggCount  int
}

type VideoStatistics struct {
	CollectCount  int // statistics.collect_count
	CommentCount  int // statistics.comment_count
	DiggCount     int // statistics.digg_count
	DownloadCount int // statistics.download_count
	PlayCount     int // statistics.play_count
	ShareCount    int // statistics.share_count
}

type Music struct {
	MusicTitle  string // music.album
	MusicUrl    string // music.play_url.uri
	MusicAuthor string // music.author
}

// TiktokAppV3Vo 给csv文件写入的数据格式，采用组合的方式
type TiktokAppV3Vo struct {
	Id    string // aweme_id
	Desc  string // desc
	Cover string
	Url   string //share_url
	author
	TagList    string    // 拼接content_desc_extra.hashtag_name
	CreateTime time.Time // create_time
	Duration   int       // video.duration
	Music
	VideoStatistics
}

// CSV表头
var csvHeaders = []string{
	"Id",               // aweme_id 0
	"Desc",             // desc 1
	"Cover",            // cover 2
	"Url",              // share_url 3
	"AuthorName",       // author.nickname 4
	"Region",           // author.region 5
	"AuthorId",         // author.uid 6
	"IsVerify",         // author.is_verify 7
	"AuthorSignature",  // author.signature 8
	"AuthorFollower",   // author.follower_count 9
	"AuthorHeartCount", // author.total_favorited 10
	"AuthorVideoCount", // author.aweme_count 11
	"AuthorDiggCount",  // author.digg_count 12
	"TagList",          // 拼接content_desc_extra.hashtag_name 13
	"CreateTime",       // create_time 14
	"Duration",         // video.duration 15
	"MusicTitle",       // music.album 16
	"MusicUrl",         // music.play_url.uri 17
	"MusicAuthor",      // music.author 18
	"CollectCount",     // statistics.collect_count 19
	"CommentCount",     // statistics.comment_count 20
	"DiggCount",        // statistics.digg_count 21
	"DownloadCount",    // statistics.download_count 22
	"PlayCount",        // statistics.play_count 23
	"ShareCount",       // statistics.share_count 24
	"SecUid",           // author.sec_uid 25
}

// TiktokAppV3Repository 数据仓库
type TiktokAppV3Repository struct {
	csvPath           string
	authorInfoBaseUrl string
	httpClient        *http.Client
}

func NewTiktokAppV3Repository() *TiktokAppV3Repository {
	return &TiktokAppV3Repository{
		csvPath:           "./tiktok_app_v3_travel.csv",
		authorInfoBaseUrl: "https://douyin.wtf/api/tiktok/web/fetch_user_profile",
		httpClient: &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				MaxIdleConns:        200,              // 总的最大空闲连接数，适用于高并发环境
				MaxIdleConnsPerHost: 100,              // 每个主机的最大空闲连接数
				IdleConnTimeout:     30 * time.Second, // 空闲连接的超时时间
			},
		},
	}
}

// SaveToCSV 将数据写入CSV文件
func (repo *TiktokAppV3Repository) SaveToCSV(datas []TiktokAppV3WebData) error {

	// 创建CSV文件
	csvFile, err := os.OpenFile(repo.csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("无法创建CSV文件: %v", err)
	}
	defer csvFile.Close()

	// 创建CSV写入器
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// write header
	info, err := csvFile.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		err = writer.Write(csvHeaders)
		if err != nil {
			return err
		}
	}

	// 写入数据
	for i := 0; i < len(datas); i++ {
		var row []string
		data := datas[i]
		for _, header := range csvHeaders {
			switch header {
			case "Id":
				row = append(row, data.AwemeId)
			case "Desc":
				row = append(row, data.Desc)
			case "Cover":
				row = append(row, "") // Cover 不存在于结构体中，留空
			case "Url":
				row = append(row, data.ShareUrl)
			case "AuthorName":
				row = append(row, data.Author.Nickname)
			case "Region":
				row = append(row, data.Author.Region)
			case "AuthorId":
				row = append(row, data.Author.Uid)
			case "IsVerify":
				row = append(row, "") // IsVerify 不存在于结构体中，留空
			case "AuthorSignature":
				row = append(row, data.Author.Signature)
			case "AuthorFollower":
				row = append(row, "") // AuthorFollower 不存在于结构体中，留空
			case "AuthorHeartCount":
				row = append(row, "") // AuthorHeartCount 不存在于结构体中，留空
			case "AuthorVideoCount":
				row = append(row, "") // AuthorVideoCount 不存在于结构体中，留空
			case "AuthorDiggCount":
				row = append(row, "") // AuthorDiggCount 不存在于结构体中，留空
			case "TagList":
				extras := data.ContentDescExtra
				var tagList string
				for _, e := range extras {
					tagList += e.HashtagName + " "
				}
				row = append(row, tagList)
			case "CreateTime":
				// 将 Unix 时间戳转换为时间格式
				t := time.Unix(int64(data.CreateTime), 0)
				row = append(row, t.Format("2006-01-02 15:04:05"))
			case "Duration":
				row = append(row, strconv.Itoa(data.Video.Duration))
			case "MusicTitle":
				row = append(row, data.Music.Album)
			case "MusicUrl":
				row = append(row, data.Music.PlayUrl.Uri)
			case "MusicAuthor":
				row = append(row, data.Music.Author) // MusicAuthor 不存在于结构体中，留空
			case "CollectCount":
				row = append(row, strconv.Itoa(data.Statistics.CollectCount))
			case "CommentCount":
				row = append(row, strconv.Itoa(data.Statistics.CommentCount))
			case "DiggCount":
				row = append(row, strconv.Itoa(data.Statistics.DiggCount))
			case "DownloadCount":
				row = append(row, strconv.Itoa(data.Statistics.DownloadCount))
			case "PlayCount":
				row = append(row, strconv.Itoa(data.Statistics.PlayCount))
			case "ShareCount":
				row = append(row, strconv.Itoa(data.Statistics.ShareCount))
			case "SecUid":
				row = append(row, data.Author.SecUid)
			default:
				row = append(row, "") // 未知字段留空
			}
		}
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadFromJson 读取JSON文件
func (repo *TiktokAppV3Repository) LoadFromJson(path string) ([]TiktokAppV3WebData, error) {
	// 读取JSON文件
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取JSON文件: %v", err)
	}
	// 将JSON解析为map
	var data []TiktokAppV3WebData
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, fmt.Errorf("解析JSON数据出错: %v", err)
	}
	return data, nil
}

// DoneOne 执行一次读取JSON并写CSV操作
func (repo *TiktokAppV3Repository) DoneOne(path string) error {
	// json路径先根据本地数据写死
	data, err := repo.LoadFromJson(path)
	if err != nil {
		return err
	}
	// 写入csv
	err = repo.SaveToCSV(data)
	if err != nil {
		return err
	}

	fmt.Printf("文件:[%s],写入成功！", path)
	return nil
}

func (repo *TiktokAppV3Repository) UpdateAuthorData() error {
	// 打开CSV文件
	file, err := os.Open("./tiktok_app_v3_travel.csv")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return err
	}
	defer file.Close()

	// 读取CSV内容
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("读取CSV文件出错:", err)
		return err
	}

	// 遍历每一行并处理
	for i, row := range records {
		if i == 0 {
			// 跳过表头
			continue
		}

		authorName := row[4]
		fmt.Println("作者名字：", authorName)
		// 发送HTTP请求获取粉丝和点赞数量
		url := repo.authorInfoBaseUrl + "?uniqueId=" + url2.PathEscape(authorName)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("创建请求时错误：", err)
			continue
		}
		res, err := repo.httpClient.Do(request)
		if err != nil {
			fmt.Println("请求失败:", err)
			continue
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println("关闭io时错误：", err)
			}
		}(res.Body)
		// read http response
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("读取响应时错误：", err)
			continue
		}
		// unmarshal
		var authData tiktokAuthorWebData
		err = json.Unmarshal(body, &authData)
		if err != nil {
			fmt.Println("解析JSON出错:", err)
			continue
		}
		if authData.Code != 200 {
			fmt.Println("爬取作者信息时发生错误，状态码是：" + strconv.Itoa(authData.Code))
			continue
		}
		fmt.Println("爬取作者信息成功：")

		// 更新CSV
		row[7] = fmt.Sprintf("%v", authData.Data.User.Verified)
		row[9] = fmt.Sprintf("%v", authData.Data.UserInfo.Stats.FollowerCount)
		row[10] = fmt.Sprintf("%v", authData.Data.UserInfo.Stats.HeartCount)
		row[11] = fmt.Sprintf("%v", authData.Data.UserInfo.Stats.VideoCount)
		row[12] = fmt.Sprintf("%v", authData.Data.UserInfo.Stats.DiggCount)
		records[i] = row

		fmt.Println("更新第", i, "行成功！")
	}

	// 将修改后的内容重新写回CSV文件
	file, err = os.Create("tiktok_app_v3_author_update.csv")
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(records)
	if err != nil {
		fmt.Println("写入CSV文件出错:", err)
		return err
	}
	fmt.Println("[成功在CSV文件中添加作者信息]")
	return nil
}

// UpdateAuthorData2 使用并发编程的方式来更新作者信息
func (repo *TiktokAppV3Repository) UpdateAuthorData2() error {
	// Open CSV file for reading
	file, err := os.Open("./tiktok_app_v3_travel_updated.csv")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return err
	}
	defer file.Close()

	// Read CSV content
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("读取CSV文件出错:", err)
		return err
	}

	var wg sync.WaitGroup
	recordChan := make(chan []string)
	errorChan := make(chan error, 10000)
	// goroutine pool
	const maxGoroutines = 100 // 设置最大 Goroutine 数量
	guard := make(chan struct{}, maxGoroutines)
	var mu sync.Mutex // 用于写操作的互斥锁

	// 定时刷盘 Goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Second) // 每10秒刷盘一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				mu.Lock()
				fmt.Println("定时刷盘...")
				err := writeUpdatedRecordsToFile(records)
				if err != nil {
					fmt.Println("刷盘时发生错误:", err)
				}
				mu.Unlock()
			}
		}
	}()

	// Start a goroutine to update records as they are processed
	go func() {
		for record := range recordChan {
			for i, rec := range records {
				if rec[0] == record[0] {
					records[i] = record
					break
				}
			}
		}
	}()

	rateLimiter := time.Tick(2000 * time.Millisecond) // 控制请求间隔

	// Process records concurrently
	for i, row := range records {
		if i == 0 || row[4] == "" || row[9] != "" {
			// Skip header or rows without authorName or updated record
			continue
		}
		<-rateLimiter // 控制请求速率
		wg.Add(1)
		guard <- struct{}{}
		go func(row []string) {
			defer wg.Done()
			defer func() { <-guard }()
			authorName := row[4]
			secUid := row[25]
			fmt.Printf("作者名字：%s，sec_uid=%s\n", authorName, secUid)
			url := repo.authorInfoBaseUrl + "?secUid=" + url2.QueryEscape(secUid)
			request, err := http.NewRequest("GET", url, nil)
			request.Header.Add("Accept", "application/json")
			if err != nil {
				fmt.Println("创建请求时错误：", err)
				errorChan <- err
				return
			}
			res, err := repo.httpClient.Do(request)
			if err != nil {
				fmt.Println("请求失败:", err)
				errorChan <- err
				return
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("读取响应时错误：", err)
				errorChan <- err
				return
			}

			var authData tiktokAuthorWebData
			err = json.Unmarshal(body, &authData)
			if err != nil {
				fmt.Printf("解析JSON出错:%v\n，请求体%v,作者名字：%s，secuid=%s\n", err, string(body), authorName, secUid)
				errorChan <- err
				return
			}
			if authData.Code != 200 {
				fmt.Println("爬取作者信息时发生错误，状态码是：" + strconv.Itoa(authData.Code))
				return
			}

			// Update row with new data
			row[7] = fmt.Sprintf("%v", authData.Data.User.Verified)
			row[9] = fmt.Sprintf("%v", authData.Data.UserInfo.Stats.FollowerCount)
			row[10] = fmt.Sprintf("%v", authData.Data.UserInfo.Stats.HeartCount)
			row[11] = fmt.Sprintf("%v", authData.Data.UserInfo.Stats.VideoCount)
			row[12] = fmt.Sprintf("%v", authData.Data.UserInfo.Stats.DiggCount)

			recordChan <- row
			fmt.Printf("更新成功！作者名字：%s\n", authorName)
		}(row)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(recordChan)
		close(errorChan)
	}()

	// Check for errors during processing
	for err := range errorChan {
		if err != nil {
			fmt.Println("发生错误：", err)
		}
	}

	// Write updated records to CSV
	writeFile, err := os.Create("./tiktok_app_v3_travel_updated.csv")
	if err != nil {
		fmt.Println("无法创建更新的文件:", err)
		return err
	}
	defer writeFile.Close()

	writer := csv.NewWriter(writeFile)
	err = writer.WriteAll(records)
	if err != nil {
		fmt.Println("写入CSV文件时出错:", err)
		return err
	}
	writer.Flush()

	fmt.Println("CSV文件更新完毕")
	return nil
}

// writeUpdatedRecordsToFile 将记录写入文件
func writeUpdatedRecordsToFile(records [][]string) error {
	writeFile, err := os.Create("./tiktok_app_v3_travel_updated.csv")
	if err != nil {
		return err
	}
	defer writeFile.Close()

	writer := csv.NewWriter(writeFile)
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
