package main

import (
	"fmt"
	"time"
	"math/rand"
)

func main() {
	for {
		time.Sleep(1*time.Millisecond)
		fmt.Printf("%dusec \n", rand.Int63n(9999999))
	}
}
