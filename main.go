package main

import (
	"fmt"
	"unsafe"
)

//go:noescape
func cpuidex(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32)

func DetectHyperV() bool {
	ret := false

	// Get Vendor String
	buf := make([]byte, 13)
	_, ebx, ecx, edx := cpuidex(0x40000000, 0)

	// Copy vendor string to buf
	*(*uint32)(unsafe.Pointer(&buf[0])) = ebx
	*(*uint32)(unsafe.Pointer(&buf[4])) = ecx
	*(*uint32)(unsafe.Pointer(&buf[8])) = edx
	buf[12] = 0

	// Check if it's "Microsoft Hv"
	if *(*uint32)(unsafe.Pointer(&buf[0])) == 0x7263694D &&
		*(*uint32)(unsafe.Pointer(&buf[4])) == 0x666F736F &&
		*(*uint32)(unsafe.Pointer(&buf[8])) == 0x76482074 {

		// Get feature flags
		eax, ebx, _, _ := cpuidex(0x40000003, 0)

		// Check specific bits
		if ((eax & 0x00001000) == 0) &&
			((ebx & 0x00000800) == 0) &&
			((ebx & 0x00000001) == 0) &&
			((ebx & 0x00000040) == 0) {
			ret = true
		}
	}

	return ret
}

func main() {
	isHyperV := DetectHyperV()
	if isHyperV {
		fmt.Println("[+] Running in Hyper-V(Guest) environment")
	} else {
		fmt.Println("[-] Not running in Hyper-V(Guest) environment")
	}
}
