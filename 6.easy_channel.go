package main

import (
	"fmt"
	"strings"
)

// channel
// 有開 channel 就一定要關 channel
// goRoutine 之間的溝通
// channel 可以被 buffer 也可以被 unbuffer
func shout(input <-chan string, output chan<- string) {
	for {
		s, found := <-input
		if !found {
			// 錯誤處理
			fmt.Print("inputChan 裡面沒有東西可以拉出來")
		}
		output <- fmt.Sprintf("%v !!!", strings.ToUpper(s))
	}
}

func shoutGame() {
	inputChan := make(chan string)
	outputChan := make(chan string)

	go shout(inputChan, outputChan) // 如果把這行註解掉，那輸入值會被塞入 inputChan 不過並沒有跟程式說要從 inputChan 拿值出來，所以就會進入 deadlock 然後就報錯
	for {
		fmt.Print("input:")
		userInput := ""
		_, _ = fmt.Scanln(&userInput)

		if userInput == "q" {
			break
		}
		inputChan <- userInput

		fmt.Printf("response:%v\n", <-outputChan)

	}
	fmt.Println("end")

	close(inputChan)
	close(outputChan)
}
