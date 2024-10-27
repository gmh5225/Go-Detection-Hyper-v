#include "textflag.h"

// func cpuidex(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32)
TEXT ·cpuidex(SB), NOSPLIT, $0-24
    MOVL leaf+0(FP), AX
    MOVL subleaf+8(FP), CX
    CPUID
    MOVL AX, eax+16(FP)
    MOVL BX, ebx+20(FP)
    MOVL CX, ecx+24(FP)
    MOVL DX, edx+28(FP)
    RET
