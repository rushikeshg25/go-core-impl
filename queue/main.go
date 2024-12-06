package main

type Queue struct{
	queue []int32;
}


func (q *Queue) Enqueue(v int32){
	q.queue=append(q.queue, v)
}

func (q *Queue) Dequeue() int32{
	if(len(q.queue)==0){
		panic("Queue is empty")
	}
	temp:=q.queue[0];
	q.queue=q.queue[1:];
	return temp;
}


func main() {
	q:=Queue{
		queue:make([]int32,0),
	}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(4)
	q.Enqueue(5)
	q.Enqueue(6)
	q.Enqueue(7)
	q.Enqueue(8)
	q.Enqueue(9)
	for i:=0;i<10;i++{
		println(q.Dequeue())
	}

}