package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const baseURL = "http://localhost:8088/api"

type Result struct {
	Endpoint    string
	AvgTime     time.Duration
	MinTime     time.Duration
	MaxTime     time.Duration
	Requests    int
	Errors      int
	ReqPerSec   float64
}

func benchmark(name, url string, count int) Result {
	var wg sync.WaitGroup
	durations := make([]time.Duration, count)
	errors := 0
	client := &http.Client{Timeout: 30 * time.Second}

	start := time.Now()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			reqStart := time.Now()
			resp, err := client.Get(url)
			durations[idx] = time.Since(reqStart)
			if err != nil || resp.StatusCode != 200 {
				errors++
				return
			}
			resp.Body.Close()
		}(i)

		if i%10 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	wg.Wait()
	totalTime := time.Since(start)

	var total, min, max time.Duration
	min = time.Hour
	for _, d := range durations {
		if d > 0 {
			total += d
			if d < min {
				min = d
			}
			if d > max {
				max = d
			}
		}
	}

	avg := total / time.Duration(count)
	rps := float64(count) / totalTime.Seconds()

	return Result{
		Endpoint:  name,
		AvgTime:   avg,
		MinTime:   min,
		MaxTime:   max,
		Requests:  count,
		Errors:    errors,
		ReqPerSec: rps,
	}
}

func main() {
	fmt.Println("🚀 GoFiver Benchmark")
	fmt.Println("====================")

	endpoints := []struct {
		name string
		url  string
	}{
		{"GET /blogs (paginated)", baseURL + "/blogs?page=1&limit=20"},
		{"GET /blogs/:id", baseURL + "/blogs/1"},
		{"GET /users", baseURL + "/users"},
		{"GET /users/:id", baseURL + "/users/1"},
	}

	results := []Result{}

	for _, ep := range endpoints {
		fmt.Printf("\n⏱️  Testing %s...\n", ep.name)
		r := benchmark(ep.name, ep.url, 50)
		results = append(results, r)
		fmt.Printf("   Avg: %v | Min: %v | Max: %v | RPS: %.1f\n",
			r.AvgTime.Round(time.Millisecond),
			r.MinTime.Round(time.Millisecond),
			r.MaxTime.Round(time.Millisecond),
			r.ReqPerSec)
	}

	fmt.Println("\n📊 Summary")
	fmt.Println("==========")
	fmt.Printf("%-30s %12s %12s %12s %10s\n", "Endpoint", "Avg", "Min", "Max", "RPS")
	fmt.Println("---------------------------------------------------------------------------")
	for _, r := range results {
		fmt.Printf("%-30s %12v %12v %12v %10.1f\n",
			r.Endpoint,
			r.AvgTime.Round(time.Millisecond),
			r.MinTime.Round(time.Millisecond),
			r.MaxTime.Round(time.Millisecond),
			r.ReqPerSec)
	}

	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Printf("\nJSON: %s\n", string(out))
}
