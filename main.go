package main

import (
	"log"
	"tiktok-spider/client"
	"tiktok-spider/conf"
	"tiktok-spider/model"
)

//const (
//	token = "6psPGzlj7kobkTAl6zGQyVoY5BoriJmC7olf5XLB6k954GQ5Yt9Z+y2JLg=="
//	ck    = "tt_chain_token=6rmHtivV+TvBsGBL0+YSZA==; tt_csrf_token=WOpXMdiJ--Hnh1k7ibcnOYziQVJVddXxYQNc; ak_bmsc=10DFD09720AEBAD644FF5A6F7979A960~000000000000000000000000000000~YAAQ3WLBF6fxi6CSAQAAnk3yqRmHId2doQk2gF+UaHknZSqJnJF9bgu81a/TndFOIINicyuS6rk+ReQBzpkzxOiYLLxHj6wD8cRiMOsFSNNSNVF3V39ShND7BPPx3V2MkxaSYmOlaegrcGBoIHUXTULb42VvWJ9+3wQnE7MShO4/6Rk3yRJ8CO0d0lpnjIVmzDz9wrb+9FKuZdB667Yl0TGEs5fA74bA0JMyN5NPJ2+op16piDus6HN1ZjU0LTxnLJ/XurJt1s7kT0ophlGpEt9o9E9eYV0KT+3aeLMk4P1pjaFIjiv/0CQJkihGymc5ygoqr5pOmSQEX3UsCLdQWFUDDmMsTw/CMdRdhVxr8JuTGKeG87Jv8kE1Qj/iIOIiFqbRhY2AnIu64Cc=; mp_838c65c3e2afe9d50264505a75298594_mixpanel=%7B%22distinct_id%22%3A%20%221929df3b28c2360-0c37132650000a-1f525636-157188-1929df3b28d3962%22%2C%22%24device_id%22%3A%20%22192a52461fc1587-0e086d0f223f17-1f525636-157188-192a52461fd35bf%22%2C%22%24initial_referrer%22%3A%20%22%24direct%22%2C%22%24initial_referring_domain%22%3A%20%22%24direct%22%2C%22%24user_id%22%3A%20%221929df3b28c2360-0c37132650000a-1f525636-157188-1929df3b28d3962%22%2C%22Platform%22%3A%20%22Extension%22%2C%22Is%20Logged%20In%22%3A%20false%2C%22Tier%22%3A%20%22Anonymous%22%7D; passport_csrf_token=7e79d9a0e54d9f919a839050d2d70cbe; passport_csrf_token_default=7e79d9a0e54d9f919a839050d2d70cbe; s_v_web_id=verify_m2hliqoo_0szV3uE9_03zY_4wmI_8BVx_jK2SJvW5csvj; store-country-code=tw; store-country-code-src=uid; multi_sids=7426417624936023048%3A2974484ac056b432832c3b230ac1a7f2; cmpl_token=AgQQAPOHF-RO0rcxAdwgMd0x_XImsQsZ_5A3YNRJ0A; passport_auth_status=c7a95c48aa438155791a7a09a379f770%2C8282e1668a607143c2fb6c977cd566f3; passport_auth_status_ss=c7a95c48aa438155791a7a09a379f770%2C8282e1668a607143c2fb6c977cd566f3; sid_guard=2974484ac056b432832c3b230ac1a7f2%7C1729430695%7C15552000%7CFri%2C+18-Apr-2025+13%3A24%3A55+GMT; uid_tt=7cd5d68a333350c69a1faefc82f9d532bf8e62d266c679d5d539910cdc3baf31; uid_tt_ss=7cd5d68a333350c69a1faefc82f9d532bf8e62d266c679d5d539910cdc3baf31; sid_tt=2974484ac056b432832c3b230ac1a7f2; sessionid=2974484ac056b432832c3b230ac1a7f2; sessionid_ss=2974484ac056b432832c3b230ac1a7f2; sid_ucp_v1=1.0.0-KDIzZDNiNzRmODIyYWRlMzkwZjRmYmMxMjRkYzU4NGQ5YWVjYTUxYzYKIQiIiIzAuu_7h2cQp4nUuAYYswsgDDCX37-4BjgIQBJIBBADGgZtYWxpdmEiIDI5NzQ0ODRhYzA1NmI0MzI4MzJjM2IyMzBhYzFhN2Yy; ssid_ucp_v1=1.0.0-KDIzZDNiNzRmODIyYWRlMzkwZjRmYmMxMjRkYzU4NGQ5YWVjYTUxYzYKIQiIiIzAuu_7h2cQp4nUuAYYswsgDDCX37-4BjgIQBJIBBADGgZtYWxpdmEiIDI5NzQ0ODRhYzA1NmI0MzI4MzJjM2IyMzBhYzFhN2Yy; store-idc=maliva; tt-target-idc=alisg; tt-target-idc-sign=p6Ab3Z2p6E_EzAO24Jfw0KqBQrG-dtpvdzCZ7yW40kEyKZLGBXzHiuEw_dYr1LjiL1GyivoTMDtt1pJc56kvxAdId8gQJoUN82sY7uCI4gjDxJvKJQebeByokD3-i2D1SCao3nuq2Il6TdW1-kf13DXU78TMLsv2y2S6kHB1AYEsle74ose0cZ3lSe2t6xo1NHoBYSi9h39Un-P1fWK7i-XOausen5vvJnOUPj7bYrmMIA8sOpI-TT-ZfoceyAYSwCOw0M_8IcGoyS2bFt47S0OzBRpspANt9sJP_eDI0YfJV7-Q72Rg5mh5fPKyV1ifq98S5mAvlLL9NDpWnZ3YOap6xmp6eIGGVBYROQHB3VkFQRaUrk3J5PaoxDjFvHHq-GI8XL9zgTLp6jfwKEOxbxFvp-b334nEpuR-vLgIuJ9BBqn7VBlcsNAWUhNIiIzLYqD2l2MEBJhxKgd5XqPLT0ad2IvVJpYOicK7dYa6UhVGqvH4baSEz4DE4dj1VSLh; bm_sv=F653B9306700D7B1314A383537ADAB88~YAAQ5GLBF5LWNJySAQAANT0aqhnHn5Bg5zTyVXNGFK75KGC4G3WhUawsg/JibwxeJv9get1jWC/1Vm2H3HqSzR2Uj0mSHhIWxia/cu0zm11bOsxXteTuJ6T4LIbFLPunextzxtfuqi6KfL0pfqDPm8BUvGWbPAQ5TSf/OkvvRn4oSPO2ZP7gyEgV3irpdIDytemc14t+kg6QF4FMeW5AcejL933DajXs7Nn+YKVsHlWLDd6QC7PRRfw5vMe72D8J~1; ttwid=1%7CKVpbRGNETldWPzZvY734Xpy9-POuEscV1jqCKRfA3MM%7C1729430701%7Cb7b8a01e6c418a4b65785b3394c7e145251be93f410cde93e4a62f8b6522f2c8; odin_tt=a00c1d757485fbfeb7798c2f2f34e88b8935c5b6c5c2578ac4b310677ef4741c7f5d211a9739d65b80b530131003f77afeffb23494c3eeb1e6c4dd3cc14596d299c8240b4cba1e62550c52cb46b36b47; msToken=ZTT95D-Q2iwVRywxC5e262yW61dKkeM6bEj9aXd9RoGyUOsi9QLxd382_Rtal8oqYvpPlyaYBbTDzEUo9odcxTegRzR8t_P092hIL5lU5trBBzEwRKO6rtUvQ8-RrCoe1zYF0eASiVCzX21VqL-tpu7Keeg="
//)
//
//type Params struct {
//	keyword   string
//	count     int32
//	offset    int32
//	search_id string `json:"search_id"`
//	cookie    string
//}

//type AppParams struct {
//	keyword      string
//	offset       int32
//	count        int32
//	soft_type    int32
//	publish_time int32
//}
//
//var nomore int

// 加载配置
func loadConf() {
	err := conf.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}

func main() {
	// 加载配置
	loadConf()
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

//func main() {
//
//	loadConf()
//
//	var crawlCount int
//	var wt sync.WaitGroup
//	var mx sync.Mutex
//
//	nomore = 0
//
//	// app query
//	params := &AppParams{
//		keyword:      url2.PathEscape("china travel"),
//		count:        20,
//		offset:       0,
//		soft_type:    0,
//		publish_time: 0,
//	}
//
//	crawlCount = 0
//	for {
//		// get data
//		resp, err := getData(params)
//		if err != nil {
//			fmt.Println("getData出现错误，已退出循环，错误信息如下：")
//			fmt.Println(err)
//			break
//		}
//		// output csv
//		wt.Add(1)
//		go func() {
//			mx.Lock()
//			// todo fix bug: 1. 出现多goroutine并发写文件的问题
//			err := writeCSV(resp.ItemList)
//			if err != nil {
//				fmt.Println("写入csv时错误：", err)
//			}
//			defer func() {
//				mx.Unlock()
//				wt.Done()
//			}()
//		}()
//		crawlCount += len(resp.ItemList)
//		// fix bug: 2. 多次爬取的数据都是一样的
//		fmt.Println("正在修改参数")
//		params.offset = resp.Cursor
//		params.search_id = resp.Logid
//		fmt.Println("当前已爬取数据：", crawlCount, "条")
//		time.Sleep(1 * time.Second)
//	}
//	wt.Wait()
//
//	fmt.Println("nomore count =", nomore)
//}
//
//func getData(p *Params) (*model.VideoResp, error) {
//	url := "https://api.tikhub.io/api/v1/tiktok/web/fetch_search_video?keyword=" +
//		p.keyword + "&count=" + strconv.Itoa(int(p.count)) +
//		"&offset=" + strconv.Itoa(int(p.offset)) +
//		"&search_id=" + p.search_id + "&cookie=" + p.cookie
//	// user client
//	client := &http.Client{
//		Timeout: time.Second * 20,
//	}
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		fmt.Println("创建请求时错误：", err)
//		return nil, err
//	}
//	req.Header.Add("Authorization", "Bearer "+token)
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println("发送请求时错误：", err)
//		return nil, err
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println("关闭io时错误：", err)
//		}
//	}(res.Body)
//
//	// resp:read
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println("读取响应时错误：", err)
//		return nil, err
//	}
//
//	// resp:unmarshal
//	var resp model.VideoResp
//	err = json.Unmarshal(body, &resp)
//	if err != nil {
//		fmt.Println("反序列化时错误：", err)
//		fmt.Println("响应体：", string(body))
//		return nil, err
//	}
//	fmt.Println("抓取数据成功，一共", len(resp.ItemList), "条，正在检查...")
//
//	if resp.Code != 200 {
//		//  保存当前爬虫参数
//		fmt.Println("保存爬虫参数：", p)
//		fmt.Println("保存最新结果：", string(body))
//		return nil, errors.New("当前状态码不是200:" + string(body))
//	}
//	if resp.HasMore != 1 {
//		//  保存当前爬虫参数
//		fmt.Println("保存爬虫参数：", p)
//		fmt.Println("保存最新结果：", string(body))
//
//		return &resp, nil
//	}
//	if len(resp.ItemList) == 0 {
//		//  保存当前爬虫参数
//		fmt.Println("保存爬虫参数：", p)
//		fmt.Println("保存最新结果：", string(body))
//		return nil, errors.New("数据为0，响应体：" + string(body))
//	}
//
//	return &resp, nil
//}
//
//// write a batch of data to csv
//func writeCSV(itemList []model.Item) error {
//	tabHeader := []string{
//		"id",
//		"desc",
//		"cover",
//		"video_addr", // 自己合成
//		"author_id",
//		"author_name",
//		"digg_count",
//		"share_count",
//		"comment_count",
//		"play_count",
//		"collect_count",
//		"tag_list",
//		"is_verify",
//		"author_signature",
//		"author_follower",
//		"author_heart_count",
//		"author_video_count",
//		"author_digg_count",
//		"create_time",
//		"duration",
//		"music_title",
//		"music_url",
//		"music_author",
//	}
//	file, err := os.OpenFile("output_tiktok_1.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		return err
//	}
//	defer func(file *os.File) {
//		err := file.Close()
//		if err != nil {
//			fmt.Println(err)
//		}
//	}(file)
//
//	// init csv writer
//	writer := csv.NewWriter(file)
//	defer writer.Flush()
//
//	// write header
//	info, err := file.Stat()
//	if err != nil {
//		return err
//	}
//	if info.Size() == 0 {
//		err = writer.Write(tabHeader)
//		if err != nil {
//			return err
//		}
//	}
//
//	// write to csv
//	for _, item := range itemList {
//		_ = writer.Write([]string{
//			item.ID,
//			item.Desc,
//			item.Video.Cover,
//			formatVideoURL(item.Author.Nickname, item.Video.ID),
//			item.Author.ID,
//			item.Author.Nickname,
//			strconv.Itoa(item.Stats.DiggCount),
//			strconv.Itoa(item.Stats.ShareCount),
//			strconv.Itoa(item.Stats.CommentCount),
//			strconv.Itoa(item.Stats.PlayCount),
//			strconv.Itoa(item.Stats.CollectCount),
//			getTagList(item),
//			strconv.FormatBool(item.Author.Verified),
//			item.Author.Signature,
//			strconv.Itoa(item.AuthorStats.FollowerCount),
//			strconv.Itoa(item.AuthorStats.HeartCount),
//			strconv.Itoa(item.AuthorStats.VideoCount),
//			strconv.Itoa(item.AuthorStats.DiggCount),
//			time.Unix(item.CreateTime, 0).Format("2006-01-02 15:04:05"),
//			strconv.Itoa(item.Video.Duration),
//			item.Music.Title,
//			item.Music.PlayUrl,
//			item.Music.AuthorName,
//		})
//	}
//	fmt.Println("爬虫成功，数据已经写入csv文件")
//	return nil
//}
//func formatVideoURL(authorName string, videoId string) string {
//	return "https://www.tiktok.com/@" + authorName + "/video/" + videoId
//}
//func getTagList(item model.Item) string {
//	a := ""
//	for _, extra := range item.TextExtra {
//		a += extra.HashtagName
//		a += " "
//	}
//	return a
//}
