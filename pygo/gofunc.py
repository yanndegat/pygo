import ctypes
import inspect
import re
from os.path import exists


class gofunc(object):
    """
    gofunc annotation decorates any func to call a go
    lib func with the same name

    Ex:
    ```
    // mylib.go
    func myfunc(a, b string) {
       ....
    }
    ```
    will give

    ```
    //mylib.py
    import pygo
    @pygo.gofunc(lib="mylib.so", sig=string,string,void)
    ```

    :param lib: The name of the golang lib.
    :type lib: string

    :param sig: The type signature of the golang lib func.
    :type lib: string

    :return: The result of the multiplication.
    :rtype: int
    """

    def __init__(self, lib=None, sig=None, fname=None):
        if lib is None or not isinstance(lib, str) or not exists(lib):
            raise Exception("lib is mandatory and has to be a string"
                            " representing the file path of a go lib.")
        if sig is not None and not isinstance(sig, str):
            raise Exception("sig has to be a string representing the type"
                            " signature of a go lib func.")
        if fname is not None and not isinstance(fname, str):
            raise Exception("fname has to be a string representing a valid"
                            " function name of a go lib func.")

        self.libPath = lib
        self.fname = fname
        self.sig = sig

        return

    def __call__(self, f):
        try:
            self.lib = ctypes.cdll.LoadLibrary(self.libPath)
        except Exception as e:
            raise e

        # method name is overriden by annotation "fname" arg
        if self.fname is None:
            self.fname = f.__name__

        try:
            self.func = getattr(self.lib, self.fname)
        except AttributeError:
            raise AttributeError(
                f"func {self.f.name} not found in {self.libPath}")

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
        self.func.restype = _map_ctype(self.sig[-1])

        def wrapped_f(*args):
            enc_args = [_enc_type_value(arg, self.sig[i])
                        for i, arg in enumerate(args)]

            return _dec_type_value(self.func(*enc_args), self.sig[-1])

        return wrapped_f


def _enc_type_value(value, type, enc="utf-8"):
    if type == "string":
        return value.encode(enc)

    return value


def _dec_type_value(value, type, enc="utf-8"):
    if type == "string":
        return value.decode(enc)

    return value


def _trim_sigtype(t):
    # remove ending indexes such as _1, _2 in arg types
    return re.sub('_[0-9]+', '', t)


_ctypes = inspect.getmembers(ctypes, lambda a: not(inspect.isroutine(a)))


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
        return ctypes.c_char_p
    elif t == "void":
        return ctypes.c_void_p
    else:
        if isinstance(t, str) and t.startswith('c_'):
            found = [a[1] for a in _ctypes if a[0] == t]
            if len(found) == 1:
                return found[0]

        raise Exception(f"unkwon type {t}.")
