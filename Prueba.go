package main

import (
	"fmt"
)

func main() {
	numeros := []int{10, 20, 30, 40, 50}
	for ind, val := range numeros {
		fmt.Printf("[%d] = %d\n", ind, val)
	
	}

}
