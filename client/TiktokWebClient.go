package client

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
	"tiktok-spider/conf"
	"tiktok-spider/model"
	"tiktok-spider/utils"
	"time"
)

type TiktokWebClient struct {
	params         model.TiktokWebParams // 爬虫参数
	baseUrl        string                // api基础路径
	retry          int32                 // 重试次数
	baseFilePrefix string                // 输出文件前缀

	httpClient *http.Client // http客户端

	wt sync.WaitGroup
	mx sync.Mutex
}

func NewTiktokWebClient(params model.TiktokWebParams) *TiktokWebClient {
	return &TiktokWebClient{
		params:  params,
		baseUrl: "https://api.tikhub.io/api/v1/tiktok/web/fetch_search_video",
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseFilePrefix: time.UnixDate,
		retry:          0,
	}
}

// todo 这段逻辑大部分可以抽象，因为换了client后代码是一样的。
// 除了修改参数、统计爬虫数量根据接口而不同

func (yc *TiktokWebClient) SearchVideoAndStore() {
	i := 1
	total := 0
	fmt.Println("[begin] 开始爬虫")
	for {
		// search
		resp, err := yc.SearchVideo()
		if err != nil {
			fmt.Println("[SearchVideoAndStore] SearchVideo错误，信息：", err)
			break
		}
		// store
		yc.wt.Add(1)
		go func() {
			yc.mx.Lock()
			err := yc.writeCSV(resp)
			if err != nil {
				fmt.Println("[SearchVideoAndStore] writeCSV错误，信息：", err)
			}
			defer func() {
				yc.wt.Done()
				yc.mx.Unlock()
			}()
		}()
		total += len(resp.Data.ItemList)
		fmt.Println("第", i, "轮爬虫结束，当前一共爬取", total, "条")
		// change params
		yc.params.SearchId = resp.Data.LogPb.ImprId
		yc.params.Offset = resp.Data.Cursor
		i += 1
		// sleep
		time.Sleep(1 * time.Second)
	}
	yc.wt.Wait()
	fmt.Println("[end] 结束爬虫,数据一共有", total)
}

func (yc *TiktokWebClient) SearchVideo() (*model.TiktokWebResponse, error) {
	// send http request
	url := yc.baseUrl + "?keyword=" + url2.PathEscape(yc.params.Keyword) +
		"&count=" + yc.params.Count +
		"&offset=" + strconv.Itoa(int(yc.params.Offset)) +
		"&search_id=" + yc.params.SearchId +
		"&cookie=" + yc.params.Cookie
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求时错误：", err)
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+conf.AppConfig.Token)
	res, err := yc.httpClient.Do(req)
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

	// read http response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("读取响应时错误：", err)
		return nil, err
	}

	// resolve http response
	var resp model.TiktokWebResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("反序列化时错误：", err)
		fmt.Println("响应体：", string(body))
		return nil, err
	}
	fmt.Println("抓取数据成功，一共", len(resp.Data.ItemList), "条，正在检查...")

	// status check
	if resp.Code != 200 {
		//  保存当前爬虫参数
		fmt.Println("警告：api出现错误，状态码=", resp.Code)
		fmt.Println("保存爬虫参数：")
		_ = utils.SaveToJSON(resp.Params, time.UnixDate+"-params-code_err"+".json")
		fmt.Println("保存最新结果:")
		_ = utils.SaveToJSON(resp.Data, time.UnixDate+"-crawl_data-code_err"+".json")
		return nil, errors.New("当前状态码不是200,code =" + string(resp.Code))
	}
	if resp.Data.HasMore == 0 {
		fmt.Println("警告：爬取到的数据为0条，正在保存爬虫参数和最新结果")
		//  保存当前爬虫参数
		_ = utils.SaveToJSON(resp.Params, time.UnixDate+"-params-num_err"+".json")
		_ = utils.SaveToJSON(resp.Data, time.UnixDate+"-crawl_data-num_err"+".json")
		if yc.retry == conf.AppConfig.MaxRetry {
			return nil, errors.New("爬到的数据为0并且超过最大重试次数")
		}
		yc.retry += 1
	}

	return &resp, nil
}

func (yc *TiktokWebClient) writeCSV(data *model.TiktokWebResponse) error {
	outputName := yc.baseFilePrefix + "_output_tiktok_web.csv"
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
		"duration",
		"music_title",
		"music_url",
		"music_author",
	}
	file, err := os.OpenFile(outputName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		err = writer.Write(tabHeader)
		if err != nil {
			return err
		}
	}

	// write to csv
	for _, item := range data.Data.ItemList {
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
			strconv.Itoa(item.Video.Duration),
			item.Music.Title,
			item.Music.PlayUrl,
			item.Music.AuthorName,
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
