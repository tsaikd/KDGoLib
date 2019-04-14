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
		return
	}
	if importPath == "" {
		var pkg *build.Package
		if pkg, err = GuessPackageFromDir(srcDir); err != nil {
			if !ErrorImpossibleImportDir1.Match(err) {
				return
			}
			if IsGoModDir(srcDir) {
				var mods []Module
				mods, err = ParseGoMod(srcDir)
				if err != nil {
					return
				}
				err = collectModule(pkglist, mods)
			}
			return nil, ErrorImpossibleImportDir1.New(nil, srcDir)
		}
		importPath = pkg.ImportPath
	}
	if err = collectPackage(pkglist, importPath, srcDir); err != nil {
		return
	}
	return
}

// ParsePackagePaths return PackageList by parse paths, ignore vendor
func ParsePackagePaths(srcDir string, paths ...string) (pkglist *PackageList, err error) {
	var pkg *build.Package
	pkglist = &PackageList{}
	if srcDir, err = filepath.Abs(srcDir); err != nil {
		return
	}

	if len(paths) < 1 {
		if pkg, err = GuessPackageFromDir(srcDir); err != nil {
			return
		}
		pkglist.AddPackage(pkg)
		return
	}

	for _, pkgpath := range paths {
		importPath := strings.TrimSuffix(pkgpath, "/...")
		if importPath != pkgpath {
			if err = collectPackage(pkglist, importPath, srcDir); err != nil {
				return
			}
		} else {
			if pkg, err = build.Import(importPath, srcDir, 0); err != nil {
				return
			}
			pkglist.AddPackage(pkg)
		}
	}

	return
}

func collectPackage(result *PackageList, importPath string, srcDir string) (err error) {
	pkg, err := build.Import(importPath, srcDir, 0)
	if err == nil {
		result.AddPackage(pkg)
		if len(pkg.Dir) > len(srcDir) && strings.HasPrefix(pkg.Dir, srcDir) {
			if err = collectPackageFromSubDir(result, pkg.ImportPath, pkg.Dir); err != nil {
				return
			}
		}
		importPath = pkg.ImportPath
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

func collectModule(result *PackageList, modules []Module) (err error) {
	for _, mod := range modules {
		result.AddPackage(&build.Package{
			Dir:        mod.Dir,
			Name:       mod.GoMod,
			ImportPath: mod.Path,
		})
	}
	return
}
