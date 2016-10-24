package pkgutil

import (
	"go/build"
	"regexp"
	"sort"
)

// PackageList contains Package info and support lookup efficiently
type PackageList struct {
	pkgpool  map[*build.Package]bool
	dirpool  map[string]*build.Package
	namepool map[string]*build.Package
}

func (t *PackageList) ensureInit() {
	if t.pkgpool == nil {
		t.pkgpool = map[*build.Package]bool{}
	}
	if t.dirpool == nil {
		t.dirpool = map[string]*build.Package{}
	}
	if t.namepool == nil {
		t.namepool = map[string]*build.Package{}
	}
}

// AddPackage add pkg into PackageList
func (t *PackageList) AddPackage(pkg *build.Package) {
	t.ensureInit()
	t.pkgpool[pkg] = true
	t.dirpool[pkg.Dir] = pkg
	t.namepool[getPackageName(pkg)] = pkg
}

// RemovePackage remove pkg from PackageList
func (t *PackageList) RemovePackage(pkg *build.Package) {
	if pkg == nil {
		return
	}
	t.RemoveByDir(pkg.Dir)
}

// RemoveByDir remove pkg from PackageList by dir
func (t *PackageList) RemoveByDir(dir string) {
	t.ensureInit()
	if pkg := t.LookupByDir(dir); pkg != nil {
		delete(t.pkgpool, pkg)
		delete(t.dirpool, pkg.Dir)
		delete(t.namepool, getPackageName(pkg))
	}
}

// RemoveByName remove pkg from PackageList by import path
func (t *PackageList) RemoveByName(name string) {
	t.ensureInit()
	if pkg := t.LookupByName(name); pkg != nil {
		delete(t.pkgpool, pkg)
		delete(t.dirpool, pkg.Dir)
		delete(t.namepool, getPackageName(pkg))
	}
}

// LookupByDir lookup Package by dir
func (t *PackageList) LookupByDir(dir string) *build.Package {
	t.ensureInit()
	return t.dirpool[dir]
}

// LookupByName lookup Package by import path
func (t *PackageList) LookupByName(name string) *build.Package {
	t.ensureInit()
	return t.namepool[name]
}

// Map return packages in map form
func (t *PackageList) Map() map[*build.Package]bool {
	t.ensureInit()
	return t.pkgpool
}

// Sorted return packages in slice form and sorted by package directory
func (t *PackageList) Sorted() []*build.Package {
	t.ensureInit()

	keys := make([]string, len(t.dirpool))
	i := 0
	for k := range t.dirpool {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	result := make([]*build.Package, len(t.dirpool))
	i = 0
	for _, k := range keys {
		result[i] = t.dirpool[k]
		i++
	}

	return result
}

// Len return length of PackageList
func (t *PackageList) Len() int {
	t.ensureInit()
	return len(t.pkgpool)
}

func getPackageName(pkg *build.Package) string {
	regVendor := regexp.MustCompile(`^.*/vendor/`)
	return regVendor.ReplaceAllString(pkg.ImportPath, "")
}
