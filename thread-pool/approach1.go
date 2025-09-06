package main

import (
	"fmt"
	"time"
)

func worker(workerId int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		time.Sleep(time.Second * 2)
		res := job * 2
		results <- res
		fmt.Printf("Worker %d did job %d\n", workerId, job)
	}

}

func main() {
	workerNum := 3
	jobNum := 10
	jobs := make(chan int, jobNum)
	results := make(chan int, jobNum)
	for i := 0; i < workerNum; i++ {
		go worker(i, jobs, results)
	}
	for i := 0; i < jobNum; i++ {
		jobs <- i
	}
	close(jobs)
	for i := 0; i < jobNum; i++ {
		res := <-results
		fmt.Println(res)
	}
	fmt.Println("done")
}
