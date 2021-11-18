# PyGo
A small library exposising a helper decorator to ease the call of a go shared library
from python.

### Installation
```
pip install pygo
```

### Get started
How to call go from python with this lib:

```Python
import pygo

@pygo.gofunc(lib="mygolib.so", sig="string,string,void")
def myGoFunc():
    return
```
