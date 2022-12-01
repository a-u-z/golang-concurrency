package main

import (
	"log"
	"sync"
)

type aa struct {
	msg string
}

// 這樣會有 race condition 產生， go run -race .，可以偵測出來
func raceCondition() {
	var wg sync.WaitGroup
	newAa := &aa{
		msg: "Good morning",
	}

	wg.Add(12)
	go newAa.updateMsg("Good afternoon", &wg)
	go newAa.updateMsg("Good night", &wg)
	go newAa.updateMsg("Good afternoon", &wg)
	go newAa.updateMsg("Good night", &wg)
	go newAa.updateMsg("Good afternoon", &wg)
	go newAa.updateMsg("Good night", &wg)
	go newAa.updateMsg("Good afternoon", &wg)
	go newAa.updateMsg("Good night", &wg)
	go newAa.updateMsg("Good afternoon", &wg)
	go newAa.updateMsg("Good night", &wg)
	go newAa.updateMsg("Good afternoon", &wg)
	go newAa.updateMsg("Good night", &wg)

	log.Printf("here is msg:%+v", newAa) // 有時會是 afternoon, 有時會是 night
}

func (a *aa) updateMsg(s string, wg *sync.WaitGroup) string {
	defer wg.Done()
	a.msg = s
	return s
}
