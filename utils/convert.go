package utils

import (
	"reflect"
	"unsafe"
)

// S2B StringToBytes
func S2B(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}

// B2S BytesToString
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
