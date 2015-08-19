package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const numGoroutines = 10
const maxSleepTimeSecs = 1

var urls = []string{
	"https://groupon.com",
	"https://groupon.com/getaways",
	"https://groupon.com/goods",
	"http://google.com",
	"http://amazon.com",
	"http://yahoo.com",
	"http://microsoft.com",
	"http://example.com",
	"http://unreachable.com/noexist",
}

func main() {
	demo1()
	// demo2()
	// demo3()
}

func demo1() {
	ch := make(chan string)

	for _, url := range urls {
		go func(url string) {
			ch <- httpGetCheck(url)
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Println(<-ch)
	}

	fmt.Println("Done!")
}

func httpGetCheck(url string) string {
	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		return fmt.Sprintf("*** (%v) is NOT 200 OK", url)
	}

	return fmt.Sprintf("    (%v) is 200 OK", url)
}

func demo2() {
	ch := make(chan string)

	for _, url := range urls {
		go func(url string) {
			ch <- httpGetCheck(url)
		}(url)
	}

	secondMarks := ticker(1 * time.Second)
	timeout := time.After(6 * time.Second)

	done := false
	for messagesReceived := 0; messagesReceived < len(urls) && !done; {
		select {
		case message := <-ch:
			fmt.Println(message)
			messagesReceived++
		case tick := <-secondMarks:
			fmt.Printf("======== %v =======\n", tick)
		case <-timeout:
			fmt.Println("******** Bored ********")
			done = true
		}
	}

	fmt.Println("Done!")
}

func ticker(period time.Duration) chan int {
	ch := make(chan int)

	go func() {
		i := 0
		for {
			ch <- i
			time.Sleep(period)
			i++
		}
	}()

	return ch
}

func demo3() {
	jeff := steadySpeaker(400*time.Millisecond, "jeff")
	sara := steadySpeaker(600*time.Millisecond, "\tsara")
	pickles := steadySpeaker(1000*time.Millisecond, "\t\tpickles")
	random := randomSpeaker()
	timeout := time.After(5 * time.Second)

	done := false
	for !done {
		select {
		case phrase := <-jeff:
			fmt.Printf("phrase: %v\n", phrase)
		case phrase := <-sara:
			fmt.Printf("phrase: %v\n", phrase)
		case phrase := <-pickles:
			fmt.Printf("phrase: %v\n", phrase)
		case phrase := <-random:
			fmt.Printf("phrase: \t\t\t\t%v\n", phrase)
		case <-timeout:
			done = true
		}
	}

	fmt.Println("Enough!")
}

func steadySpeaker(period time.Duration, phrase string) chan string {
	ch := make(chan string)

	go func() {
		for {
			time.Sleep(period)
			ch <- phrase
		}
	}()

	return ch
}

func randomSpeaker() chan int {
	ch := make(chan int)

	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			ch <- i
		}
	}()

	return ch
}
