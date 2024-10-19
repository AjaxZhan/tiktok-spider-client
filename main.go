package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"sync"
	"tiktok-spider/model"
	"time"
)

const (
	token = "6psPGzlj7kobkTAl6zGQyVoY5BoriJmC7olf5XLB6k954GQ5Yt9Z+y2JLg=="
	stop  = 10000
)

type params struct {
	keyword   string
	count     int32
	offset    int32
	search_id string
	cookie    string
}

func main() {
	var crawlCount int
	var wt sync.WaitGroup
	// user query
	params := &params{
		keyword:   url2.PathEscape("#china travel"),
		count:     20,
		offset:    0,
		search_id: "",
		cookie:    "",
	}

	for crawlCount = 0; crawlCount < stop; {
		// get data
		resp, err := getData(params)
		if err != nil {
			fmt.Println(err)
			return
		}
		// output csv
		wt.Add(1)
		go func() {
			err := writeCSV(resp.ItemList)
			if err != nil {
				fmt.Println("写入csv时错误：", err)
			}
			defer wt.Done()
		}()
		crawlCount += len(resp.ItemList)
		fmt.Println("当前已爬取数据：", crawlCount, "条")
		time.Sleep(4 * time.Second)
	}
	wt.Wait()
}

func getData(p *params) (*model.VideoResp, error) {

	url := "https://api.tikhub.io/api/v1/tiktok/web/fetch_search_video?keyword=" +
		p.keyword + "&count=" + strconv.Itoa(int(p.count)) +
		"&offset=" + strconv.Itoa(int(p.offset))
	// user client
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求时错误：", err)
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时错误：", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭io时错误：", err)
		}
	}(res.Body)

	// resp:read
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("读取响应时错误：", err)
		return nil, err
	}

	// resp:unmarshal
	var resp model.VideoResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("反序列化时错误：", err)
		return nil, err
	}
	fmt.Println("抓取数据成功，一共", len(resp.ItemList), "条，正在检查...")
	if resp.Code != 200 {
		return nil, errors.New("当前状态码不是200")
	}
	if resp.HasMore != 1 {
		return nil, errors.New("没有更多数据了")
	}

	return &resp, nil
}

// write a batch of data to csv
func writeCSV(itemList []model.Item) error {
	tabHeader := []string{
		"id",
		"desc",
		"cover",
		"video_addr", // 自己合成
		"author_id",
		"author_name",
		"digg_count",
		"share_count",
		"comment_count",
		"play_count",
		"collect_count",
		"tag_list",
		"is_verify",
		"author_signature",
		"author_follower",
		"author_heart_count",
		"author_video_count",
		"author_digg_count",
		"create_time",
	}
	file, err := os.Create("output.csv")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	// init csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write header
	err = writer.Write(tabHeader)
	if err != nil {
		return err
	}

	// write to csv
	for _, item := range itemList {
		_ = writer.Write([]string{
			item.ID,
			item.Desc,
			item.Video.Cover,
			formatVideoURL(item.Author.Nickname, item.Video.ID),
			item.Author.ID,
			item.Author.Nickname,
			strconv.Itoa(item.Stats.DiggCount),
			strconv.Itoa(item.Stats.ShareCount),
			strconv.Itoa(item.Stats.CommentCount),
			strconv.Itoa(item.Stats.PlayCount),
			strconv.Itoa(item.Stats.CollectCount),
			getTagList(item),
			strconv.FormatBool(item.Author.Verified),
			item.Author.Signature,
			strconv.Itoa(item.AuthorStats.FollowerCount),
			strconv.Itoa(item.AuthorStats.HeartCount),
			strconv.Itoa(item.AuthorStats.VideoCount),
			strconv.Itoa(item.AuthorStats.DiggCount),
			time.Unix(item.CreateTime, 0).Format("2006-01-02 15:04:05"),
		})
	}
	fmt.Println("爬虫成功，数据已经写入csv文件")
	return nil
}
func formatVideoURL(authorName string, videoId string) string {
	return "https://www.tiktok.com/@" + authorName + "/video/" + videoId
}
func getTagList(item model.Item) string {
	a := ""
	for _, extra := range item.TextExtra {
		a += extra.HashtagName
		a += " "
	}
	return a
}
