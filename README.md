# Medium multiply
A small demo library for a Medium publication about publishing libraries.

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
```
