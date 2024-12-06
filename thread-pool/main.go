package main

import "sync"


type Job func()

type Pool struct{
	workQueue chan Job
	wg sync.WaitGroup
}

func NewPool(worker int)*Pool{
	pool:=Pool{
	workQueue: make(chan Job),
	}
	pool.wg.Add(worker);
	for i:=0;i<worker;i++{
		go func(){
			defer pool.wg.Done()
			
		}()
	}
	return &pool;
}

func main(){
	pool:=NewPool

}