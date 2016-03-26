package apimgr

import "reflect"

func newSorter(manager *Manager) *sorter {
	apis := []Definition{}
	for _, api := range manager.apiMethodPatternMap {
		apis = append(apis, api)
	}
	return &sorter{
		Manager: manager,
		apis:    apis,
	}
}

type sorter struct {
	*Manager

	apis []Definition
}

func (t sorter) Len() int {
	return len(t.apis)
}

func (t sorter) Swap(i int, j int) {
	t.apis[i], t.apis[j] = t.apis[j], t.apis[i]
}

func (t sorter) Less(i int, j int) bool {
	ki := t.getSortKey(t.apis[i])
	kj := t.getSortKey(t.apis[j])
	return ki < kj
}

func (t sorter) getSortKey(api Definition) string {
	pkgpath := getPackagePath(reflect.ValueOf(api.Request))
	key := pkgpath + " " + t.GetMethodPatternKey(t.Manager, api)
	return key
}
