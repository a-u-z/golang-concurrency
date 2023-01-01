package main

import (
	"fmt"

	"github.com/iloveanimal/cutedog/utils"
)

type Car struct {
	Status        bool // 如果這個 status 是 true 代表車子完成
	Name          string
	ExteriorParts bool
	InteriorParts bool
}

// 需要這個 init 嗎？ var 出來的時候不就都是給定一個初始值了嗎？
// func (c *Car) initCar() {
// 	c.Name = ""
// 	c.ExteriorParts = false
// 	c.InteriorParts = false
// 	c.Status = false
// }

var (
	nameChannel          = make(chan Car, 10)
	exteriorPartsChannel = make(chan Car, 10)
	interiorPartsChannel = make(chan Car, 10)
	statusChannel        = make(chan Car, 10)
	doneChannel          = make(chan []Car)
)

func carAssembly() {
	carsCount := 20
	// 要把車子送上固定產線

	go addName()
	go addExteriorParts()
	go addInteriorParts()
	go addStatus(carsCount)
	for {
		select {
		default:
			nameChannel <- Car{}
		case cars := <-doneChannel:
			for i := range cars {
				if !cars[i].ExteriorParts || !cars[i].InteriorParts || !cars[i].Status {
					fmt.Print("something went wrong")
					return
				}
			}
			fmt.Printf("Good job. You made %v cars\n", len(cars))
			close(doneChannel)
			return
		}
	}

}

func addName() {
	for {
		// 還要關閉 channel
		car := <-nameChannel
		car.Name = utils.RandomString(10)
		exteriorPartsChannel <- car
	}
}

func addExteriorParts() {
	for {
		car := <-exteriorPartsChannel
		car.ExteriorParts = true
		interiorPartsChannel <- car
	}
}

func addInteriorParts() {
	for {
		car := <-interiorPartsChannel
		car.InteriorParts = true
		statusChannel <- car
	}
}

func addStatus(carsCount int) {
	var carCollect []Car
	for {
		car := <-statusChannel
		car.Status = true
		carCollect = append(carCollect, car)
		if len(carCollect) == carsCount {
			doneChannel <- carCollect
		}
	}
}
