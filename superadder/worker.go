package superadder

import (
	"time"
)

// Worker represents a worker that processes workloads.
// Each worker has a unique ID, a channel for receiving workloads,
// a delay duration to simulate processing time, and a channel to signal termination.
type Worker struct {
	id int
	workloads chan *Workload
	delay time.Duration
	die chan bool
}

// Run starts the worker's main loop, processing workloads from the workloads channel.
// It listens for two types of events:
// 1. A signal from the die channel, which causes the worker to stop running and return.
// 2. A workload from the workloads channel, which is processed by spawning a new goroutine
//    that calls the work method with the received load.
func (w *Worker) Run() {
	for {
		select {
		case <- w.die:
			return
		case load := <- w.workloads:
			go w.work(load)
		}
	}
}

func (w* Worker) work(load *Workload) {
	time.Sleep(w.delay)
	// fmt.Println("Worker ", w.id, " is working on slice ", load.item, " for sum of V", load.sumId)
	load.dump <- sum(load.item)
}

func sum(slice []int) int {
    total := 0
	for _, value := range slice {
        total += value
    }
    return total
}