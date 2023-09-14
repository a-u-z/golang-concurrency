package main

import (
	"fmt"
	"log"
	"sync"
)

type P1 struct {
	One int
}

type P2 struct {
	Two int
}

func Practice() {
	var wg sync.WaitGroup
	p1Channel := make(chan P1) // 初始化 p1Channel
	p2Channel := make(chan P2) // 初始化 p2Channel
	errorChannel := make(chan error)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		wg.Add(1)
		go process1(p1Channel, i, errorChannel) // 這邊要用 go 去執行，不然會死鎖
		go process2(p2Channel, i, errorChannel)
	}
	go func() {
		for p1 := range p1Channel {
			log.Printf("here is p1:%+v", p1)
			wg.Done()
		}
	}()
	go func() {
		for p2 := range p2Channel {
			log.Printf("here is p2:%+v", p2)
			wg.Done()
		}
	}()
	errString := []string{}
	go func() {
		for err := range errorChannel {
			errString = append(errString, err.Error())
			wg.Done()
		}
	}()

	wg.Wait()
	defer close(p1Channel)
	defer close(p2Channel)
	defer close(errorChannel)

	log.Printf("here is errString:%+v", errString)
}

func process1(p1Channel chan P1, i int, errorChannel chan error) {
	if i == 5 {
		errorChannel <- fmt.Errorf("this is error test, now is int %v", i)
		return
	}
	pOne := P1{
		One: i,
	}
	p1Channel <- pOne
}

func process2(p2Channel chan P2, i int, errorChannel chan error) {
	if i == 6 {
		errorChannel <- fmt.Errorf("this is error test, now is int %v", i)
		return
	}

	pTwo := P2{
		Two: i,
	}
	p2Channel <- pTwo
}
