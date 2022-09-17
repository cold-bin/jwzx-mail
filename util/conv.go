// @author cold bin
// @date 2022/9/16

package util

import (
	"unsafe"
)

// QuickB2S []byte 快速转换出 string
func QuickB2S(bs []byte) string {
	bp := (*[3]uintptr)(unsafe.Pointer(&bs))
	sp := [3]uintptr{bp[0], bp[1], bp[1]}
	return *(*string)(unsafe.Pointer(&sp))
}

// QuickS2B string 快速转换出 []byte
func QuickS2B(s string) []byte {
	sp := (*[2]uintptr)(unsafe.Pointer(&s))
	bp := [2]uintptr{sp[0], sp[1]}
	return *(*[]byte)(unsafe.Pointer(&bp))
}
