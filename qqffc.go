package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jakecoffman/cron"
)

/// come from https://hk.saowen.com/a/9de9acc6f0f77dd0636e88b2d721ae18e78a8799fbf8d862cb401f95ab61055a
const (
	URL1 = "https://cgi.im.qq.com/cgi-bin/minute_city"
	URL2 = "https://cgi.im.qq.com/data/1min_city.dat"

	URL3 = "https://mma.qq.com/cgi-bin/im/online"
)

/// 城市在线人数汇总数据,使用URL1或者URL2所采集的数据.
type CityOnlineData struct {
	Time   string `json:"time"`
	Minute []int  `json:"minute"`
}

/// 总在线人数数据,使用URL3所采集的数据.
type TotalOnlineData struct {
	Current int `json:"c"`
	History int `json:"h"`
	EC      int `json:"ec"`
}

/// 分分彩数据
type LotteryData struct {
	Issue       string ///< 期号
	OpenNumbers []int  ///< 开奖号码
	OnlineCount int    ///< 在线人数
	Fluctuating int    ///< 波动值
}

/// 计算开奖号码
func computeLotteryNumbers(onlineCount int) []int {

	var arrOpenNumbers []int

	// 分解在线人数到数组
	strOnlineCount := fmt.Sprintf("%d", onlineCount)
	var arrNumbers []int
	for i := 0; i < len(strOnlineCount); i++ {
		num, _ := strconv.Atoi(string(strOnlineCount[i]))
		arrNumbers = append(arrNumbers, num)
	}
	// fmt.Println(arrNumbers)

	if len(arrNumbers) < 4 {
		return arrOpenNumbers
	}

	// 计算总和
	nTotal := 0
	for i := 0; i < len(arrNumbers); i++ {
		nTotal += arrNumbers[i]
	}
	// fmt.Println("total:", nTotal)

	//取总和的个位数作为万位数
	arrOpenNumbers = append(arrOpenNumbers, nTotal%10)

	//在线人数最后的4位作为开奖号码的千百十个位
	for i := len(arrNumbers) - 4; i < len(arrNumbers); i++ {
		arrOpenNumbers = append(arrOpenNumbers, arrNumbers[i])
	}
	// fmt.Println("open numbers:", arrOpenNumbers)

	return arrOpenNumbers
}

/// 计算彩票期数
func computeLotteryIssue() string {
	currentTime := time.Now()

	nSequence := (currentTime.Hour() * 60) + currentTime.Minute()
	if nSequence == 0 {
		nSequence = 1440
		yesTime := currentTime.AddDate(0, 0, -1)
		currentTime = yesTime
	}
	strSequence := fmt.Sprintf("%.4d", nSequence)

	strIssue := fmt.Sprintf("%d%d%d-%s", currentTime.Year(), currentTime.Month(), currentTime.Day(), strSequence)

	return strIssue
}

/// 打印彩票信息
func printLotteryDetail(strLotteryName string, lotteryData LotteryData) {
	t := time.Now()
	strFluctuating := ""
	if lotteryData.Fluctuating <= 0 {
		strFluctuating = fmt.Sprintf("%d", lotteryData.Fluctuating)
	} else {
		strFluctuating = fmt.Sprintf("+%d", lotteryData.Fluctuating)
	}
	fmt.Printf("%s 期号:%s 开奖号码:%d,%d,%d,%d,%d 腾讯在线人数:%d 波动值:%s 时间:%4d-%02d-%02d %02d:%02d:%02d\n",
		strLotteryName, lotteryData.Issue,
		lotteryData.OpenNumbers[0], lotteryData.OpenNumbers[1], lotteryData.OpenNumbers[2], lotteryData.OpenNumbers[3], lotteryData.OpenNumbers[4],
		lotteryData.OnlineCount, strFluctuating,
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

///
func writeToFile(strLotteryName string, lotteryData LotteryData) {
	strFilePath := fmt.Sprintf("./%s/", strLotteryName)

	exist, err := pathExists(strFilePath)
	if err != nil {
		log.Println("get dir error![%v]", err)
		return
	}

	if !exist {
		err := os.Mkdir(strFilePath, os.ModePerm)
		if err != nil {
			log.Println("mkdir failed![%v]", err)
		}
	}

	t := time.Now()

	strFileName := fmt.Sprintf("%s%4d-%02d-%02d.txt", strFilePath, t.Year(), t.Month(), t.Day())
	file, err := os.OpenFile(strFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	strFluctuating := ""
	if lotteryData.Fluctuating <= 0 {
		strFluctuating = fmt.Sprintf("%d", lotteryData.Fluctuating)
	} else {
		strFluctuating = fmt.Sprintf("+%d", lotteryData.Fluctuating)
	}

	strText := fmt.Sprintf("%s\t%d,%d,%d,%d,%d\t%d\t%s\t%4d-%02d-%02d %02d:%02d:%02d\n",
		lotteryData.Issue,
		lotteryData.OpenNumbers[0], lotteryData.OpenNumbers[1], lotteryData.OpenNumbers[2], lotteryData.OpenNumbers[3], lotteryData.OpenNumbers[4],
		lotteryData.OnlineCount, strFluctuating,
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	if _, err = file.WriteString(strText); err != nil {
		log.Println(err)
	}
}

/// 采集数据
func fetchTotalOnlineData(url string) TotalOnlineData {
	var data TotalOnlineData

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println(err)
		return data
	}

	htmlText := doc.Text()
	//fmt.Println("html response text:", htmlText)
	if len(htmlText) < 12 {
		return data
	}

	strJSON := htmlText[12 : len(htmlText)-1]
	//fmt.Println(strJSON)

	// 解析json
	err = json.Unmarshal([]byte(strJSON), &data)
	if err != nil {
		log.Println(err)
		return data
	}
	//fmt.Println(data)

	return data
}

/// 执行采集
/// result like : http://www.off0.com/fenfencai.php
func crawlTotalOnlineData(preData *LotteryData) LotteryData {
	var lotteryData LotteryData

	// 采集数据
	var crawlerData TotalOnlineData
	crawlerData = fetchTotalOnlineData(URL3)
	lotteryData.OnlineCount = crawlerData.Current

	// 计算开奖号码
	var arrOpenNumbers []int
	arrOpenNumbers = computeLotteryNumbers(crawlerData.Current)
	lotteryData.OpenNumbers = arrOpenNumbers

	/// 计算开奖期号
	strIssue := computeLotteryIssue()
	lotteryData.Issue = strIssue

	if preData != nil {
		lotteryData.Fluctuating = lotteryData.OnlineCount - preData.OnlineCount
	}

	printLotteryDetail("腾讯分分彩", lotteryData)

	writeToFile("腾讯分分彩", lotteryData)

	return lotteryData
}

/// 采集数据
func fetchCityOnlineData(url string) CityOnlineData {
	var data CityOnlineData

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println(err)
		return data
	}

	htmlText := doc.Text()
	//fmt.Println("html response text:", htmlText)

	// 解析json
	err = json.Unmarshal([]byte(htmlText), &data)
	if err != nil {
		log.Println(err, htmlText)
		return data
	}
	//fmt.Println(data)

	return data
}

/// 执行采集
/// result like https://1680118.com/view/tencent_ffc/ssc_index.html
func crawlCityOnlineData(preData *LotteryData) LotteryData {
	var lotteryData LotteryData

	// 采集数据
	var crawlerData CityOnlineData
	crawlerData = fetchCityOnlineData(URL2)
	lotteryData.OnlineCount = crawlerData.Minute[0]

	// 计算开奖号码
	var arrOpenNumbers []int
	arrOpenNumbers = computeLotteryNumbers(lotteryData.OnlineCount)
	lotteryData.OpenNumbers = arrOpenNumbers

	// 计算开奖期号
	strIssue := computeLotteryIssue()
	lotteryData.Issue = strIssue

	if preData != nil {
		lotteryData.Fluctuating = lotteryData.OnlineCount - preData.OnlineCount
	}

	// 输出结果
	printLotteryDetail("QQ  分分彩", lotteryData)

	writeToFile("QQ分分彩", lotteryData)

	return lotteryData
}

func testComputeLotteryNumbers() {
	onlineCount := 123456789
	computeLotteryNumbers(onlineCount)
}

/// 采集一期
func crawlOne() {
	crawlTotalOnlineData(nil)
	crawlCityOnlineData(nil)
}

/// 采集多期
func crawlAll() {
	var qqffc LotteryData
	var txffc LotteryData

	qqffcPre := LotteryData{Issue: "", OpenNumbers: []int{0, 0, 0, 0, 0}, OnlineCount: 0, Fluctuating: 0}
	txffcPre := LotteryData{Issue: "", OpenNumbers: []int{0, 0, 0, 0, 0}, OnlineCount: 0, Fluctuating: 0}

	spec := "35 * * * * ?"
	cronJob := cron.New()
	cronJob.AddFunc(spec, func() {
		txffc = crawlTotalOnlineData(&txffcPre)
		txffcPre = txffc

		qqffc = crawlCityOnlineData(&qqffcPre)
		qqffcPre = qqffc
	}, "crawl")
	cronJob.Start()

	defer cronJob.Stop()

	select {}
}

func main() {
	//crawlOne()
	crawlAll()
}
