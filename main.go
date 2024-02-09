package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	dataFilePath = "request_counter.gob"
	windowSize   = 60
)

type CounterData struct {
	Count      int
	StartTimes []time.Time
	WindowSize int
	mu         sync.Mutex
}

func main() {
	counter := loadCounterData()
	go cleanupOldCounters(counter)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter.mu.Lock()
		defer counter.mu.Unlock()

		counter.Count++
		counter.StartTimes = append(counter.StartTimes, time.Now())
		fmt.Fprintf(w, "Total requests in the last %d seconds: %d", windowSize, counter.Count)
	})

	// Handle interrupt signal to save counter data before exiting
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		saveCounterData(counter)
		os.Exit(0)
	}()

	port := 8080
	fmt.Printf("Server listening on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func loadCounterData() *CounterData {
	file, err := os.Open(dataFilePath)
	if err != nil {
		return &CounterData{
			Count:      0,
			StartTimes: make([]time.Time, 0),
			WindowSize: windowSize,
		}
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var counter CounterData
	err = decoder.Decode(&counter)
	if err != nil {
		fmt.Println("Error decoding counter data:", err)
		return &CounterData{
			Count:      0,
			StartTimes: make([]time.Time, 0),
			WindowSize: windowSize,
		}
	}

	return &counter
}

func cleanupOldCounters(counter *CounterData) {
	for {
		time.Sleep(time.Second)

		counter.mu.Lock()

		currentTime := time.Now()
		for i := 0; i < len(counter.StartTimes); {
			if currentTime.Sub(counter.StartTimes[i]).Seconds() >= float64(windowSize) {
				counter.Count--
				counter.StartTimes = append(counter.StartTimes[:i], counter.StartTimes[i+1:]...)
			} else {
				i++
			}
		}

		saveCounterData(counter)
		counter.mu.Unlock()
	}
}

func saveCounterData(counter *CounterData) {
	file, err := os.Create(dataFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(counter)
	if err != nil {
		fmt.Println("Error encoding counter data:", err)
	}
}
