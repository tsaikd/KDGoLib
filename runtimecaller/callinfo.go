package runtimecaller

import "runtime"

// CallInfo contains runtime caller information
type CallInfo struct {
	// builtin data
	PC       uintptr
	FilePath string
	Line     int

	// extra info after some process
	PCFunc      *runtime.Func
	PackageName string
	FileDir     string
	FileName    string
	FuncName    string
}
