package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	c1 := make(chan string, 10)
	c2 := make(chan string, 10)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go receive(c1, c2, &wg)

	for i := 0; i < 10; i++ {
		c1 <- fmt.Sprint(i)
	}

	time.Sleep(time.Second * 2)
	close(c1)
	close(c2)
	wg.Wait()
}

func receive(ch1, ch2 chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case s1, ok := <-ch1:
			if ok {
				fmt.Println(s1)
			} else {
				return
			}
		case s2 := <-ch2:
			fmt.Println(s2)
		default:
			fmt.Println("default")
			time.Sleep(time.Second * 1)
		}
	}
}
