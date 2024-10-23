package repo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// 此文件将爬取到的数据持久化数据库或者csv文件中

type TiktokAppV3WebData struct {
	Author struct {
		Nickname  string `json:"nickname"`
		Region    string `json:"region"`
		Uid       string `json:"uid"`
		Signature string `json:"signature"`
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

// 数据模型
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

// 表头
var csvHeaders = []string{
	"Id",               // aweme_id
	"Desc",             // desc
	"Cover",            // cover
	"Url",              // share_url
	"AuthorName",       // author.nickname
	"Region",           // author.region
	"AuthorId",         // author.uid
	"IsVerify",         // author.is_verify
	"AuthorSignature",  // author.signature
	"AuthorFollower",   // author.follower_count
	"AuthorHeartCount", // author.total_favorited
	"AuthorVideoCount", // author.aweme_count
	"AuthorDiggCount",  // author.digg_count
	"TagList",          // 拼接content_desc_extra.hashtag_name
	"CreateTime",       // create_time
	"Duration",         // video.duration
	"MusicTitle",       // music.album
	"MusicUrl",         // music.play_url.uri
	"MusicAuthor",      // music.author
	"CollectCount",     // statistics.collect_count
	"CommentCount",     // statistics.comment_count
	"DiggCount",        // statistics.digg_count
	"DownloadCount",    // statistics.download_count
	"PlayCount",        // statistics.play_count
	"ShareCount",       // statistics.share_count
}

// TiktokAppV3Repository 数据仓库
type TiktokAppV3Repository struct {
	csvPath string
}

func NewTiktokAppV3Repository() *TiktokAppV3Repository {
	return &TiktokAppV3Repository{
		csvPath: "./tiktok_app_v3.csv",
	}
}

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

func (repo *TiktokAppV3Repository) DoneBatch() error {
	// json路径先根据本地数据写死
	return nil
}
