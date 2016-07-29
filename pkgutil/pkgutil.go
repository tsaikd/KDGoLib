package pkgutil

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorImpossibleImportDir1 = errutil.NewFactory("guess the import package info failed from %q")
)

// GuessPackageFromDir return importable package from dir if possible
func GuessPackageFromDir(dir string) (pkg *build.Package, err error) {
	absdir, err := filepath.Abs(dir)
	if err != nil {
		return
	}

	dirsplit := strings.Split(absdir, string(os.PathSeparator))
	for i, n := 1, len(dirsplit); i < n; i++ {
		trypath := strings.Join(dirsplit[i:], string(os.PathSeparator))
		if pkg, err = build.Import(trypath, absdir, 0); err == nil {
			return
		}
	}

	return nil, ErrorImpossibleImportDir1.New(nil, absdir)
}

// FindAllSubPackages return all sub packages under importPath
func FindAllSubPackages(importPath string, srcDir string) (pkglist *PackageList, err error) {
	pkglist = &PackageList{}
	if srcDir, err = filepath.Abs(srcDir); err != nil {
		return nil, err
	}
	if err = collectPackage(pkglist, importPath, srcDir); err != nil {
		return nil, err
	}
	return pkglist, nil
}

func collectPackage(result *PackageList, importPath string, srcDir string) (err error) {
	pkg, err := build.Import(importPath, srcDir, 0)
	if err == nil {
		result.AddPackage(pkg)
		if len(pkg.Dir) > len(srcDir) && strings.HasPrefix(pkg.Dir, srcDir) {
			if err = collectPackageFromSubDir(result, importPath, pkg.Dir); err != nil {
				return
			}
		}
	}

	return collectPackageFromSubDir(result, importPath, srcDir)
}

func collectPackageFromSubDir(result *PackageList, importPath string, srcDir string) (err error) {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		name := file.Name()
		if name == "vendor" {
			continue
		}
		if strings.HasPrefix(name, ".") {
			continue
		}

		childPath := filepath.Join(importPath, name)
		childSrcDir := filepath.Join(srcDir, name)
		if err = collectPackage(result, childPath, childSrcDir); err != nil {
			return
		}
	}

	return nil
}
