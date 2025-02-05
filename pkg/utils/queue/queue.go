package queue

type (
	Queue struct {
		last   *node
		first  *node
		length int
	}
	node struct {
		value interface{}
		prev  *node
	}
)

// Create a new queue
func New() *Queue {
	return &Queue{nil, nil, 0}
}

// Return the number of items in the queue
func (q *Queue) Len() int {
	return q.length
}

// Remove the last item of the queue and return it
func (q *Queue) Dequeue() interface{} {
	if q.length == 0 {
		return nil
	}

	n := q.last
	q.last = n.prev

	if q.last == nil {
		q.first = nil
	}

	q.length = q.length - 1
	return n.value
}

// Add a value onto the top of the queue
func (q *Queue) Enqueue(value interface{}) {
	n := &node{value, nil}

	if q.last == nil {
		q.last = n
	} else {
		q.first.prev = n
	}

	q.first = n

	q.length = q.length + 1
}
