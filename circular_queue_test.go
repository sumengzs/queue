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
	"sort"
	"sync"
	"testing"
)

func prepareCircularQueue(size uint32, length uint32) *CircularQueue {
	queue := NewCircularQueue(size)
	for i := 0; i < int(length); i++ {
		queue.Put(i)
	}
	return queue
}

// prepareWantData
// exp:
// [0:5) [0,1,2,3,4]
// (5:0] [4,3,2,1,0]
func prepareWantData(start, end int) []interface{} {
	data := make([]interface{}, 0)
	if start <= end {
		for i := start; i < end; i++ {
			data = append(data, i)
		}
	} else {
		for i := start - 1; i >= end; i-- {
			data = append(data, i)
		}
	}
	return data
}

func compareWithInt32(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].(int) != b[i].(int) {
			return false
		}
	}
	return true
}

func TestCircularQueue_Get(t *testing.T) {
	tests := []struct {
		name  string
		queue *CircularQueue
		want  interface{}
	}{
		{
			name:  "queue nodata",
			queue: prepareCircularQueue(5, 0),
			want:  nil,
		},
		{
			name:  "queue size equal to data length",
			queue: prepareCircularQueue(5, 5),
			want:  4,
		},
		{
			name:  "queue size less than data length",
			queue: prepareCircularQueue(3, 5),
			want:  4,
		},
		{
			name:  "queue size more than data length",
			queue: prepareCircularQueue(5, 3),
			want:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.queue
			if got := q.Get(); got != tt.want && got.(int) != tt.want.(int) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_GetAll(t *testing.T) {
	tests := []struct {
		name  string
		queue *CircularQueue
		want  []interface{}
	}{
		{
			name:  "queue nodata",
			queue: prepareCircularQueue(5, 0),
			want:  prepareWantData(0, 0),
		},
		{
			name:  "queue size equal to data length",
			queue: prepareCircularQueue(5, 5),
			want:  prepareWantData(0, 5),
		},
		{
			name:  "queue size less than data length",
			queue: prepareCircularQueue(3, 5),
			want:  prepareWantData(2, 5),
		},
		{
			name:  "queue size more than data length",
			queue: prepareCircularQueue(5, 3),
			want:  prepareWantData(0, 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.queue
			if got := q.GetAll(); !compareWithInt32(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_GetPoint(t *testing.T) {
	tests := []struct {
		name  string
		queue *CircularQueue
		args  uint32
		want  interface{}
	}{
		{
			name:  "queue nodata",
			queue: prepareCircularQueue(5, 0),
			args:  3,
			want:  nil,
		},
		{
			name:  "queue size equal to data length",
			queue: prepareCircularQueue(5, 5),
			args:  3,
			want:  1,
		},
		{
			name:  "queue size less than data length",
			queue: prepareCircularQueue(3, 5),
			args:  3,
			want:  nil,
		},
		{
			name:  "queue size more than data length",
			queue: prepareCircularQueue(5, 3),
			args:  3,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.queue
			if got := q.GetPoint(tt.args); got != tt.want && got.(int) != tt.want.(int) {
				t.Errorf("GetPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_Gets(t *testing.T) {
	tests := []struct {
		name  string
		queue *CircularQueue
		args  uint32
		want  []interface{}
	}{
		{
			name:  "queue nodata",
			queue: prepareCircularQueue(5, 0),
			args:  3,
			want:  nil,
		},
		{
			name:  "queue size equal to data length",
			queue: prepareCircularQueue(5, 5),
			args:  3,
			want:  prepareWantData(2, 5),
		},
		{
			name:  "queue size less than data length",
			queue: prepareCircularQueue(3, 5),
			args:  3,
			want:  prepareWantData(2, 5),
		},
		{
			name:  "queue size less than data length, get size more than data length",
			queue: prepareCircularQueue(3, 5),
			args:  4,
			want:  prepareWantData(2, 5),
		},
		{
			name:  "queue size more than data length",
			queue: prepareCircularQueue(5, 3),
			args:  3,
			want:  prepareWantData(0, 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.queue
			if got := q.Gets(tt.args); !compareWithInt32(got, tt.want) {
				t.Errorf("Gets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_Len(t *testing.T) {
	tests := []struct {
		name  string
		queue *CircularQueue
		want  uint32
	}{
		{
			name:  "queue nodata",
			queue: prepareCircularQueue(5, 0),
			want:  0,
		},
		{
			name:  "queue size equal to data length",
			queue: prepareCircularQueue(5, 5),
			want:  5,
		},
		{
			name:  "queue size less than data length",
			queue: prepareCircularQueue(3, 5),
			want:  3,
		},
		{
			name:  "queue size more than data length",
			queue: prepareCircularQueue(5, 3),
			want:  3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.queue
			if got := q.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTestPut(b *testing.B) {
	b.ResetTimer()
	queue := NewCircularQueue(1024)
	for i := 0; i < b.N; i++ {
		queue.Put(i)
	}
}

func BenchmarkTestSafePut(b *testing.B) {
	b.ResetTimer()
	queue := NewCircularQueue(1024)
	for i := 0; i < b.N; i++ {
		queue.SafePut(i)
	}
}

func BenchmarkTestGet(b *testing.B) {
	queue := prepareCircularQueue(1024, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Get()
	}
}

func BenchmarkTestGetAll(b *testing.B) {
	queue := prepareCircularQueue(1024, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.GetAll()
	}
}

func BenchmarkTestGets(b *testing.B) {
	b.ReportAllocs()
	queue := prepareCircularQueue(1024, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Gets(uint32(i) % 1024)
	}
}

func TestCircularQueue_SafePutConcurrent(t *testing.T) {
	numGoroutines := 1000
	q := NewCircularQueue(uint32(numGoroutines * numGoroutines))

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()

			for j := i * numGoroutines; j < i*numGoroutines+numGoroutines; j++ {
				q.SafePut(j)
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	actLen := q.Len()
	expLen := numGoroutines * numGoroutines
	if actLen != uint32(expLen) {
		t.Errorf("actLen=%d expLen=%d", actLen, expLen)
	}
	data := q.GetAll()
	sort.Slice(data, func(i, j int) bool {
		return data[i].(int) < data[j].(int)
	})
	for i := range data {
		if i != data[i].(int) {
			t.Errorf("want=%d act:=%d", i, data[i])
		}
	}
}
