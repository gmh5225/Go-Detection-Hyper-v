package main

import (
	"testing"
)

// Mock CPUID data structure
type CPUIDData struct {
	leaf    uint32
	subleaf uint32
	eax     uint32
	ebx     uint32
	ecx     uint32
	edx     uint32
}

// Test cases
var testCases = []struct {
	name     string
	cpuidSeq []CPUIDData
	want     bool
}{
	{
		name: "VMWare[win with HV, chk.exe]",
		cpuidSeq: []CPUIDData{
			{
				leaf:    0x40000000,
				subleaf: 0,
				eax:     0x4000000C,
				ebx:     0x7263694D, // "Micr"
				ecx:     0x666F736F, // "osof"
				edx:     0x76482074, // "t Hv"
			},
			{
				leaf:    0x40000003,
				subleaf: 0,
				eax:     0x00003FFF,
				ebx:     0x002BB9FF,
				ecx:     0x00000002,
				edx:     0x71FEFFF2,
			},
		},
		want: false,
	},
	{
		name: "VMWare[win with HV + wsbox[chk.exe]]",
		cpuidSeq: []CPUIDData{
			{
				leaf:    0x40000000,
				subleaf: 0,
				eax:     0x4000000B,
				ebx:     0x7263694D, // "Micr"
				ecx:     0x666F736F, // "osof"
				edx:     0x76482074, // "t Hv"
			},
			{
				leaf:    0x40000003,
				subleaf: 0,
				eax:     0x00006E7F,
				ebx:     0x003B8030,
				ecx:     0x00000002,
				edx:     0xE0BED7B2,
			},
		},
		want: true,
	},
}

func setupMockCPUID(t *testing.T, cpuidSeq []CPUIDData) {
	var callCount int
	MockCPUID = func(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32) {
		if callCount >= len(cpuidSeq) {
			t.Fatalf("Too many CPUID calls")
		}
		data := cpuidSeq[callCount]
		callCount++

		// Verify leaf and subleaf match expected values
		if leaf != data.leaf || subleaf != data.subleaf {
			t.Errorf("Unexpected CPUID leaf/subleaf: got (%x,%x), want (%x,%x)",
				leaf, subleaf, data.leaf, data.subleaf)
		}

		return data.eax, data.ebx, data.ecx, data.edx
	}
}

func TestDetectHyperV(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setupMockCPUID(t, tc.cpuidSeq)
			got := DetectHyperV()
			if got != tc.want {
				t.Errorf("DetectHyperV() = %v, want %v", got, tc.want)
			}
		})
	}
}

func FuzzDetectHyperV(f *testing.F) {
	// Add the test cases as seeds
	for _, tc := range testCases {
		f.Add(tc.name)
	}

	f.Fuzz(func(t *testing.T, testName string) {
		// Find the test case
		var tc *struct {
			name     string
			cpuidSeq []CPUIDData
			want     bool
		}
		for _, testCase := range testCases {
			if testCase.name == testName {
				tc = &testCase
				break
			}
		}
		if tc == nil {
			t.Skip("Unknown test case")
		}

		setupMockCPUID(t, tc.cpuidSeq)
		_ = DetectHyperV()
	})
}
