package main

import (
	"fmt"
	"math/rand"
)

func main() {

	array := rand.Perm(100)

	fmt.Println("array:", array)
	taskCount := 20

	intCh := make(chan int)
	defer close(intCh)

	fmt.Println("hello parallel")

	for i := 0; i < taskCount; i++ {
		subArray := array[i*5 : (i+1)*5]
		go arraySum(subArray, intCh)
	}

	sum := 0
	for i := 0; i < taskCount; i++ {
		sum += <-intCh
	}

	sum2 := 0
	for i := 0; i < len(array); i++ {
		sum2 += array[i]
	}

	fmt.Println("Sum parallel:", sum)
	fmt.Println("Sum not parallel:", sum2)
}

func arraySum(array []int, ch chan int) {

	fmt.Println("array:", array)
	result := 0
	for i := 0; i < len(array); i++ {
		result += array[i]
	}

	ch <- result
}
