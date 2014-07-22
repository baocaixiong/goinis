package goinis

import (
	"sync"
)

type keyValue struct {
	lock sync.RWMutex
	K    string
	V    interface{}
}

func newKeyValue(key string, value interface{}) *keyValue {
	kv := &keyValue{
		K: key,
	}

	if Util.IsArrayKey(key) {
		kv.V = append(make([]string, 0), value.(string))
	} else {
		kv.V = value.(string)
	}

	return kv
}

func (kv *keyValue) setValue(value string) *keyValue {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	switch v := kv.V.(type) {
	case []string:
		kv.V = append(v, value)
	case string:
		kv.V = value
	}

	return kv
}

func (kv *keyValue) getValue() interface{} {
	return kv.V
}

func (kv *keyValue) addValue(str string) *keyValue {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	switch v := kv.V.(type) {
	case string:
		kv.V = v + str
	case []string:
		last := v[len(v)-1]
		kv.V = append(v[:len(v)-1], last+str)
	}

	return kv
}
