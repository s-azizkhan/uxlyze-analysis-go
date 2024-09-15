package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
	process_worker "uxlyze/analyzer/pkg/process-worker"

	"github.com/joho/godotenv"
)

// Job represents a simple job with ID and Payload
type Job struct {
	ID      int32  `json:"id"`
	Payload string `json:"payload"`
}

// ResourceUsage logs memory usage before and after processing a job
type ResourceUsage struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
}

// RateLimiter struct to control job processing rate
type RateLimiter struct {
	jobsProcessed int
	maxJobs       int
	resetTime     time.Time
}

var jobCounter int32               // Atomic counter for job IDs
var jobQueue = make(chan Job, 100) // Buffered channel for job queue

// NewRateLimiter creates a new rate limiter with maxJobs allowed per minute
func NewRateLimiter(maxJobsPerMinute int) *RateLimiter {
	return &RateLimiter{
		jobsProcessed: 0,
		maxJobs:       maxJobsPerMinute,
		resetTime:     time.Now().Add(time.Minute),
	}
}

// CanProcess returns true if a job can be processed, otherwise false
func (rl *RateLimiter) CanProcess() bool {
	// If the reset time has passed, reset the counter
	if time.Now().After(rl.resetTime) {
		rl.jobsProcessed = 0
		rl.resetTime = time.Now().Add(time.Minute)
	}

	// Allow processing if jobsProcessed is less than maxJobs
	if rl.jobsProcessed < rl.maxJobs {
		rl.jobsProcessed++
		return true
	}
	return false
}

// logResourceUsage logs the memory usage before and after processing
func logResourceUsage() ResourceUsage {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return ResourceUsage{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
	}
}

// Worker processes jobs from the job queue, logs time and resource usage
func worker(id int, rateLimiter *RateLimiter) {
	for job := range jobQueue {
		// Wait until we can process the job according to the rate limiter
		for !rateLimiter.CanProcess() {
			time.Sleep(time.Second) // Check again in a second
		}

		// Log the start time and resource usage before processing the job
		startTime := time.Now()
		startResource := logResourceUsage()

		// Process the job
		fmt.Printf("Worker %d processing job ID: %d with payload: %s\n", id, job.ID, job.Payload)
		// job processing
		process_worker.AnalyzeReportWorker(job.Payload)

		// Log the end time and resource usage after processing the job
		endTime := time.Now()
		endResource := logResourceUsage()

		// Calculate the time taken and resource usage
		duration := endTime.Sub(startTime)
		allocDiff := endResource.Alloc - startResource.Alloc
		totalAllocDiff := endResource.TotalAlloc - startResource.TotalAlloc
		sysDiff := endResource.Sys - startResource.Sys

		// Log the details
		fmt.Printf("Worker %d finished job ID: %d\n", id, job.ID)
		fmt.Printf("Job ID: %d took %v to process\n", job.ID, duration)
		fmt.Printf("Memory usage - Alloc: %d bytes, TotalAlloc: %d bytes, Sys: %d bytes\n", allocDiff, totalAllocDiff, sysDiff)
	}
}

// SubmitJobHandler handles job submissions via API
func SubmitJobHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON payload
	var requestBody struct {
		Payload string `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil || requestBody.Payload == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Increment job ID and create a new job
	jobID := atomic.AddInt32(&jobCounter, 1)
	job := Job{
		ID:      jobID,
		Payload: requestBody.Payload,
	}

	// Send the job to the queue
	jobQueue <- job

	// Respond with a confirmation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Job added to queue",
		"job_id":  job.ID,
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
	log.Println("Starting UI/UX analysis server...")
	// Define how many jobs should be processed per minute
	jobsPerMinute := 5

	// Create the rate limiter
	rateLimiter := NewRateLimiter(jobsPerMinute)

	// Start a single worker with the rate limiter
	go worker(1, rateLimiter)

	// Handle the job submission endpoint
	http.HandleFunc("/submit-job", SubmitJobHandler)

	// Start the HTTP server
	port := "8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
