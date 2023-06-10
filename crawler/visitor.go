package main

import (
	// "fmt"

	"fmt"

	"github.com/gocolly/colly/v2"
)

type VisitorDoner interface {
	Done(string)
}

type htmlCallBack struct {
	Selector string
	Callback colly.HTMLCallback
}

type Visitor struct {
	htmlCallbacks []*htmlCallBack
	url           string
	ch            chan<- *Visitor
	doner         []VisitorDoner
	factory       *VisitorFactory
}

func (visitor *Visitor) OnFailed() {
	fmt.Println("on failed : ", visitor.url)
	visitor.ch <- visitor
}

func (visitor *Visitor) RegisterDoner(doner VisitorDoner) {
	visitor.doner = append(visitor.doner, doner)
}

func (visitor *Visitor) OnSucceed() {
	for _, doner := range visitor.doner {
		doner.Done(visitor.url)
	}
}

func NewHomepageVisitor(visitorCollector chan<- *Visitor) *Visitor {
	visitor := &Visitor{ch: visitorCollector, url: Host}
	visitor.htmlCallbacks = []*htmlCallBack{
		{
			// 获取文章详情
			Selector: "article",
			Callback: func(h *colly.HTMLElement) {
				if detailUrl := h.ChildAttr("header[class='entry-header'] > h2[class='entry-title'] > a", "href"); len(detailUrl) > 0 {
					visitorCollector <- visitor.factory.NewVisitor(Article, detailUrl)
				}
			},
		}, {
			// 分类
			Selector: `aside[class='widget widget_categories'] > ul`,
			Callback: func(h *colly.HTMLElement) {
				h.ForEach("li", func(i int, h *colly.HTMLElement) {

					if link := h.ChildAttr("a", "href"); len(link) > 0 {
						//fmt.Println("categories : ", link)
						visitorCollector <- visitor.factory.NewVisitor(Category, link)
					}

					h.ForEach("ul > li", func(i int, h *colly.HTMLElement) {
						if link := h.ChildAttr("a", "href"); len(link) > 0 {
							//fmt.Println("subCategories : ", link)
							visitorCollector <- visitor.factory.NewVisitor(Category, link)
						}
					})
				})
			},
		}, {
			// 标签
			Selector: `aside[class='widget widget_tag_cloud'] > div`,
			Callback: func(h *colly.HTMLElement) {
				h.ForEach("a", func(i int, h *colly.HTMLElement) {
					if link := h.Attr("href"); len(link) > 0 {
						//fmt.Println("tag : ", h.Text, link)
						visitorCollector <- visitor.factory.NewVisitor(Tag, link)
					}
				})
			},
		}, {
			// 归档
			Selector: `aside[class='widget widget_archive'] > ul`,
			Callback: func(h *colly.HTMLElement) {
				h.ForEach("li", func(i int, h *colly.HTMLElement) {
					if link := h.ChildAttr("li > a", "href"); len(link) > 0 {
						//fmt.Println("archive : ", h.Text, link)
						visitorCollector <- visitor.factory.NewVisitor(Archive, link)
					}

				})
			},
		}, {
			// 导航栏
			Selector: `ul[class='nav navbar-nav']`,
			Callback: func(h *colly.HTMLElement) {
				h.ForEach("li", func(i int, h *colly.HTMLElement) {
					if link := h.ChildAttr("a", "href"); len(link) > 0 {
						//fmt.Println("nav : ", h.Text, link, i)
						if i > 0 {
							visitorCollector <- visitor.factory.NewVisitor(Nav, link)
						}
					}
				})
			},
		}, {
			// next page
			Selector: `div[class='wp-pagenavi']`,
			Callback: func(h *colly.HTMLElement) {
				if link := h.ChildAttr(`a[class='nextpostslink']`, "href"); len(link) > 0 {
					//fmt.Println("pagination : ", link)
					visitorCollector <- visitor.factory.NewVisitor(Pagination, link)
				}
			},
		},
	}
	return visitor
}

func NewArticleVisitor(visitorCollector chan<- *Visitor, url string) *Visitor {
	return &Visitor{url: url}
}

func NewArchiveVisitor(visitorCollector chan<- *Visitor, url string) *Visitor {
	visitor := &Visitor{url: url}
	visitor.htmlCallbacks = []*htmlCallBack{
		{
			// next page
			Selector: `div[class='wp-pagenavi']`,
			Callback: func(h *colly.HTMLElement) {
				if link := h.ChildAttr(`a[class='nextpostslink']`, "href"); len(link) > 0 {
					//fmt.Println("pagination : ", link)
					visitorCollector <- visitor.factory.NewVisitor(Archive, link)
				}
			},
		},
	}

	return visitor
}

func NewTagVisitor(visitorCollector chan<- *Visitor, url string) *Visitor {
	visitor := &Visitor{url: url}
	visitor.htmlCallbacks = []*htmlCallBack{
		{
			// next page
			Selector: `div[class='wp-pagenavi']`,
			Callback: func(h *colly.HTMLElement) {
				if link := h.ChildAttr(`a[class='nextpostslink']`, "href"); len(link) > 0 {
					//fmt.Println("pagination : ", link)
					visitorCollector <- visitor.factory.NewVisitor(Tag, link)
				}
			},
		},
	}

	return visitor
}

func NewCategoryVisitor(visitorCollector chan<- *Visitor, url string) *Visitor {
	visitor := &Visitor{ch: visitorCollector, url: url}
	visitor.htmlCallbacks = []*htmlCallBack{
		{
			// next page
			Selector: `div[class='wp-pagenavi']`,
			Callback: func(h *colly.HTMLElement) {
				if link := h.ChildAttr(`a[class='nextpostslink']`, "href"); len(link) > 0 {
					//fmt.Println("pagination : ", link)
					visitorCollector <- visitor.factory.NewVisitor(Category, link)
				}
			},
		},
	}

	return visitor
}

func NewNavVisitor(visitorCollector chan<- *Visitor, url string) *Visitor {
	return &Visitor{url: url}
}

func NewPagination(visitorCollector chan<- *Visitor, url string) *Visitor {
	visitor := &Visitor{ch: visitorCollector, url: url}
	visitor.htmlCallbacks = []*htmlCallBack{
		{
			// next page
			Selector: `div[class='wp-pagenavi']`,
			Callback: func(h *colly.HTMLElement) {
				if link := h.ChildAttr(`a[class='nextpostslink']`, "href"); len(link) > 0 {
					//fmt.Println("pagination : ", link)
					visitorCollector <- visitor.factory.NewVisitor(Pagination, link)
				}
			},
		}, {
			// 获取文章详情
			Selector: "article",
			Callback: func(h *colly.HTMLElement) {
				if detailUrl := h.ChildAttr("header[class='entry-header'] > h2[class='entry-title'] > a", "href"); len(detailUrl) > 0 {
					visitorCollector <- visitor.factory.NewVisitor(Article, detailUrl)
				}
			},
		},
	}

	return visitor
}
