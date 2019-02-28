package txffc

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly"

	"../../db"
	"../../lottery"
)

// come from https://hk.saowen.com/a/9de9acc6f0f77dd0636e88b2d721ae18e78a8799fbf8d862cb401f95ab61055a
const (
	URL3 = "https://mma.qq.com/cgi-bin/im/online"
)

// TotalOnlineData 总在线人数数据,使用URL3所采集的数据.
type TotalOnlineData struct {
	Current int `json:"c"`
	History int `json:"h"`
	EC      int `json:"ec"`
}

var txffc lottery.Data
var txffcPre lottery.Data
var lotteryData lottery.Data

// StartCrawler 启动采集
func StartCrawler() {
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("Visiting", r.Request.URL)

		htmlText := r.Body
		//fmt.Println("html response text:", htmlText)
		if len(htmlText) < 12 {
			return
		}

		strJSON := htmlText[12 : len(htmlText)-1]
		//fmt.Println(strJSON)

		crawlData := &TotalOnlineData{}

		// 解析json
		err := json.Unmarshal([]byte(strJSON), &crawlData)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		//fmt.Println(crawlData)
		processCrawlData(crawlData)
	})

	err := c.Visit(URL3)
	if err != nil {
		log.Fatal("Visit Error: ", err)
	}

}

// processCrawlData 执行采集
// result like : http://www.off0.com/fenfencai.php
func processCrawlData(crawlData *TotalOnlineData) {

	lotteryData.OnlineCount = crawlData.Current

	// 计算开奖号码
	var arrOpenNumbers []int
	arrOpenNumbers = lottery.ComputeLotteryNumbers(crawlData.Current)
	lotteryData.OpenNumbers = arrOpenNumbers

	/// 计算开奖期号
	strIssue := lottery.ComputeLotteryIssue()
	lotteryData.Issue = strIssue

	lotteryData.Fluctuating = lotteryData.OnlineCount - txffcPre.OnlineCount

	lottery.Print("腾讯分分彩", lotteryData)

	db.WriteToFile("腾讯分分彩", lotteryData)
}
