package goinis

type KeyValue struct {
	lock sync.RWMutex
	K    string
	V    interface{}
}

func NewKeyValue(key string, value interface{}, comment ...string) *KeyValue {
	kv := &KeyValue{
		K: key,
	}

	if Util.IsArrayKey(key) {
		kv.V = append(make([]string, 0), "string")
	} else {
		kv.V = value
	}

	return kv
}

func (kv *KeyValue) SetValue(value string) *KeyValue {
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

func (kv *KeyValue) GetValue() interface{} {
	return kv.V
}

func (kv *KeyValue) AddValue(str string) *KeyValue {
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
