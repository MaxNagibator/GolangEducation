package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	taskCount := 20

	var wg sync.WaitGroup
	wg.Add(taskCount) // в группе две горутины

	intCh := make(chan int, 5)
	defer close(intCh)

	var mutex sync.Mutex

	fmt.Println("hello parallel")
	t := time.Now()

	for i := 1; i <= taskCount; i++ {
		go waitGoGo(t, i, intCh, &mutex, &wg)
	}

	wg.Wait()

	fmt.Println("The End")
}

func square(ch chan int) {

	fmt.Println("num := ", "huy")
	num := <-ch
	fmt.Println("num := ", num)
	ch <- num * num
}

func waitGoGo(t time.Time, i int, ch chan int, mutex *sync.Mutex, waitGroup *sync.WaitGroup) {

	if i%4 == 0 {
		mutex.Lock()
		defer mutex.Unlock()
	}
	defer waitGroup.Done()
	ch <- i
	seconds := time.Duration(1) * time.Second
	time.Sleep(seconds)
	fmt.Println(time.Since(t), "COMPLETE", i)
	<-ch
}
