//go:build !noasm && !appengine && !gccgo

package crc64nvme

import (
	. "github.com/klauspost/cpuid/v2"
)

var hasAsm = CPU.Supports(ASIMD) && CPU.Supports(PMULL)

func updateAsm(crc uint64, p []byte) (checksum uint64)
func updateAsmSingle()
