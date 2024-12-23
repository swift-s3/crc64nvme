//go:build !noasm && !appengine && !gccgo

package crc64nvme

var hasAsm = false

func updateAsm(crc uint64, p []byte) (checksum uint64)
