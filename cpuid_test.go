//go:build test
// +build test

package main

import (
	_ "unsafe" // required for go:linkname
)

// MockCPUID is the mock function for testing
var MockCPUID func(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32)

//go:linkname cpuidex cpuidex
func cpuidex(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32) {
	if MockCPUID != nil {
		return MockCPUID(leaf, subleaf)
	}
	return 0, 0, 0, 0
}
