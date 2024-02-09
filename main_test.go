package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestSlidingWindow(t *testing.T) {
	counter := &CounterData{
		Count:      0,
		StartTimes: make([]time.Time, 0),
		WindowSize: windowSize,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter.mu.Lock()
		defer counter.mu.Unlock()
		counter.Count++
		counter.StartTimes = append(counter.StartTimes, time.Now())
		fmt.Fprintf(w, "Total requests in the last %d seconds: %d", windowSize, counter.Count)
	}))
	defer ts.Close()

	for i := 0; i < windowSize; i++ {
		res, err := http.Get(ts.URL)
		if err != nil {
			t.Fatalf("Error making HTTP request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		expectedBody := fmt.Sprintf("Total requests in the last %d seconds: %d", windowSize, counter.Count)
		actualBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}

		if string(actualBody) != expectedBody {
			t.Errorf("Expected response body %s, got %s", expectedBody, string(actualBody))
		}

		time.Sleep(time.Second)
	}

	if counter.Count > windowSize {
		t.Errorf("Expected count not to exceed %d, got %d", windowSize, counter.Count)
	}

	if len(counter.StartTimes) > windowSize {
		t.Errorf("Expected start times length not to exceed %d, got %d", windowSize, len(counter.StartTimes))
	}

	if counter.Count != windowSize {
		t.Errorf("Expected count %d, got %d", windowSize, counter.Count)
	}
}

func TestSaveLoadCounterData(t *testing.T) {
	testFilePath := "test_request_counter.gob"
	defer os.Remove(testFilePath)

	counter := &CounterData{
		Count:      42,
		StartTimes: []time.Time{time.Now()},
		WindowSize: windowSize,
	}

	saveCounterData(counter)
	loadedCounter := loadCounterData()
	if loadedCounter.Count != counter.Count {
		t.Errorf("Expected count %d, got %d", counter.Count, loadedCounter.Count)
	}

	if len(loadedCounter.StartTimes) != 1 || !loadedCounter.StartTimes[0].Equal(counter.StartTimes[0]) {
		t.Errorf("Expected start times %v, got %v", counter.StartTimes, loadedCounter.StartTimes)
	}

	if loadedCounter.WindowSize != counter.WindowSize {
		t.Errorf("Expected window size %d, got %d", counter.WindowSize, loadedCounter.WindowSize)
	}
}
