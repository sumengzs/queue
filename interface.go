/*
Copyright 2022 The KubePort Authors.
*/

package queue

type Interface interface {
	Len() int
	Cap() int
	Get() interface{}
	Put(item interface{})
}
