package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

const (
	seatingCapacity = 100
	arrivalRate     = 100
	cutDuration     = 1000 * time.Millisecond
	timeOpen        = 10 * time.Second
)

func main() {
	rand.Seed(time.Now().UnixNano())
	color.HiGreen("Welcome to the Barber Shop\n")
	color.HiGreen("--------------------------\n")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		BarbersDoneChan: doneChan,
		ClientsChan:     clientChan,
		Open:            true,
	}

	color.Yellow("The shop is open for the day")
	shop.addBarber("frank")
	shop.addBarber("john")
	time.Sleep(5 * time.Second)

	shopClosing := make(chan bool)
	closed := make(chan bool)
	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()
	i := 1
	go func() {
		for {
			randomMillisecond := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMillisecond)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()
	<-closed
}
