# SuperAdder

SuperAdder is a Go project that demonstrates concurrent processing of workloads using a pool of worker goroutines. The project generates random sequences of integers, distributes the work among multiple workers, and accumulates the results. Each worker does not know which global sum it is working on.

## Project Structure

The project is organized into the following files and directories:

- `main.go`: The main entry point of the application that initializes the pool of workers, generates random sequences of integers, and dispatches the work to the workers.
- `worker.go`: The worker package that defines the worker pool and the worker goroutines that process the work.
- `accumulator.go`: The accumulator package that defines the accumulator goroutine that accumulates the results from the workers.
- `README.md`: The project README file that provides an overview of the project and its structure.