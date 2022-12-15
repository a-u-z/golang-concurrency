package main

import (
	"log"
	"sync"
)

type bb struct {
	msg string
}

// mutex 是防止同時改動到一個變數
// waitGroup 是要讓等全部有被登記數量(wg.Add)的事件完成後（wg.Done），才會結束程式運行
// 兩個都需要
func (a *bb) updateMsg(s string, m *sync.Mutex, w *sync.WaitGroup) {
	defer w.Done()
	m.Lock()
	a.msg = s

	m.Unlock()
}

// 這樣會有 race condition 產生， go run -race .，可以偵測出來
// 這個要到 main 才可以跑，用 -race 才確定沒有 race condition
func mutexPractice() {
	newBb := &bb{
		msg: "Good morning",
	}
	var mutex sync.Mutex
	var wg sync.WaitGroup

	wg.Add(12)
	go newBb.updateMsg("Good afternoon", &mutex, &wg)
	go newBb.updateMsg("Good night", &mutex, &wg)
	go newBb.updateMsg("Good afternoon", &mutex, &wg)
	go newBb.updateMsg("Good night", &mutex, &wg)
	go newBb.updateMsg("Good afternoon", &mutex, &wg)
	go newBb.updateMsg("Good night", &mutex, &wg)
	go newBb.updateMsg("Good afternoon", &mutex, &wg)
	go newBb.updateMsg("Good night", &mutex, &wg)
	go newBb.updateMsg("Good afternoon", &mutex, &wg)
	go newBb.updateMsg("Good night", &mutex, &wg)
	go newBb.updateMsg("Good afternoon", &mutex, &wg)
	go newBb.updateMsg("Good night", &mutex, &wg)

	log.Printf("here is msg:%+v", newBb) // 有時會是 afternoon, 有時會是 night，因為哪個 go 會先完成不確定
}
