package orderedmap

import "encoding/json"

type OrderedMap interface {
	Add(v interface{}) bool
	Remove(v interface{}) bool
	IsEmpty() bool
	Exist(v interface{}) bool
	Map() map[interface{}]bool
	Slice() []interface{}
	MarshalJSON() ([]byte, error)
}

type OrderedMapData struct {
	datamap   map[interface{}]bool
	dataslice []interface{}
}

func New(vs ...interface{}) (t OrderedMap) {
	t = &OrderedMapData{
		datamap:   make(map[interface{}]bool),
		dataslice: []interface{}{},
	}
	for _, v := range vs {
		t.Add(v)
	}
	return
}

func (t *OrderedMapData) Add(v interface{}) bool {
	_, ok := t.datamap[v]
	if ok {
		return false
	} else {
		t.datamap[v] = true
		t.dataslice = append(t.dataslice, v)
		return true
	}
}

func (t *OrderedMapData) Remove(v interface{}) bool {
	_, ok := t.datamap[v]
	if !ok {
		return false
	}
	delete(t.datamap, v)
	for i, vinslice := range t.dataslice {
		if v == vinslice {
			t.dataslice = append(t.dataslice[:i], t.dataslice[i+1:]...)
			break
		}
	}
	return true
}

func (t OrderedMapData) IsEmpty() bool {
	return len(t.dataslice) < 1
}

func (t OrderedMapData) Exist(v interface{}) bool {
	_, exist := t.datamap[v]
	return exist
}

func (t OrderedMapData) Map() map[interface{}]bool {
	return t.datamap
}

func (t OrderedMapData) Slice() []interface{} {
	return t.dataslice
}

func (t OrderedMapData) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Slice())
}
