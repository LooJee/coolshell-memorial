package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Collector struct {
	collector *colly.Collector
	ch        chan *Visitor
}

func NewCollector(ch chan *Visitor) *Collector {
	return &Collector{
		collector: colly.NewCollector(
			colly.AllowedDomains(Host, "coolshell.cn"),
		),
		ch: ch,
	}
}

func (collector *Collector) Collect() {
	collector.collector.SetRequestTimeout(time.Minute)

	for visitor := range collector.ch {
		var (
			cloner    = collector.collector.Clone()
			onSucceed bool
		)

		for _, cb := range visitor.htmlCallbacks {
			cloner.OnHTML(cb.Selector, cb.Callback)
		}

		cloner.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		cloner.OnError(func(r *colly.Response, err error) {
			fmt.Println("something error...", err)
			go visitor.OnFailed()
		})

		cloner.OnResponse(func(r *colly.Response) {
			// fmt.Println("response : ", r.Request.URL, r.Request.URL.RequestURI())

			if r.StatusCode == 200 {
				file := "./dist" + r.Request.URL.RequestURI()
				if !strings.HasSuffix(file, ".html") {
					file += "/index.html"
				}
				os.MkdirAll(filepath.Dir(file), 0777)

				if err := r.Save(file); err != nil {
					panic(err)
				}
				onSucceed = true
			} else {
				fmt.Println("+++++++++++++++++++++++failed at response : ", r.Request.URL.RequestURI(), r.StatusCode)
				panic("failed at response")
			}
		})

		cloner.Visit(visitor.url)
		if onSucceed {
			visitor.OnSucceed()
		}
		fmt.Println("collector channel size : ", len(collector.ch))
	}
}
