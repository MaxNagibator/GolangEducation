package main

import "fmt"

func main() {
	defer finish("1")
	defer finish("2")
	defer finish("3")

	test := 1
	test1 := 0
	test2 := test / test1
	fmt.Println(test2)
	panic("trevoga trevoga volk ukral zaychat")
	var numbers [5]int = [5]int{7, 2, 3, 4, 5}

	for i := 1; i < len(numbers); i++ {
		for i2 := i; i2 < len(numbers); i2++ {
			if numbers[i2] < numbers[i2-1] {
				swap(&numbers, i2)
			}
		}
	}

	fmt.Println(numbers)

	for j := 0; j < len(numbers); j++ {
		fmt.Println(numbers[j])
	}
}

func swap(numbers *[5]int, i2 int) {
	var temp int = numbers[i2]
	numbers[i2] = numbers[i2-1]
	numbers[i2-1] = temp
}

func finish(bla string) {
	fmt.Println(bla)
}
