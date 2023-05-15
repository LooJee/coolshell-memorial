package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("coolshell.cn", "www.coolshell.cn"),
	)

	c.SetRequestTimeout(60 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML(`span.pages`, func(h *colly.HTMLElement) {
		// h.Text should be 第 1 / 74 页
		texts := strings.Split(h.Text, " ")
		if len(texts) != 5 {
			panic("invalid texts")
		}

		totalPages, err := strconv.Atoi(texts[3])
		if err != nil {
			panic(fmt.Sprintf("parse total pages failed: %v", err))
		}

		fmt.Println("total pages : ", totalPages)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("something error...", err)
	})

	c.Visit("https://coolshell.cn/")
}
