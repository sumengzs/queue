/*
Copyright 2022 The KubePort Authors.
*/

package queue

import (
	"testing"
)

func prepareCircularQueue(size int, length int) *CircularQueue {
	queue := NewCircularQueue(size)
	for i := 0; i < length; i++ {
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
	if start < end {
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
		args  int
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
			want:  2,
		},
		{
			name:  "queue size less than data length",
			queue: prepareCircularQueue(3, 5),
			args:  3,
			want:  2,
		},
		{
			name:  "queue size more than data length",
			queue: prepareCircularQueue(5, 3),
			args:  3,
			want:  0,
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
		args  int
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
		want  int
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
	b.ReportAllocs()
	b.ResetTimer()
	queue := NewCircularQueue(1024)
	for i := 0; i < b.N; i++ {
		queue.Put(i)
	}
}
func BenchmarkTestGet(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	queue := prepareCircularQueue(1024, 1024)
	for i := 0; i < b.N; i++ {
		queue.Get()
	}
}

func BenchmarkTestGetAll(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	queue := prepareCircularQueue(1024, 1024)
	for i := 0; i < b.N; i++ {
		queue.GetAll()
	}
}

func BenchmarkTestGets(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	queue := prepareCircularQueue(1024, 1024)
	for i := 0; i < b.N; i++ {
		queue.Gets(i % 1024)
	}
}
