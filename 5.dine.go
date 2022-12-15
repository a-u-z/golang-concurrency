package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

type PhilosopherAndForks struct {
	name      string
	rightFork int
	leftFork  int
}

var philosophers = []PhilosopherAndForks{
	{name: "A", leftFork: 4, rightFork: 0},
	{name: "B", leftFork: 0, rightFork: 1},
	{name: "C", leftFork: 1, rightFork: 2},
	{name: "D", leftFork: 2, rightFork: 3},
	{name: "E", leftFork: 3, rightFork: 4},
}

var eatRound = 3                     // how many times a philosopher eats
var eatTime = 500 * time.Millisecond // how long it takes to eatTime
var thinkTime = 2 * time.Second      // how long a philosopher thinks

func dine() {
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))
	// 這是因為到時候會讓每個哲學家都啟動吃東西 func 不過還是要等大家都坐好，才啟動吃東西程序
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// 將每一把叉子都打上可以上鎖的 mutex
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// 讓每一個哲學家都啟動吃東西程序
	for i := 0; i < len(philosophers); i++ {
		go seatAndEat(philosophers[i], wg, forks, seated)
	}
	wg.Wait()
}

// seatAndEat is the function fired off as a goroutine for each of our philosophers. It takes one
// philosopher, our WaitGroup to determine when everyone is done, a map containing the mutexes for every
// fork on the table, and a WaitGroup used to pause execution of every instance of this goroutine
// until everyone is seated at the table.
func seatAndEat(philosopher PhilosopherAndForks, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	// 有把等待的 wg 傳進來，所以先預設這個 func 結束後，wg.Done()，這樣就會減一
	defer wg.Done()

	// 要寫在裡面的原因是 seated.Wait() 要等待的是後面的程序，如果獨立出來一個 func 的話沒有辦法做到等待後面的程序
	seated.Done()
	seated.Wait()

	// 每一個人要吃三輪
	for i := 1; i < eatRound; i++ {

		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			colorPrint(philosopher.name, fmt.Sprintf("\t%s takes fork No.%v\n", philosopher.name, philosopher.rightFork))
			forks[philosopher.leftFork].Lock()
			colorPrint(philosopher.name, fmt.Sprintf("\t%s takes fork No.%v\n", philosopher.name, philosopher.leftFork))

		} else {
			forks[philosopher.leftFork].Lock()
			colorPrint(philosopher.name, fmt.Sprintf("\t%s takes fork No.%v\n", philosopher.name, philosopher.leftFork))
			forks[philosopher.rightFork].Lock()
			colorPrint(philosopher.name, fmt.Sprintf("\t%s takes fork No.%v\n", philosopher.name, philosopher.rightFork))
		}
		colorPrint(philosopher.name, fmt.Sprintf("\t%s eat round %v.\n", philosopher.name, i))

		time.Sleep(eatTime)
		// 哲學家需要思考，不過叉子沒有放下，所以其他人還沒有辦法吃
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()
		colorPrint(philosopher.name, fmt.Sprintf("\t%s put down the forks.\n", philosopher.name))
	}

	fmt.Println(philosopher.name, "is done.")
}

func colorPrint(philosopherName, word string) {
	switch philosopherName {

	case "A":
		color.Cyan(word)
	case "B":
		color.Yellow(word)
	case "C":
		color.Red(word)
	case "D":
		color.Green(word)
	case "E":
		color.Magenta(word)
	}
}
