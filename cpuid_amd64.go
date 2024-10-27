//go:build !test
// +build !test

package main

//go:noescape
func cpuidex(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32)
