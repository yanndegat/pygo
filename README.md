# PyGo [WIP]

A small library exposing a helper decorator to ease the call of a go shared library
from python.

## Installation
```
pip install pygo
```

## Get started
How to call go from python with this lib:

### edit your go lib `mygolib.go`:
``` Go
package main

import (
	"C"
	"fmt"
)

//export myGoFunc
func myGoFunc(carg *C.char) *C.char {
	goarg := C.GoString(carg)
	return C.CString(fmt.Sprintf("hello %s", goarg))
}

func main() {}
```

### , build it:

``` sh
$ go build -o ./mygolib.so -buildmode=c-shared mygolib.go
```

### , use it

edit your main.py
```Python
import pygo

@pygo.gofunc(lib="mygolib.so")
def myGoFunc(string_1, *string): pass

@pygo.gofunc(lib="mygolib.so", sig="string,string", fname="myGoFunc")
def test1(): pass

@gofunc(lib="tests/mygolib.so", fname="myGoFunc")
def test2(c_char_p_1, *c_char_p): pass


if __name__ == '__main__':
    res = myGoFunc("world")
    print( res))
    
    res = test1("world")
    print( res))

    res = test2("world".encode('utf-8'))
    print( res.decode('utf-8'))
```

### , run it

``` sh
$ python3 main.py
hello world
hello world
hello world
```


## Motivation

Hopefully, this lib could help one convert a python codebase to go incrementally.
