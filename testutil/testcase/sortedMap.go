package testcase

import "sort"

// SortedMap is a map and can be sorted
type SortedMap struct {
	datamap map[string]interface{}
	datakey []string
}

// Get return value by key
func (t *SortedMap) Get(name string) interface{} {
	if t.datamap == nil {
		panic("provider " + name + " not found")
	}
	return t.datamap[name]
}

// Set value to map by key
func (t *SortedMap) Set(name string, provider interface{}) {
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
func (t *SortedMap) Remove(name string) {
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
func (t *SortedMap) IsExists(name string) bool {
	if t.datamap == nil {
		return false
	}
	_, exist := t.datamap[name]
	return exist
}

// Last return last element, return nil if len == 0
func (t *SortedMap) Last() interface{} {
	if t.Len() < 1 {
		return nil
	}

	key := t.datakey[t.Len()-1]
	return t.datamap[key]
}

// Shift an element, return nil if len == 0
func (t *SortedMap) Shift() interface{} {
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
func (t *SortedMap) Walk(callback func(name string, element interface{}) (stop bool)) {
	for _, key := range t.datakey {
		element := t.datamap[key]
		if callback(key, element) {
			return
		}
	}
}

// Sort by keys
func (t *SortedMap) Sort() {
	sort.Sort(t)
}

// Keys return all keys
func (t *SortedMap) Keys() []string {
	return t.datakey
}

// Len length of map, used for sort
func (t *SortedMap) Len() int {
	return len(t.datamap)
}

// Less compare key, used for sort
func (t *SortedMap) Less(i, j int) bool {
	return t.datakey[i] < t.datakey[j]
}

// Swap key position, used for sort
func (t *SortedMap) Swap(i, j int) {
	t.datakey[i], t.datakey[j] = t.datakey[j], t.datakey[i]
}
