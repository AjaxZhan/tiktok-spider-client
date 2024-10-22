package client

// 鉴于Tiktok的app端数据太多了，采取直接下载的方式，而不是保存为CSV

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"strconv"
	"tiktok-spider/conf"
	"tiktok-spider/model"
	"tiktok-spider/utils"
	"time"
)

type TiktokAppV3Client struct {
	params         model.TiktokAppV3Params // 爬虫参数
	baseUrl        string                  // api基础路径
	retry          int32                   // 重试次数
	baseFilePrefix string                  // 输出文件前缀

	httpClient *http.Client // http客户端
}

func NewTiktokV3Client(params model.TiktokAppV3Params) *TiktokAppV3Client {
	return &TiktokAppV3Client{
		params:  params,
		baseUrl: "https://api.tikhub.io/api/v1/tiktok/app/v3/fetch_video_search_result",
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseFilePrefix: "./output_tiktok_app/" + time.Now().Format("2006-01-02 15:01:05"),
		retry:          0,
	}
}

func (yc *TiktokAppV3Client) SearchVideoAndStore() {
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
		dataLen := len(resp.Data.Data)
		// store
		// 直接保存为json文件
		err = utils.SaveToJSON(resp, "/tiktok_app_v3/"+strconv.Itoa(i)+"-"+
			strconv.Itoa(dataLen)+".json")
		if err != nil {
			fmt.Println("警告：保存json错误:", err)
			break
		}
		total += dataLen
		fmt.Println("第", i, "轮爬虫结束，当前一共爬取", total, "条")
		// change params
		yc.params.Offset = resp.Data.Cursor
		i += 1
		// sleep
		time.Sleep(1 * time.Second)
	}
	fmt.Println("[end] 结束爬虫,数据一共有", total)
}

func (yc *TiktokAppV3Client) SearchVideo() (*model.TiktokAppV3Response, error) {
	// send http request
	url := yc.baseUrl + "?keyword=" + url2.PathEscape(yc.params.Keyword) +
		"&count=" + strconv.Itoa(yc.params.Count) +
		"&offset=" + strconv.Itoa(yc.params.Offset) +
		"&sort_type=" + strconv.Itoa(yc.params.SortType) +
		"&publish_time=" + strconv.Itoa(yc.params.PublishTime)

	fmt.Println("[SearchVideo] 爬虫请求路径为=", url)

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
	var resp model.TiktokAppV3Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("反序列化时错误：", err)
		fmt.Println("响应体：", string(body))
		return nil, err
	}
	fmt.Println("抓取数据成功，cursor = ", resp.Data.Cursor, "，正在检查...")

	// status check
	if resp.Code != 200 {
		//  保存当前爬虫参数
		fmt.Println("警告：api出现错误，状态码=", resp.Code)
		fmt.Println("保存爬虫参数：")
		_ = utils.SaveToJSON(yc.params, "./log/"+time.Now().Format("2006-01-02 15:01:05")+"-params-code_err"+".json")
		fmt.Println("保存最新结果:")
		_ = utils.SaveToJSON(resp.Data, "./log/"+time.Now().Format("2006-01-02 15:01:05")+"-crawl_data-code_err"+".json")
		return nil, errors.New("当前状态码不是200,code =" + string(resp.Code))
	}
	if len(resp.Data.Data) == 0 {
		fmt.Println("警告：爬取到的数据为0条，正在保存爬虫参数和最新结果")
		//  保存当前爬虫参数
		_ = utils.SaveToJSON(yc.params, "./log/"+time.Now().Format("2006-01-02 15:01:05")+"-params-num_err"+".json")
		_ = utils.SaveToJSON(resp.Data, "./log/"+time.Now().Format("2006-01-02 15:01:05")+"-crawl_data-num_err"+".json")
		if yc.retry == conf.AppConfig.MaxRetry {
			return nil, errors.New("爬到的数据为0并且超过最大重试次数")
		}
		yc.retry += 1
	}

	return &resp, nil
}
