package pantsu

import (
	"reflect"
	"unsafe"
)

//his is an unsafe way, the result string and []byte buffer share the same bytes.
//Please make sure not to modify the bytes in the []byte buffer if the string still survives!
func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// func b2s(b []byte) string {
// 	return *(*string)(unsafe.Pointer(&b))
// }

const splitter = ':'

func findPathIndex(s1 string) int {
	b := s2b(s1)
	return bFindPathIndex(b)
}

func bFindPathIndex(b1 []byte) int {
	for idx, v := range b1 {
		if v == splitter {
			return idx
		}
	}
	return -1
}
