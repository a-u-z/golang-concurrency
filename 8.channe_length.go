package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

func listenToChannel(c <-chan int, wg *sync.WaitGroup) {
	for {
		color.Blue(fmt.Sprintf("%v out \n", <-c))
		wg.Done()
		time.Sleep(time.Second)
	}
}

func channelLength() {
	loops := 15
	wg := sync.WaitGroup{}
	wg.Add(loops)
	// intChannel := make(chan int)
	intChannel := make(chan int, 10) // 一開始直接給 channel 10 個空間，
	go listenToChannel(intChannel, &wg)
	// 會一直塞東西進去，直到塞了 10 個，到達上限，要等 listenToChannel 拉出來後，才能再塞 int 進去
	for i := 1; i <= 15; i++ {
		intChannel <- i
		color.Yellow(fmt.Sprintf("%v in \n", i))
	}
	wg.Wait()
	close(intChannel)
}
