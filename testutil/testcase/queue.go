package testcase

// Queue is a data struct support map and slice feature
type Queue interface {
	Get(name string) interface{}
	Set(name string, provider interface{})
	Remove(name string)
	IsExists(name string) bool
	Last() interface{}
	Shift() interface{}
	Walk(callback func(name string, element interface{}) (stop bool))
	Keys() []string
	Len() int
}

// QueueType implement Queue
type QueueType struct {
	datamap map[string]interface{}
	datakey []string
}

// Get return value by key
func (t *QueueType) Get(name string) interface{} {
	if t.datamap == nil {
		panic("provider " + name + " not found")
	}
	return t.datamap[name]
}

// Set value to map by key
func (t *QueueType) Set(name string, provider interface{}) {
	if t.datamap == nil {
		t.datamap = map[string]interface{}{}
		t.datakey = []string{}
	}
	if _, exist := t.datamap[name]; exist {
		panic("provider " + name + " already exist")
	}
	t.datamap[name] = provider
	t.datakey = append(t.datakey, name)
}

// Remove key from map
func (t *QueueType) Remove(name string) {
	if t.datamap == nil {
		return
	}
	if _, exist := t.datamap[name]; !exist {
		panic("provider " + name + " not found in map")
	}
	delete(t.datamap, name)

	idx := -1
	for i, key := range t.datakey {
		if key == name {
			idx = i
			break
		}
	}
	if idx < 0 {
		panic("provider " + name + " not found in slice")
	}
	t.datakey = append(t.datakey[:idx], t.datakey[idx+1:]...)
}

// IsExists return true if name exists in map
func (t *QueueType) IsExists(name string) bool {
	if t.datamap == nil {
		return false
	}
	_, exist := t.datamap[name]
	return exist
}

// Last return last element, return nil if len == 0
func (t *QueueType) Last() interface{} {
	if t.Len() < 1 {
		return nil
	}

	key := t.datakey[t.Len()-1]
	return t.datamap[key]
}

// Shift an element, return nil if len == 0
func (t *QueueType) Shift() interface{} {
	if t.Len() < 1 {
		return nil
	}

	var firstKey string
	firstKey, t.datakey = t.datakey[0], t.datakey[1:]
	element := t.datamap[firstKey]
	delete(t.datamap, firstKey)
	return element
}

// Walk all data for callback
func (t *QueueType) Walk(callback func(name string, element interface{}) (stop bool)) {
	for _, key := range t.datakey {
		element := t.datamap[key]
		if callback(key, element) {
			return
		}
	}
}

// Keys return all keys
func (t *QueueType) Keys() []string {
	return t.datakey
}

// Len length of map
func (t *QueueType) Len() int {
	return len(t.datamap)
}
