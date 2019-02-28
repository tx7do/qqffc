package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"../lottery"
)

// pathExists 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}

// WriteToFile 写入文件
func WriteToFile(strLotteryName string, lotteryData lottery.Data) {
	strFilePath := fmt.Sprintf("./%s/", strLotteryName)

	exist, err := pathExists(strFilePath)
	if err != nil {
		fmt.Printf("get dir error![%s]\n", err.Error())
		return
	}

	if !exist {
		err := os.Mkdir(strFilePath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%s]\n", err.Error())
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
	} else {
		strFluctuating = "+"
	}
	strFluctuating += strconv.FormatInt(int64(lotteryData.Fluctuating), 10)

	strText := fmt.Sprintf("%s\t%d,%d,%d,%d,%d\t%d\t%s\t%4d-%02d-%02d %02d:%02d:%02d\n",
		lotteryData.Issue,
		lotteryData.OpenNumbers[0], lotteryData.OpenNumbers[1], lotteryData.OpenNumbers[2], lotteryData.OpenNumbers[3], lotteryData.OpenNumbers[4],
		lotteryData.OnlineCount, strFluctuating,
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	if _, err = file.WriteString(strText); err != nil {
		log.Println(err)
	}
}
