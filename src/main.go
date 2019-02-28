package main

import (
	"github.com/jakecoffman/cron"

	"./crawler/qqffc"
	"./crawler/txffc"
)

func main() {
	spec := "10 * * * * ?"
	cronJob := cron.New()
	cronJob.AddFunc(spec, func() {
		qqffc.StartCrawler()
		txffc.StartCrawler()
	}, "crawl")
	cronJob.Start()

	defer cronJob.Stop()

	select {}
}
