package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Collector struct {
	collector *colly.Collector
	ch        chan *Visitor
	number    int
}

func NewCollector(number int, ch chan *Visitor) *Collector {
	return &Collector{
		collector: colly.NewCollector(
			colly.AllowedDomains(Host, "coolshell.cn"),
		),
		ch:     ch,
		number: number,
	}
}

func (collector *Collector) Collect() {
	for visitor := range collector.ch {
		fmt.Println(collector.number, "get visitor :", visitor.url)
		var (
			cloner = collector.collector.Clone()
		)

		for _, cb := range visitor.htmlCallbacks {
			cloner.OnHTML(cb.Selector, cb.Callback)
		}

		cloner.OnRequest(func(r *colly.Request) {
			fmt.Println(collector.number, "Visiting", r.URL.String())
		})

		cloner.OnError(func(r *colly.Response, err error) {
			fmt.Println(collector.number, r.Request.URL.RequestURI(), "something error...", err)
			visitor.OnFailed()
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
			} else {
				fmt.Println("+++++++++++++++++++++++failed at response : ", r.Request.URL.RequestURI(), r.StatusCode)
				panic("failed at response")
			}
		})

		cloner.Visit(visitor.url)
		visitor.Done()
		fmt.Println("collector channel size : ", len(collector.ch))
	}
}
