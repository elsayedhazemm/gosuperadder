package superadder

// Workload represents a structure that holds a slice of integers, for workers to use
// an identifier for the global sum operation the slice is under,
// and a channel corresponding to sumId for dumping the sum of this slice.
type Workload struct {
	item  []int
	sumId int
	dump  chan int
}