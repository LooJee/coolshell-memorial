package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type visitorType int

const (
	Home visitorType = iota
	Archive
	Tag
	Nav
	Article
	Category
	Pagination
)

type VisitorFactory struct {
	wg          sync.WaitGroup
	collectorCh chan *Visitor
	cnt         atomic.Int32
}

func NewVisitorFactory(collectorCh chan *Visitor) *VisitorFactory {
	return &VisitorFactory{
		wg:          sync.WaitGroup{},
		collectorCh: collectorCh,
	}
}

func (factory *VisitorFactory) NewVisitor(vt visitorType, url string) *Visitor {
	factory.wg.Add(1)
	factory.cnt.Add(1)

	var visitor *Visitor

	switch vt {
	case Home:
		visitor = NewHomepageVisitor(factory.collectorCh)
	case Archive:
		visitor = NewArchiveVisitor(factory.collectorCh, url)
	case Tag:
		visitor = NewTagVisitor(factory.collectorCh, url)
	case Nav:
		visitor = NewNavVisitor(factory.collectorCh, url)
	case Article:
		visitor = NewArticleVisitor(factory.collectorCh, url)
	case Category:
		visitor = NewCategoryVisitor(factory.collectorCh, url)
	case Pagination:
		visitor = NewPagination(factory.collectorCh, url)
	}

	visitor.RegisterDoner(factory)
	visitor.factory = factory
	visitor.vt = vt

	fmt.Println("new visitor : ", vt, url)

	return visitor
}

func (factory *VisitorFactory) Done(url string) {
	factory.cnt.Add(-1)
	fmt.Println(url, "done........................, left : ", factory.cnt.Load())
	factory.wg.Done()
}

func (factory *VisitorFactory) Wait() {
	factory.wg.Wait()
}
