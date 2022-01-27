package util

type CircleQueue struct {
	data     []interface{}
	capacity int
	head     int
	tail     int
}

func NewCircleQueue(cap int) *CircleQueue {
	if cap <= 0 {
		return nil
	}

	return &CircleQueue{
		data:     make([]interface{}, cap),
		capacity: cap,
		head:     0,
		tail:     0,
	}
}

func (q *CircleQueue) IsEmpty() bool {
	return q.head == q.tail
}

func (q *CircleQueue) InitQueue() {
	q.head = 0
	q.tail = 0
}

func (q *CircleQueue) OverwriteTail(v interface{}) {
	if q.IsEmpty() {
		return
	}

	q.data[q.tail] = v
}

//큐 갯수가 가득 찼다면 head를 빼고 추가
func (q *CircleQueue) Enqueue(v interface{}) interface{} {
	var h interface{}
	if q.isFull() {
		h = q.Dequeue()
	}

	q.data[q.tail] = v
	q.tail = q.calcPos(q.tail)

	return h
}

func (q *CircleQueue) Dequeue() interface{} {
	if q.IsEmpty() {
		return nil
	}

	v := q.data[q.head]
	q.head = q.calcPos(q.head)
	return v
}

func (q *CircleQueue) PeekHead() interface{} {
	return q.peek(q.head)
}
func (q *CircleQueue) PeekTail() interface{} {
	return q.peek((q.tail - 1) % q.capacity)
}

func (q *CircleQueue) Count() int {
	return (q.capacity + q.tail - q.head) % q.capacity
}

func (q *CircleQueue) ExportToArray() []interface{} {
	if q.IsEmpty() {
		return make([]interface{}, 0)
	}

	arr := make([]interface{}, q.Count())
	pos := q.head

	for ii := 0; ii < q.Count(); ii++ {
		arr[ii] = q.data[pos]
		pos = q.calcPos(pos)
	}

	return arr
}

func (q *CircleQueue) PeekArrayFromTail(cnt int) []interface{} {
	if q.IsEmpty() {
		return nil
	}

	if q.Count() < cnt {
		cnt = q.Count()
	}

	arr := make([]interface{}, cnt)
	pos := q.calcPos(q.head + (q.Count() - cnt - 1))

	for ii := cnt - 1; ii >= 0; ii-- {
		arr[ii] = q.data[pos]

		pos = q.calcPos(pos)

	}

	return arr
}

//private
func (q *CircleQueue) calcPos(cur int) int {
	return (cur + 1) % q.capacity
}

func (q *CircleQueue) isFull() bool {
	return q.head == q.calcPos(q.tail)
}

func (q *CircleQueue) peek(pos int) interface{} {
	if q.IsEmpty() {
		return nil
	}

	if pos < 0 {
		pos = q.head
	}
	return q.data[pos]
}
