// Author: Jasmin A Smith
// Date: 06/20/2020

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Lets tourists use computers
func useComputer(tourist string, leave chan string) {
	fmt.Println(tourist, "is online")
	timeSpent := rand.Intn(105) + 15
	time.Sleep(time.Duration(timeSpent) * time.Millisecond * 10)
	fmt.Println(tourist, "is done, having spent", timeSpent, "minutes online.")
	leave <- tourist
}

// make sure no more than 8 users, and keeps going until all tourist have gone
func manageTourists(enter, leave chan string, stop chan struct{}) {
	notUsed := 0
	line := make([]string, 0)
	for {
		select {
		case tourist := <-enter:
			if notUsed < 8 {
				notUsed++
				go useComputer(tourist, leave)
			} else {
				time.Sleep(10 * time.Millisecond)
				fmt.Println(tourist, "waiting for turn.")
				line = append(line, tourist)
			}
		case <-leave:
			notUsed--
			if len(line) > 0 {
				next := line[0]
				line = line[1:]
				go func() {
					enter <- next
				}()
			} else if notUsed == 0 {
				close(stop)
				return
			}

		}
	}
}

// Shuffles the slice to make sure there was a random order
func Shuffle(slice []int) []int {

	newSlice := make([]int, len(slice))
	length := len(slice)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(slice))
		newSlice[i] = slice[randomIndex]
		slice = append(slice[:randomIndex], slice[randomIndex+1:]...)
	}
	return newSlice
}
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	slice := make([]int, 25)
	for i := 0; i < 25; i++ {
		slice[i] = i + 1
	}

	slice = Shuffle(slice)
	enter := make(chan string)
	leave := make(chan string)
	stop := make(chan struct{})
	go manageTourists(enter, leave, stop)

	for i := 0; i < 25; i++ {
		enter <- "Tourist " + strconv.Itoa(slice[i])
	}
	<-stop
	fmt.Println("The place is empty, let's close up and go to the beach!")
}
