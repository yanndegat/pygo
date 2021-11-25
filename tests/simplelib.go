package main

// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"
)

//export test1
func test1(carg *C.char) *C.char {
	goarg := C.GoString(carg)
	return C.CString(fmt.Sprintf("hello %s", goarg))
}

//export freeCString
func freeCString(c *C.char) {
	C.free(unsafe.Pointer(c))
}

func main() {}
