package main

import (
	"fmt"
	"unsafe"
)

// Hyper-V Feature Bits (ebx)
const (
	HV_CREATE_PARTITIONS     = 0x00000001
	HV_ACCESS_PARTITION_ID   = 0x00000002
	HV_ACCESS_MEMORY_POOL    = 0x00000004
	HV_ADJUST_MEMORY_BUFFERS = 0x00000008
	HV_POST_MESSAGES         = 0x00000010
	HV_SIGNAL_EVENTS         = 0x00000020
	HV_CREATE_PORT           = 0x00000040
	HV_CONNECT_PORT          = 0x00000080
	HV_ACCESS_STATS          = 0x00000100
	// 0x00000200 Reserved
	// 0x00000400 Reserved
	HV_DEBUGGING                    = 0x00000800
	HV_CPU_MANAGEMENT               = 0x00001000
	HV_CONFIGURE_PROFILER           = 0x00002000
	HV_ENABLE_EXPANDED_STACKWALKING = 0x00004000
	// 0x00008000 Reserved
	HV_ACCESS_VSM          = 0x00010000
	HV_ACCESS_VP_REGISTERS = 0x00020000
	// 0x00040000 Reserved
	// 0x00080000 Reserved
	HV_ENABLE_EXTENDED_HYPERCALLS = 0x00100000
	HV_START_VIRTUAL_PROCESSOR    = 0x00200000
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
		_, ebx, _, _ := cpuidex(0x40000003, 0)

		// Check specific bits
		if ((ebx & HV_DEBUGGING) == 0) &&
			((ebx & HV_CREATE_PARTITIONS) == 0) &&
			((ebx & HV_CREATE_PORT) == 0) {
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
