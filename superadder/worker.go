package superadder

import (
	"time"
)

type Worker struct {
	id int
	workloads chan *Workload
	delay time.Duration
	die chan bool
}

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