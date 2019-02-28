package lottery

import (
	"fmt"
	"strconv"
	"time"
)

// Data 分分彩数据
type Data struct {
	Issue       string // 期号
	OpenNumbers []int  // 开奖号码
	OnlineCount int    // 在线人数
	Fluctuating int    // 波动值
}

// ComputeLotteryNumbers 计算开奖号码
func ComputeLotteryNumbers(onlineCount int) []int {

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

// ComputeLotteryIssue 计算彩票期数
func ComputeLotteryIssue() string {
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

// Print 打印彩票信息
func Print(strLotteryName string, lotteryData Data) {
	t := time.Now()
	strFluctuating := ""
	if lotteryData.Fluctuating <= 0 {
	} else {
		strFluctuating = "+"
	}
	strFluctuating += strconv.FormatInt(int64(lotteryData.Fluctuating), 10)

	fmt.Printf("%s 期号:%s 开奖号码:%d,%d,%d,%d,%d 腾讯在线人数:%d 波动值:%s 时间:%4d-%02d-%02d %02d:%02d:%02d\n",
		strLotteryName, lotteryData.Issue,
		lotteryData.OpenNumbers[0], lotteryData.OpenNumbers[1], lotteryData.OpenNumbers[2], lotteryData.OpenNumbers[3], lotteryData.OpenNumbers[4],
		lotteryData.OnlineCount, strFluctuating,
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
