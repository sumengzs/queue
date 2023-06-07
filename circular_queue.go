/*
Copyright 2023 The Queue Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package queue

import (
	"runtime"
	"sync/atomic"
)

type CircularQueue struct {
	items  []interface{}
	size   uint32
	len    uint32
	offset uint32
}

func NewCircularQueue(size uint32) *CircularQueue {
	return &CircularQueue{
		items:  make([]interface{}, size),
		size:   size,
		len:    0,
		offset: 0,
	}
}

func (q *CircularQueue) Len() uint32 {
	return q.len
}

func (q *CircularQueue) Cap() uint32 {
	return q.size
}

func (q *CircularQueue) Get() interface{} {
	if q.len == 0 {
		return nil
	}
	if q.offset == 0 {
		return q.items[q.size-1]
	}
	return q.items[q.offset-1]
}

func (q *CircularQueue) GetPoint(point uint32) interface{} {
	if q.len == 0 || q.len <= point || q.size <= point || point < 0 {
		return nil
	}
	point = (q.offset + q.size - point - 1) % q.size
	return q.items[point]
}

func (q *CircularQueue) Gets(size uint32) []interface{} {
	if q.len == 0 {
		return nil
	}
	if q.len < size {
		size = q.len
	}
	data := q.GetAll()
	data = data[uint32(len(data))-size:]
	return data
}

func (q *CircularQueue) GetAll() []interface{} {
	data := make([]interface{}, q.len)
	if q.len < q.size {
		copy(data, q.items[:q.len])
	} else {
		copy(data, q.items[q.offset:])
		if q.offset > 0 {
			copy(data[q.size-q.offset:], q.items[:q.offset])
		}
	}
	return data
}

func (q *CircularQueue) Put(item interface{}) {
	q.items[q.offset] = item
	if q.len < q.size {
		q.len++
	}
	q.offset = (q.offset + 1) % q.size
}

func (q *CircularQueue) SafePut(item interface{}) {
	var offset, newOffset uint32
	for {
		offset = q.offset
		newOffset = (offset + 1) % q.size
		if atomic.CompareAndSwapUint32(&q.offset, offset, newOffset) {
			if q.len < q.size {
				atomic.AddUint32(&q.len, 1)
			}
			q.items[offset] = item
			break
		} else {
			runtime.Gosched()
		}
	}
}
