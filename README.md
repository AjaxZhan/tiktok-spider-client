# TikTok-Spider-Client

## Introduction

TikTok-Spider-Client is a web-spider program used to get TikTok video information, only for educational usage.

This program is basically used to finish academic social analysis.

This program is not directly sending http request to TikTok, instead, it sends http request to TikHub,a api server provider.

This repository is currently only for my own use, with relatively complete comments,
can provide basic shortcut code templates for friends who are too lazy to write their own http code.


## Quick start

1. Edit `config.json`, especially `token`.
2. In `main.go`, Load the configuration file firstly.
3. Create the client as we did in `main.go`, then call the `xxxAndStore` method.

For example:
```go
package main

import (
	"log"
	"tiktok-spider/client"
	"tiktok-spider/conf"
	"tiktok-spider/model"
)

func main() {
	loadConf()
	params := model.TiktokWebParamsSend{
		Keyword: "China travel",
		Offset:  0,
		Count:   1,
	}
	tiktokWebClient := client.NewTiktokWebClient(params)
	tiktokWebClient.SearchVideoAndStore()
}

// 加载配置
func loadConf() {
	err := conf.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}
```

Run the program:
```bash
go run main.go
```

## Disclaimer

This project is developed solely for educational purposes and is not intended for any illegal or unethical use.
It demonstrates the use of third-party APIs and does not aim to bypass any legal restrictions or terms of service set by the API providers.

The author holds no responsibility for how others may use this code.

Please ensure that your use of this project complies with all relevant laws and the terms of service of any APIs you utilize.