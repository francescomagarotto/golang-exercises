package main

import (
	"github.com/fatih/color"
	"time"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barberName string) {
	shop.NumberOfBarbers++
	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients", barberName)
		for {
			if len(shop.ClientsChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap", barberName)
				isSleeping = true
			}
			client, shopOpen := <-shop.ClientsChan
			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barberName)
				}
				shop.cutHair(barberName, client)
			} else {
				shop.sendBarberHome(barberName)
				return
			}
		}
	}()
}

func (b *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair.", barber, client)
	time.Sleep(b.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)
}
func (b *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home", barber)
	b.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("Closing shop for the day")
	close(shop.ClientsChan)
	for i := 0; i < shop.NumberOfBarbers; i++ {
		<-shop.BarbersDoneChan
	}
	close(shop.BarbersDoneChan)
	color.Green("The shop is closed for the day")
}

func (shop *BarberShop) addClient(clientName string) {
	color.Green("Client %s arrived", clientName)
	if shop.Open {
		select {
		case shop.ClientsChan <- clientName:
			color.Yellow("%s takes a seat in the waiting room ", clientName)
		default:
			color.Red("The waiting room is full, so %s leaves", clientName)
		}
	} else {
		color.Red("The shop is closed. %s leaves", clientName)
	}
}
