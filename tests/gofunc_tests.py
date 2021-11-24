import unittest
import ctypes
from pygo import gofunc
from pygo import _map_ctype

# package name is different from dir path on purpose
from mylibgo.pygo import mygolib


class GoFuncTestCase(unittest.TestCase):

    def setUp(self):
        return

    def test_map_ctype(self):
        """Test string maps to ctypes.c_char_p"""

        result = _map_ctype("string")
        self.assertEqual(result, ctypes.c_char_p)

        result = _map_ctype("c_char_p")
        self.assertEqual(result, ctypes.c_char_p)

        with self.assertRaises(Exception):
            result = _map_ctype("error")

    def test_func_sig(self):
        """Test call go func"""

        @gofunc(lib="mygolib.so", sig="c_char_p,c_char_p")
        def test1(): pass

        result = test1("world".encode('utf-8'))
        self.assertEqual(result.decode('utf-8'), "hello world")

    def test_func_args(self):
        """Test call go func"""

        @gofunc(lib="mygolib.so")
        def test1(c_char_p_1, *c_char_p): pass

        result = test1("world".encode('utf-8'))
        self.assertEqual(result.decode('utf-8'), "hello world")

    def test_func_enc(self):
        """Test call go func"""

        @gofunc(lib="mygolib.so")
        def test1(string_1, *string): pass

        result = test1("world")
        self.assertEqual(result, "hello world")

    def test_func_rename(self):
        """Test call go func"""

        @gofunc(lib="mygolib.so")
        def test1(string_1, *string): pass

        @gofunc(lib="mygolib.so", fname="test1")
        def test2(string_1, *string): pass

        result = test1("world")
        self.assertEqual(result, "hello world")

    def test_mylibgo(self):
        """Test call go func"""

        mygolib.Test0()
        mygolib.Test1("world")
        print(mygolib.Test2("world"))
        print(mygolib.Test3("hello", "world", 0))



if __name__ == '__main__':
    unittest.main()
