package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// 學習要點：
// 1. 用 channel 來判斷是否程式執行結束
// 2. 線程的同步與互斥的操作
type BarberShop struct {
	ShopCapacity            int
	HairCutDuration         time.Duration
	NumberOfBarbers         int
	CanBarbersGoHomeChannel chan bool
	ClientsWaitingChannel   chan string
	Open                    bool
	CutCount                int
	ShopClosingChannel      chan bool
}

func (shop *BarberShop) initBarberShop() {
	shop.ShopCapacity = seatingCapacity
	shop.HairCutDuration = cutDuration
	shop.NumberOfBarbers = 0
	shop.ClientsWaitingChannel = make(chan string, seatingCapacity)
	shop.CanBarbersGoHomeChannel = make(chan bool)
	shop.Open = true
	shop.CutCount = 0
	shop.ShopClosingChannel = make(chan bool)
	color.Green("The shop is open for the day!")
}

func (shop *BarberShop) addBarber(barberName string) {
	shop.NumberOfBarbers++
}

func (shop *BarberShop) barberWork(barberName string) {
	isSleeping := false // 初始值是醒著
	color.Yellow("%s goes to the waiting room to check for clients.", barberName)

	for {
		// if there are no clients, the barber goes to sleep
		if len(shop.ClientsWaitingChannel) == 0 {
			color.Yellow("There is nothing to do, so %s takes a nap.", barberName)
			isSleeping = true
		}

		// 從 clients 拉東西出來，有可能 waitingChannel 會是空的，所以要用兩個值檢查
		// 當真的有人在等待的時候，才讓 barber 做事情
		client, foundSomeoneWaiting := <-shop.ClientsWaitingChannel

		if foundSomeoneWaiting {
			// 當 barber 在睡覺，叫醒他
			if isSleeping {
				color.Yellow("%s wakes %s up.", client, barberName)
				isSleeping = false
			}
			// cut hair
			shop.cutHair(barberName, client)
		} else {
			// 其實應該要設定一個時間，當時間到了下班時間才可以回家
			shop.sendBarberHome(barberName)
			return
		}
	}
}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair.", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)
	shop.CutCount++
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.CanBarbersGoHomeChannel <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("Closing shop for the day.")

	close(shop.ClientsWaitingChannel)
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.CanBarbersGoHomeChannel
	}

	close(shop.CanBarbersGoHomeChannel)

	color.Green("---------------------------------------------------------------------")
	color.Green("The barbershop is now closed for the day, and everyone has gone home.")
}

func (shop *BarberShop) addClient(client string) {
	// print out a message
	color.Green("*** %s arrives!", client)

	if shop.Open {
		select {
		case shop.ClientsWaitingChannel <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default:
			color.Red("The waiting room is full, so %s leaves.", client)
		}
	} else {
		color.Red("The shop is already closed, so %s leaves!", client)
	}
}

func (shop *BarberShop) selectShopClosingChannel(shopClosingChannel chan bool, i int) {
	for {
		// get a random number with average arrival rate
		randomMilliseconds := rand.Int() % (2 * arrivalRate)
		select {
		case <-shopClosingChannel:
			// 然後要等到所有等待區的人剪完頭髮，才可以結束
			return // 結束執行這個 func 的意思
		case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
			shop.addClient(fmt.Sprintf("Client #%d", i))
			i++
			// 這邊沒有 default 需要修改
		}
	}
}

func (shop *BarberShop) closeShopAndMainFunc(shopClosingChannel chan bool, mainFuncCloseChannel chan bool) {
	<-time.After(openTime)     // 這個是說這間店開的秒數
	shopClosingChannel <- true // 用這個 channel 來控制是否要讓店休息
	shop.closeShopForDay()
	log.Printf("here is shop:%+v", shop.CutCount)
	mainFuncCloseChannel <- true
}
