package main

import (
	"log"
	"sync"
)

func waitGroup() {

	var wg sync.WaitGroup
	words := []string{
		"alpha",
		"beta",
		"delta",
		"gamma",
		"zeta",
	}
	wg.Add(len(words)) // 設定要等待的數量

	for i, v := range words {
		go printSomething(i, v, &wg)
	}
	wg.Wait()
}

func printSomething(index int, word string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("%v:%+v", index, word)
}
