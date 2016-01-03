package apimgr

import "reflect"

func getPackagePath(pkg reflect.Value) (pkgpath string) {
	if pkg.Kind() == reflect.String {
		return pkg.String()
	}
	pkgpath = pkg.Type().PkgPath()
	if pkgpath != "" {
		return
	}
	if pkg.Kind() == reflect.Ptr {
		return getPackagePath(pkg.Elem())
	}
	return ""
}
