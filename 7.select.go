package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func inputSomethingToChannel(s string, c chan<- string) {
	for {
		switch s {
		case "channel 1":
			time.Sleep(3 * time.Second)
		case "channel 2":
			time.Sleep(1 * time.Second)
		default:
			time.Sleep(200 * time.Millisecond)
		}
		c <- fmt.Sprintf("here is %v", s)
	}
}

func selectPractice() {
	channel1 := make(chan string)
	channel2 := make(chan string)

	go inputSomethingToChannel("channel 1", channel1)
	go inputSomethingToChannel("channel 2", channel2)
	for {
		// channel 版本的 switch case
		// 可以把 channel 內拿出的東西，進行判斷，分配到不同的邏輯
		// 一定要有一個 default 來避免死鎖，因為如果 channel 1 跟 channel 2 都沒有東西拉出來的話，這個 select 就爆了
		select {
		case s := <-channel1:
			color.Cyan(fmt.Sprintf("select 1: %v\n", s))
		case s := <-channel1:
			color.Yellow(fmt.Sprintf("select 2: %v\n", s))
		case s := <-channel2:
			color.Magenta(fmt.Sprintf("select 3: %v\n", s))
		case s := <-channel2:
			color.Green(fmt.Sprintf("select 4: %v\n", s))

		default:
			time.Sleep(2 * time.Second)
			color.Red("here is default, nothing from channel 1 and channel 2\n")
		}
	}
}
