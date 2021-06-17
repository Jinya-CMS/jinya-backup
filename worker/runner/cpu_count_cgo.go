// +build cgo

package runner

import "runtime"

var CpuCount = runtime.NumCPU()
