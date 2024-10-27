#include "textflag.h"

// func cpuidex(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32)
TEXT ·cpuidex(SB), NOSPLIT, $0-24
    MOVL    leaf+0(FP), AX      
    MOVL    subleaf+4(FP), CX   
    CPUID                       
    MOVL    AX, ret+8(FP)      
    MOVL    BX, ret+12(FP)     
    MOVL    CX, ret+16(FP)     
    MOVL    DX, ret+20(FP)     
    RET

