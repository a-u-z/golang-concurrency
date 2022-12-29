package main

import (
	"math/rand"
	"time"
)

// variables
var seatingCapacity = 5
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var openTime = 5 * time.Second

func barberMain() {
	// seed our random number generator
	rand.Seed(time.Now().UnixNano())

	// 用 init 實例出 barberShop
	var shop = BarberShop{}
	shop.initBarberShop()

	// 用兩個 channel 來控制 shop 跟 main func 是否需要關閉
	mainFuncCloseChannel := make(chan bool)
	go shop.closeShopAndMainFunc(shop.ShopClosingChannel, mainFuncCloseChannel)

	// 新增每一位理髮師
	// 到這個步驟就可以營業了
	shop.addBarber("Frank")
	go shop.barberWork("Frank")

	// 新增需要理髮的客戶
	i := 1 // 客戶編號
	go shop.selectShopClosingChannel(shop.ShopClosingChannel, i)

	// 這個是把這個 main func 卡住的 channel
	<-mainFuncCloseChannel
}

// goroutine
// 要馬是 select 然後裡面有 default
// 要馬是 for 裡面，不過有 return 機制
// 要馬是會接收到結束訊號
// go 的 func 其實放在哪邊沒有關係，只有他的參數要先宣告好就好
