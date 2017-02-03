package main

import (
	"github.com/codahale/hdrhistogram"
	"fmt"
	"bufio"
	"os"
	"regexp"
	"strconv"
	"os/signal"
	"syscall"
	"time"
)

var (
	hist *hdrhistogram.Histogram
	q50, q75, q90, q95, q99 int64
)

func report() {
	q99 = hist.ValueAtQuantile(99.0)
	q95 = hist.ValueAtQuantile(95.0)
	q90 = hist.ValueAtQuantile(90.0)
	q75 = hist.ValueAtQuantile(75.0)
	q50 = hist.ValueAtQuantile(50.0)
	fmt.Printf("\n\nQuantiles: %d %d %d %d %d", q99, q95, q90, q75, q50)
}

func reporter(c chan os.Signal, d chan int64) {
	ticker := time.NewTicker(5 * time.Second)
	var latency int64
	for {
		select {
		case latency = <-d:
			hist.RecordValue(latency)

		case <-ticker.C:
			report()
			hist.Reset()

		case sig := <-c:
			fmt.Println(sig)
		}
	}

}

func scanner(exp string, d chan int64) {
	// max assumed latency: 10 sec, min assumed latency: 1 usec
	hist = hdrhistogram.New(1, 10000000, 2)

	in := bufio.NewScanner(os.Stdin)
	re := regexp.MustCompile(exp)

	for in.Scan() {
		for _, match := range re.FindAllStringSubmatch(in.Text(), -1) {
			//fmt.Printf("%s,", match[1])
			latency, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				fmt.Println("Conversion error")
			}
			d <- latency
		}

	}
	report()
}

func main() {
	sigs := make(chan os.Signal, 1)
	data := make(chan int64)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go reporter(sigs, data)
	scanner(`(\d+)usec`, data)
}
