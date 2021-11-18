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

@pygo.gofunc(lib="mygolib.so", sig="string,string")
def myGoFunc():
    return
    
if __name__ == '__main__':
    res = myGoFunc("world".encode('utf-8'))
    print( res.decode('utf-8'))
```

### , run it

``` sh
$ python3 main.py
hello world
```


## Motivation

Hopefully, this lib could help one convert a python codebase to go incrementally.
