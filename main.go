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

// Microsoft Hv signature in hex
const (
	MS_HV_SIGNATURE_1 = 0x7263694D // "Micr"
	MS_HV_SIGNATURE_2 = 0x666F736F // "osof"
	MS_HV_SIGNATURE_3 = 0x76482074 // "t Hv"
)

// CPUID leaf values
const (
	CPUID_HV_VENDOR_AND_MAX_FUNCTIONS = 0x40000000
	CPUID_HV_INTERFACE                = 0x40000003
)

// isHyperVSignature checks if the buffer contains Hyper-V signature
func isHyperVSignature(buf []byte) bool {
	return *(*uint32)(unsafe.Pointer(&buf[0])) == MS_HV_SIGNATURE_1 &&
		*(*uint32)(unsafe.Pointer(&buf[4])) == MS_HV_SIGNATURE_2 &&
		*(*uint32)(unsafe.Pointer(&buf[8])) == MS_HV_SIGNATURE_3
}

// hasRequiredFeatures checks if the required features are present
func hasRequiredFeatures(ebx uint32) bool {
	return (ebx&HV_DEBUGGING) == 0 &&
		(ebx&HV_CREATE_PARTITIONS) == 0 &&
		(ebx&HV_CREATE_PORT) == 0
}

func DetectHyperV() bool {
	// Get Vendor String
	buf := make([]byte, 13)
	_, ebx, ecx, edx := cpuidex(CPUID_HV_VENDOR_AND_MAX_FUNCTIONS, 0)

	// Copy vendor string to buffer
	*(*uint32)(unsafe.Pointer(&buf[0])) = ebx
	*(*uint32)(unsafe.Pointer(&buf[4])) = ecx
	*(*uint32)(unsafe.Pointer(&buf[8])) = edx
	buf[12] = 0

	// Return false if not Hyper-V signature
	if !isHyperVSignature(buf) {
		return false
	}

	// Get feature flags
	_, ebx, _, _ = cpuidex(CPUID_HV_INTERFACE, 0)

	// Check required features
	return hasRequiredFeatures(ebx)
}

func main() {
	if isHyperVGuest := DetectHyperV(); isHyperVGuest {
		fmt.Println("[+] Running in Hyper-V(Guest) environment")
	} else {
		fmt.Println("[-] Not running in Hyper-V(Guest) environment")
	}
}
