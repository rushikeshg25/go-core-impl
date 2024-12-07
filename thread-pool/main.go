package main

import (
	"fmt"
	"sync"
	"time"
)

type Task func()

type Pool struct{
	workers chan Task
	wg sync.WaitGroup
}

func InitPool(workersCnt int)Pool{
	pool:=Pool{
		workers: make(chan Task),
	}
	pool.wg.Add(workersCnt);
	for i:=0;i<workersCnt;i++{
		go func(){
			defer pool.wg.Done();
			for job:=range pool.workers{
				job()
			}
		}()
	}
	return pool;
}

func (p *Pool) AddTask(task Task){
	p.workers<-task
}

func (p *Pool) Wait(){
	close(p.workers)
	p.wg.Wait()
}

func main(){
	pool:=InitPool(5);

	for i:=0;i<30;i++{
		task:=func(){
			time.Sleep(1*time.Second)
			fmt.Println("Task Done!")
		}
		pool.AddTask(task)
	}
	pool.Wait()
}