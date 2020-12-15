package main

import (
	"C"
	"fmt"
)

//export test1
func test1(carg *C.char) *C.char {
	goarg := C.GoString(carg)
	return C.CString(fmt.Sprintf("hello %s", goarg))
}

func main() {}
