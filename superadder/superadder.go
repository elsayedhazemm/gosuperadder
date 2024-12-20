package superadder

import (
	"fmt"
	"math/rand"
	"time"
)

type Result struct {
	SumId int
	Value int
}

type Pool []*Worker

// SuperAdder is a struct that manages a pool of workers to process workloads
// and produce results. It includes a channel for workloads, a channel for
// results, and an upper bound for delay in milliseconds.
//
// Fields:
// - pool: A Pool that manages the worker goroutines.
// - workloads: A channel through which Workload pointers are sent to be processed. In this context, workloads are slices of integers.
// - Results: A channel through which Result pointers are sent after processing.
// - delayUpperBound: An integer representing the upper bound for delay in milliseconds.
type SuperAdder struct {
	pool      Pool
	workloads chan *Workload
	Results   chan *Result
	delayUpperBound int // in ms
}

func InitSuperAdder(nWorkers int, submissionBufSize int, taskPoolBufSize int, delayUpperBoundMs int) *SuperAdder {
	pool := make([]*Worker, nWorkers)

	// print configuration
	fmt.Println("SuperAdder initialized with ", nWorkers, " workers")
	fmt.Println("Submission buffer size: ", submissionBufSize)
	fmt.Println("Task pool buffer size: ", taskPoolBufSize)
	fmt.Println("Delay upper bound: ", delayUpperBoundMs, "ms")

	workloads := make(chan *Workload, submissionBufSize)
	results := make(chan *Result, taskPoolBufSize)

	sa := &SuperAdder{
		pool:            pool,
		workloads:       workloads,
		Results:         results,
		delayUpperBound: delayUpperBoundMs,
	}

	sa.expandWorkerPool(nWorkers)
	return sa
}


// Test runs the SuperAdder with nSequences sequences of length up to seqLengthUpperBound.
func (sa *SuperAdder) Test(nSequences int, seqLengthUpperBound int) {
	fmt.Println("Testing Super Adder with ", nSequences, " sequences of length up to ", seqLengthUpperBound)
	start := time.Now()

	seqChan := make(chan []int, nSequences)
	receivedAll := make(chan bool)

	go receiveResults(sa.Results, nSequences, receivedAll)
	go generateSequences(seqChan, nSequences, seqLengthUpperBound)

	go func() {
		i := 0
		for seq := range seqChan {
			go sa.Sum(seq, i)
			i++
		}
	}()

	<-receivedAll
	sa.ShutDown()

	elapsed := time.Since(start)
	fmt.Printf("Super Adder Shut Down\n")
	fmt.Printf("Super Adder Execution time: %s\n", elapsed)
}

func (sa *SuperAdder) ShutDown() {
	for _, w := range sa.pool {
		if w != nil {
			w.die <- true
		}
	}

	close(sa.workloads)
	close(sa.Results)
}

// Sum divides the input slice into batches and distributes the work among the available workers
// in the SuperAdder pool. It then accumulates the partial sums from each batch into a final sum.
//
// Parameters:
//   v - A slice of integers to be summed.
//   sumId - An identifier for the sum operation.
//
// Panics:
//   If there are no workers available in the SuperAdder pool, it will panic with the message "WORKERS ARE ALL DEAD!".
func (sa *SuperAdder) Sum(v []int, sumId int) {		
	if len(sa.pool) == 0 {
		panic("WORKERS ARE ALL DEAD!")
	}

	batchSize := len(v)  + 1 / len(sa.pool)
	numBatches := (len(v) + batchSize - 1) / batchSize

	partialSums := make(chan int, numBatches)
	accumulated := make(chan bool)
	go sa.accumulate(sumId, partialSums, accumulated)
	
	go func() {
		<-accumulated
		close(accumulated)
		close(partialSums)
	}()

	for i := 0; i < len(v); i += batchSize {
		end := i + batchSize
		if end > len(v) {
			end = len(v)
		}

		// Send workload to task pool
		sa.workloads <- &Workload{
			item: v[i:end],
			dump: partialSums,
			sumId: sumId,
		}
	}
}

// accumulate listens for n partial sums on the channel of sequence v and sends the total sum to the results channel.
func (sa *SuperAdder) accumulate(sumId int, partialSums chan int, done chan bool) {
	total, c := 0, 0
	n := cap(partialSums)

	// listen for n partial sums on the channel of sequence v
	for c < n {
		psum, ok := <- partialSums
		if !ok {
			break
		}
		total += psum
		c++
	}

	// Send the total sum to the results channel
	sa.Results <- &Result{
		SumId: sumId,
		Value: total,
	}
	done <- true
}


func (sa *SuperAdder) expandWorkerPool(n int) {
	var delay time.Duration
	for i := 0; i < n; i++ {
		if sa.delayUpperBound > 0 {
			delay = time.Duration(rand.Intn(sa.delayUpperBound)) * time.Millisecond
		} else {
			delay = 0
		}
		w := &Worker{id: i, workloads: sa.workloads, die: make(chan bool), delay: delay}
		sa.pool = append(sa.pool, w)
		go w.Run()
	}
}


// receiveResults reads n results from the channel and closes the done channel
func receiveResults(ch chan *Result, n int, done chan bool) {
	for i := 0; i < n; i++ {
		<- ch
		// res := <-ch
		// fmt.Println("Sum of V", res.SumId, ": ", res.Value)
	}
	done <- true
}

// GenerateRandomSequence generates a random sequence of integers with a length up to lengthUpperBound.
func GenerateRandomSequence(lengthUpperBound int) []int {
	size := rand.Intn(lengthUpperBound) + 1
	seq := make([]int, size)
	for j := 0; j < size; j++ {
		seq[j] = rand.Intn(1000)
	}
	return seq
}

func generateSequences(ch chan []int, nSequences int, seqLengthUpperBound int) {
	for i := 0; i < nSequences; i++ {
		seq := GenerateRandomSequence(seqLengthUpperBound)
		ch <- seq
		// fmt.Println("V", i, ": ", seq)
	}
}