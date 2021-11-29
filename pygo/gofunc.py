import inspect
import re
import os

import ctypes
from threading import RLock

_LIBS = {}
_LIBS_LOCK = RLock()


class GoSlice(ctypes.Structure):
    _fields_ = [("data", ctypes.POINTER(ctypes.c_void_p)),
                ("len", ctypes.c_longlong), ("cap", ctypes.c_longlong)]


class CSlice(ctypes.Structure):
    _fields_ = [("data", ctypes.POINTER(ctypes.c_void_p)),
                ("len", ctypes.c_longlong), ("cap", ctypes.c_longlong)]


class GoString(ctypes.Structure):
    _fields_ = [("p", ctypes.c_char_p), ("n", ctypes.c_longlong)]


class gofunc(object):
    """
    gofunc annotation decorates any func to call a go
    lib func with the same name.

    Example:

    ```
    //mylib.go
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
    will give

    ```
    //mylib.py
    import pygo
    @pygo.gofunc(lib="mygolib.so")
    def myGoFunc(string_1, *string): pass
    ```

    The signature of the python method is used to infer golang func argument
    types and return type. You can use _IX to index the arg names if types have
    to be repeated.

    Examples:


    go sig -> python sig

    - (carg *C.char) {} -> (string): pass
    - (carg *C.char) {} -> (string, *void): pass
    - (carg *C.char) {} -> (string, *void): pass
    - (carg *C.char) *C.char {} -> (c_char_p_1, *c_char_p): pass
    - (carg *C.char) *C.char {} -> (c_char_p_1, *c_char_p): pass
    - (carg1, carg2 *C.char) *C.char {} -> (c_char_p_1, *c_char_p): pass
    - (carg1, carg2 *C.char) *C.char {} -> (string_1, string_2, *string): pass

    The annotation takes parameters

    :param lib: The name of the golang lib to load.
                It must represent a valid .so library
                present in `libPath`.
    :type lib: string

    :param libPath: The path from where the golang lib will be
                    to loaded. By default, if will load the lib
                    in the same path as the calling function.
    :type libPath: string

    :param sig: The type signature of the golang lib func.
                Will override the signature of the python function.
                The last type will be used as the return type.
                If the go func has arguments and the return type is
                void, you have to specify "void" in the sig type list.

    :type sig: string

    :param freeMem: The name of a func which has to be present
                in your lib, and whose signature has to be
                `func (*C.void)`.
                It'll be called which each pointer returned by your lib
                so that mem can be freed once data has been passed to
                python.

    :type sig: string

    :param fname: The go func name to use.
                  Will override the name of the python function.

    :type fname: string
    """

    def __init__(self,
                 lib=None,
                 libPath=None,
                 sig=None,
                 fname=None,
                 freeMemFunc="freeMem"):
        if lib is None or not isinstance(lib, str):
            raise Exception("lib is mandatory and has to be a string"
                            " representing the file path of a go lib.")
        if libPath is not None and not os.path.exists(libPath):
            raise Exception("libPath must represent an existing file path.")
        if sig is not None and not isinstance(sig, str):
            raise Exception("sig has to be a string representing the type"
                            " signature of a go lib func.")
        if freeMemFunc is None or freeMemFunc == "":
            raise Exception("freeMemFunc is mandatory and must represent"
                            " the name of a `func (s *C.void)` in your go lib.")
        if fname is not None and not isinstance(fname, str):
            raise Exception("fname has to be a string representing a valid"
                            " function name of a go lib func.")

        self.lib = lib
        self.libPath = libPath
        self.fname = fname
        self.sig = sig
        self.freeMemFunc = freeMemFunc

        return

    def __call__(self, f):
        try:
            libPath = self.libPath
            if libPath is None:
                libPath = os.path.dirname(f.__code__.co_filename)
            self.lib = _load_lib(libPath, self.lib)
        except Exception as e:
            raise e

        # method name is overriden by annotation "fname" arg
        if self.fname is None:
            self.fname = f.__name__

        try:
            self.freeMem = getattr(self.lib, self.freeMemFunc)
            self.freeMem.argtypes = [ctypes.c_void_p]
            self.freeMem.restype = ctypes.c_void_p
        except AttributeError:
            raise AttributeError(
                f"func {self.freeMemFunc} not found in {self.lib}")

        try:
            self.func = getattr(self.lib, self.fname)
        except AttributeError:
            raise AttributeError(
                f"func {self.fname} not found in {self.lib}")

        # method signature is overriden by annotation "sig" arg
        if self.sig is None:
            argspec = inspect.getfullargspec(f)
            self.sig = [_trim_sigtype(t) for t in argspec[0]]
            if argspec[1] is None:
                self.sig.append("c_void_p")
            else:
                self.sig.append(_trim_sigtype(argspec[1]))
        else:
            self.sig = [_trim_sigtype(t) for t in self.sig.split(",")]

        self.func.argtypes = [_map_ctype(t) for t in self.sig[:-1]]
        self.func.restype = _map_ret_ctype(self.sig[-1])
        self.conv = [_map_conv(t) for t in self.sig[:-1]]

        def wrapped_f(*args):
            conv_args = [self.conv[i](arg) for i, arg in enumerate(args)]
            return self._handle_ret_value(self.func(*conv_args), self.sig[-1])

        return wrapped_f

    def _handle_ret_value(self, value, valueType, enc="utf-8"):
        if value is None:
            return None

        if valueType == "string" \
           or valueType == "c_char_p" \
           or valueType == "error":
            res = ctypes.cast(value, ctypes.c_char_p).value
            self.freeMem(value)
        elif valueType == "arr_int":
            cslice = ctypes.cast(value, ctypes.POINTER(CSlice)).contents
            arr = ctypes.cast(cslice.data, ctypes.POINTER(ctypes.c_int64*cslice.len))
            res = [x for x in arr.contents]
            self.freeMem(ctypes.cast(arr, ctypes.c_void_p))
            self.freeMem(ctypes.cast(value, ctypes.c_void_p))
        elif valueType == "arr_byte":
            cslice = ctypes.cast(value, ctypes.POINTER(CSlice)).contents
            arr = ctypes.cast(cslice.data, ctypes.POINTER(ctypes.c_byte*cslice.len))
            res = [x for x in arr.contents]
            self.freeMem(ctypes.cast(arr, ctypes.c_void_p))
            self.freeMem(ctypes.cast(value, ctypes.c_void_p))
        elif valueType == "arr_bool":
            cslice = ctypes.cast(value, ctypes.POINTER(CSlice)).contents
            arr = ctypes.cast(cslice.data, ctypes.POINTER(ctypes.c_bool*cslice.len))
            res = [x for x in arr.contents]
            self.freeMem(ctypes.cast(arr, ctypes.c_void_p))
            self.freeMem(ctypes.cast(value, ctypes.c_void_p))
        else:
            res = value

        if valueType == "string":
            return res.decode(enc)

        return res


def _no_conv(v):
    return v


def _string_conv(v):
    return GoString(v.encode("utf-8"), len(v))


def _arr_conv(t):
    ct = _map_ctype(t)
    def __conv(v):
        s = len(v)
        data = (ct * s)(*v)
        gslice = GoSlice(ctypes.cast(data,ctypes.POINTER(ctypes.c_void_p)), s, s)
        return gslice

    return __conv


def _map_conv(t):
    if t == "string":
        return _string_conv
    if _is_array_type(t):
        return _arr_conv(_array_type(t))

    return _no_conv


def _is_array_type(t):
    # return true if t has prefix 'arr_'
    if t is not None and isinstance(t, str):
        return re.match('^arr_', t.strip())
    else:
        raise Exception(f"sigtype must be a valid string: {t}")


def _array_type(t):
    # remove prefix 'arr_' in strings
    if t is not None and isinstance(t, str):
        return re.sub('^arr_', '', t.strip())
    else:
        raise Exception(f"sigtype must be a valid string: {t}")


def _trim_sigtype(t):
    # remove ending indexes such as _1, _2 in arg types
    # removes stars from types
    if t is not None and isinstance(t, str):
        return re.sub('(\*|_[0-9]+)', '', t.strip())
    else:
        raise Exception(f"sigtype must be a valid string: {t}")


_ctypes = inspect.getmembers(ctypes, lambda a: not(inspect.isroutine(a)))


def _map_ret_ctype(t):
    if t == "string" or t == "error" or t == "c_char_p":
        return ctypes.POINTER(ctypes.c_char)
    if _is_array_type(t):
        return ctypes.c_size_t
    return _map_ctype(t)

def _map_ctype(t):
    if t == "bool":
        return ctypes.c_bool
    elif t == "byte":
        return ctypes.c_byte
    elif t == "char":
        return ctypes.c_char
    elif t == "int":
        return ctypes.c_int
    elif t == "int64":
        return ctypes.c_int64
    elif t == "long":
        return ctypes.c_long
    elif t == "float":
        return ctypes.c_float
    elif t == "double":
        return ctypes.c_double
    elif t == "string":
        return GoString
    elif t == "error":
        return GoString
    elif t == "void":
        return ctypes.c_void_p
    elif _is_array_type(t):
        return GoSlice
    else:
        if isinstance(t, str) and t.startswith('c_'):
            found = [a[1] for a in _ctypes if a[0] == t]
            if len(found) == 1:
                return found[0]

        raise Exception(f"unkwon type {t}.")


def _load_lib(libPath, lib):
    with _LIBS_LOCK:
        if lib not in _LIBS:
            _LIBS[lib] = ctypes.cdll.LoadLibrary(os.path.join(".", libPath, lib))

        return _LIBS[lib]
