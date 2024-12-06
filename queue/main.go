package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Queue struct{
	queue []int32;
}

func (q *Queue) Enqueue(v int32,mu *sync.Mutex){
	mu.Lock();
	q.queue=append(q.queue, v)
	mu.Unlock();
}

func (q *Queue) Dequeue(mu *sync.Mutex) int32{
	if(len(q.queue)==0){
		panic("Queue is empty")
	}
	mu.Lock();
	temp:=q.queue[0];
	q.queue=q.queue[1:];
	mu.Unlock();
	return temp;
}


func main() {
	q:=Queue{
		queue:make([]int32,0),
	}
	var mu sync.Mutex;
	var wg sync.WaitGroup;
	for i:=0;i<10000;i++{
		wg.Add(1);
		go func(){
			defer wg.Done();
			q.Enqueue(rand.Int31(),&mu);
			
		}();
	}
	for i:=0;i<1000;i++{
		wg.Add(1);
		go func(){
			defer wg.Done();
			q.Dequeue(&mu);
		}()
	}
	wg.Wait();
	fmt.Println("Done");
	fmt.Println(len(q.queue));

}
