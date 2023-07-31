//go:build cgo

package helper

import "runtime"

var CpuCount = runtime.NumCPU()
