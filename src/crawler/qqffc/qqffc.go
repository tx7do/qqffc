package qqffc

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"

	"../../db"
	"../../lottery"
)

// come from https://hk.saowen.com/a/9de9acc6f0f77dd0636e88b2d721ae18e78a8799fbf8d862cb401f95ab61055a
const (
	URL1 = "https://cgi.im.qq.com/cgi-bin/minute_city"
	URL2 = "https://cgi.im.qq.com/data/1min_city.dat"
)

// CityOnlineData 城市在线人数汇总数据,使用URL1或者URL2所采集的数据.
type CityOnlineData struct {
	Time   string `json:"time"`
	Minute []int  `json:"minute"`
}

var qqffc lottery.Data
var qqffcPre lottery.Data
var lotteryData lottery.Data

// StartCrawler 启动采集
func StartCrawler() {

	qqffcPre = lottery.Data{Issue: "", OpenNumbers: []int{0, 0, 0, 0, 0}, OnlineCount: 0, Fluctuating: 0}

	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("Visiting", r.Request.URL)

		crawlData := &CityOnlineData{}
		err := json.Unmarshal(r.Body, crawlData)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		//fmt.Println(crawlData)
		processCrawlData(crawlData)
	})

	c.Limit(&colly.LimitRule{
		Parallelism: 1,
		RandomDelay: 5 * time.Second,
	})

	err := c.Visit(URL2)
	if err != nil {
		log.Fatal("Visit Error: ", err)
	}

	c.Wait()
}

// processCrawlData 执行采集
// result like https://1680118.com/view/tencent_ffc/ssc_index.html
func processCrawlData(crawlData *CityOnlineData) {
	// 采集数据
	lotteryData.OnlineCount = crawlData.Minute[0]

	// 计算开奖号码
	var arrOpenNumbers []int
	arrOpenNumbers = lottery.ComputeLotteryNumbers(lotteryData.OnlineCount)
	lotteryData.OpenNumbers = arrOpenNumbers

	// 计算开奖期号
	strIssue := lottery.ComputeLotteryIssue()
	lotteryData.Issue = strIssue

	lotteryData.Fluctuating = lotteryData.OnlineCount - qqffcPre.OnlineCount

	// 输出结果
	lottery.Print("QQ  分分彩", lotteryData)

	db.WriteToFile("QQ分分彩", lotteryData)
}
