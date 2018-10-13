// +build go1.9

#include "textflag.h"

TEXT ·Asmxoro(SB),NOSPLIT,$0
	MOVQ	xoro+0(FP), AX	// AX = xoro
	MOVQ	(AX), CX		// CX = s0 = (*xoro)[0]
	MOVQ	8(AX), DX		// DX = s1 = (*xoro)[1]
	MOVQ	CX, BX			// BX = s0
	XORQ	DX, CX			// CX = s1 ^= s0
	MOVQ	BX, SI			// SI = s0
	ROLQ	$55, BX			// BX = rotl(s0, 55)
	XORQ	CX, BX			// BX = rotl(s0, 55) ^ s1
	MOVQ	CX, DI			// DI = s1
	SHLQ	$14, CX			// CX = s1 << 14
	XORQ	BX, CX          // CX = rotl(s0, 55) ^ s1 ^ s1<<14
	MOVQ	CX, (AX)		// (*xoro)[0] = rotl(s0, 55) ^ s1 ^ s1<<14
	ROLQ	$36, DI			// DI = rotl(s1, 36)
	MOVQ	DI, 8(AX)		// (*xoro)[1] = rotl(s1, 36)
	LEAQ	(SI)(DX*1), AX	// DX = x = s0 + (old)s1
	MOVQ	AX, ret+8(FP)	// return x
	RET
