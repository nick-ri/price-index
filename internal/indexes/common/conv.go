package common

import (
	"bytes"
	"reflect"
	"strconv"
	"unsafe"
)

func ZeroAllocStringToInt(val string, decimals int) int64 { // zero alloc translation
	var btes = []byte(val)

	foundPos := bytes.Index(btes, []byte{'.'})

	if foundPos != -1 {
		fraction := btes[foundPos+1:]
		btes = btes[:foundPos]

		if len(fraction) >= decimals {
			btes = append(btes, fraction[:decimals]...)
		} else { // rare case
			btes = append(btes, fraction...)

			for i := 0; i < decimals-len(fraction); i++ {
				btes = append(btes, '0')
			}
		}
	}

	conv, err := strconv.Atoi(zeroAllocBytes2String(btes))
	if err != nil {
		panic(err)
	}

	return int64(conv)
}

func zeroAllocBytes2String(bs []byte) (s string) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	stringHeader := reflect.StringHeader{Data: sliceHeader.Data, Len: sliceHeader.Len}
	return *(*string)(unsafe.Pointer(&stringHeader))
}
