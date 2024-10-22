package client

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"tiktok-spider/conf"
	"tiktok-spider/model"
	"tiktok-spider/utils"
	"time"
)

type YoutubeClient struct {
	params         model.YoutubeParams // 爬虫参数
	baseUrl        string              // api基础路径
	retry          int32               // 重试次数
	baseFilePrefix string              // 输出文件前缀
	httpClient     *http.Client        // http客户端

	wt sync.WaitGroup
	mx sync.Mutex
}

func NewYoutubeClient(params model.YoutubeParams) *YoutubeClient {
	return &YoutubeClient{
		params:  params,
		baseUrl: "https://api.tikhub.io/api/v1/youtube/web/search_video",
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseFilePrefix: "./output/" + time.Now().Format("2006-01-02 15:01:05"),
		retry:          0,
	}
}

func (yc *YoutubeClient) SearchVideoAndStore() {
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
		total += resp.Data.NumberOfVideos
		fmt.Println("第", i, "轮爬虫结束，当前一共爬取", total, "条")
		// change params
		yc.params.ContinuationToken = resp.Data.ContinuationToken
		i += 1
		// sleep
		time.Sleep(1 * time.Second)
	}
	yc.wt.Wait()
	fmt.Println("[end] 结束爬虫,数据一共有", total)
}

func (yc *YoutubeClient) SearchVideo() (*model.YoutubeResponse, error) {
	// send http request
	url := yc.baseUrl + "?search_query=" + yc.params.SearchQuery +
		"&language_code=" + yc.params.LanguageCode +
		"&order_by=" + yc.params.OrderBy +
		"&country_code=" + yc.params.CountryCode
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
	var resp model.YoutubeResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("反序列化时错误：", err)
		fmt.Println("响应体：", string(body))
		return nil, err
	}
	fmt.Println("抓取数据成功，一共", resp.Data.NumberOfVideos, "条，正在检查...")

	// status check
	if resp.Code != 200 {
		//  保存当前爬虫参数
		fmt.Println("警告：api出现错误，状态码=", resp.Code)
		fmt.Println("保存爬虫参数：")
		_ = utils.SaveToJSON(resp.Params, "./log/"+time.Now().Format("2006-01-02 15:01:05")+"-params-code_err"+".json")
		fmt.Println("保存最新结果:")
		_ = utils.SaveToJSON(resp.Data, "./log/"+time.Now().Format("2006-01-02 15:01:05")+"-crawl_data-code_err"+".json")
		return nil, errors.New("当前状态码不是200,code =" + string(resp.Code))
	}
	if resp.Data.NumberOfVideos == 0 {
		fmt.Println("警告：爬取到的数据为0条，正在保存爬虫参数和最新结果")
		//  保存当前爬虫参数
		_ = utils.SaveToJSON(resp.Params, "./log/"+time.Now().Format("2006-01-02 15:01:05")+"-params-num_err"+".json")
		_ = utils.SaveToJSON(resp.Data, "./log/"+time.Now().Format("2006-01-02 15:01:05")+"-crawl_data-num_err"+".json")
		if yc.retry == conf.AppConfig.MaxRetry {
			return nil, errors.New("爬到的数据为0并且超过最大重试次数")
		}
		yc.retry += 1
	}

	return &resp, nil
}

func (yc *YoutubeClient) writeCSV(data *model.YoutubeResponse) error {
	outputName := yc.baseFilePrefix + "_output_youtube.csv"

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
		err = writer.Write(model.YoutubeVideoCSVHeaders)
		if err != nil {
			return err
		}
	}

	// write to csv
	for _, item := range data.Data.Videos {
		_ = writer.Write([]string{
			item.VideoID,
			item.Title,
			item.Author,
			strconv.Itoa(item.NumberOfViews),
			item.VideoLength,
			item.Description,
			"",
			//item.IsLiveContent.(string),
			item.PublishedTime,
			item.ChannelID,
			//item.Category.(string),
			"",
			item.Type,
			strings.Join(item.Keywords, " "),
			"",
		})
	}
	fmt.Println("爬虫成功，数据已经写入csv文件")
	return nil
}
