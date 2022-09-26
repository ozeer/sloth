package tool

import "hash/crc32"

func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
