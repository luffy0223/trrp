package main

import (
	"fmt"
	"time"
)

var highChan = make(chan bool)
var lowChan = make(chan bool)

func main() {
	fmt.Println("start...")
	go changeLowChan()
	go changeHighChan()
LOOP:
	for {
		select {
		case <-highChan:
			fmt.Print("highChan is true")
			break LOOP
		case <-lowChan:
			fmt.Println("lowChan is true")
			for {
				select {
				case <-highChan:
					fmt.Print("highChan is true")
					break LOOP
				default:
					break
				}
				time.Sleep(time.Second)
				fmt.Println("do something")

			}

		}
	}
}

func changeHighChan() {
	time.Sleep(time.Second * 10)
	highChan <- true
}

func changeLowChan() {
	time.Sleep(time.Second * 5)
	lowChan <- true

}
