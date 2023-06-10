package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan *Visitor, 10000)

	factory := NewVisitorFactory(ch)

	homeVisitor := factory.NewVisitor(Home, "/")

	ch <- homeVisitor

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		collector := NewCollector(ch)
		go func() {
			defer wg.Done()
			collector.Collect()
		}()
	}

	factory.Wait()
	fmt.Println("after wait.......................###############")

	close(ch)

	wg.Wait()
}
