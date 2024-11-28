package main

import (
	"fmt"
	"gosuperadder/superadder"
	"time"
)

const (
	N_SEQUENCES = 100000
	SEQ_LENGTH_UPPER_BOUND = 100000

	N_WORKERS = 1000
	SUBMISSION_BUF_SIZE = SEQ_LENGTH_UPPER_BOUND / N_WORKERS
	TASK_POOL_BUF_SIZE = SUBMISSION_BUF_SIZE
	DELAY_UPPER_BOUND_MS = 0
)

func naiveAdder(seq []int) int {
	total := 0
	for _, value := range seq {
		total += value
	}
	return total
}

func testNaiveAdder(n int, lengthUpperBound int) {
	for i := 0; i < n; i++ {
		seq := superadder.GenerateRandomSequence(lengthUpperBound)
		naiveAdder(seq)
	}
}

func main() {
	sa := superadder.InitSuperAdder(N_WORKERS, SUBMISSION_BUF_SIZE, TASK_POOL_BUF_SIZE, DELAY_UPPER_BOUND_MS)
	sa.Test(N_SEQUENCES, SEQ_LENGTH_UPPER_BOUND)

	start := time.Now()
	testNaiveAdder(N_SEQUENCES, SEQ_LENGTH_UPPER_BOUND)
	elapsed := time.Since(start)
	fmt.Printf("Naive adder executed %d sequences of length up to %d in %s\n", N_SEQUENCES, SEQ_LENGTH_UPPER_BOUND, elapsed)
}