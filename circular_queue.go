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
	"sync"
	"sync/atomic"
)

type CircularQueue struct {
	mu     sync.RWMutex
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

func (q *CircularQueue) Len() int {
	return int(atomic.LoadUint32(&q.len))
}

func (q *CircularQueue) Cap() int {
	return int(q.size)
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
	return data[(len(data))-int(size):]
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

// AtomicGet 原子读取
// 需要注意的是，AtomicGet 方法只保证了对偏移量的原子读取，而不保证对元素本身的读取过程中的并发安全性。
// 如果在同一时间段内有其他协程在往队列中放入或取出元素，那么使用原子操作读取元素时仍然可能会遇到并发访问的问题。
// 如果需要保证绝对的准确性，建议配合使用 LockGet 与 LockPut。
func (q *CircularQueue) AtomicGet() interface{} {
	offset := atomic.LoadUint32(&q.offset)
	if q.len == 0 {
		return nil
	}
	if offset == 0 {
		return q.items[q.size-1]
	}
	return q.items[offset-1]
}

// AtomicPut 使用 atomic.CompareAndSwapUint32 保证了对偏移量的原子读取与操作，继而保证该偏移位置只能用此次操作赋值。
// 经过测试，是协程安全的。
func (q *CircularQueue) AtomicPut(item interface{}) {
	offset := atomic.LoadUint32(&q.offset)
	newOffset := (offset + 1) % q.size
	for {
		if atomic.CompareAndSwapUint32(&q.offset, offset, newOffset) {
			if q.len < q.size {
				atomic.AddUint32(&q.len, 1)
			}
			q.items[offset] = item
			break
		} else {
			offset = atomic.LoadUint32(&q.offset)
			newOffset = (offset + 1) % q.size
			runtime.Gosched()
		}
	}
}

func (q *CircularQueue) MutexPut(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Put(item)
}

func (q *CircularQueue) MutexGet() interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.Get()
}
