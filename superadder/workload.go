package superadder

type Workload struct {
	item  []int
	sumId int
	dump  chan int
}