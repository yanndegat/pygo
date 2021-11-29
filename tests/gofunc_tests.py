import unittest
import ctypes
import time

from pygo import gofunc, GoString, _map_ctype

# package name is different from dir path on purpose
from mylibgo.pygo import mygolib


class GoFuncTestCase(unittest.TestCase):

    def setUp(self):
        return

    def test_map_ctype(self):
        """Test string maps to ctypes.c_char_p"""

        result = _map_ctype("string")
        self.assertEqual(result, GoString)

        result = _map_ctype("c_char_p")
        self.assertEqual(result, ctypes.c_char_p)

        with self.assertRaises(Exception):
            result = _map_ctype("unknown")

    def test_func_sig(self):
        """Test call go func"""

        @gofunc(lib="simplelib.so", sig="c_char_p,c_char_p")
        def test1(): pass

        result = test1("world".encode('utf-8'))
        self.assertEqual(result.decode('utf-8'), "hello world")

    def test_func_args(self):
        """Test call go func"""

        @gofunc(lib="simplelib.so")
        def test1(c_char_p_1, *c_char_p): pass

        result = test1("world".encode('utf-8'))
        self.assertEqual(result.decode('utf-8'), "hello world")

    def test_func_enc(self):
        """Test call go func"""

        @gofunc(lib="simplelib.so")
        def test1(string_1, *string): pass

        result = test1("world")
        self.assertEqual(result, "hello world")

    def test_func_rename(self):
        """Test call go func"""

        @gofunc(lib="simplelib.so")
        def test1(string_1, *string): pass

        @gofunc(lib="simplelib.so", fname="test1")
        def test2(string_1, *string): pass

        result = test1("world")
        self.assertEqual(result, "hello world")

    def test_mylibgo_void(self):
        """Test call go func"""

        mygolib.Test0()

    def test_mylibgo_args(self):
        """Test call go func"""

        mygolib.Test1("world")

    def test_mylibgo_return_int(self):
        """Test call go func"""

        self.assertEqual(mygolib.Test2("world"), 42)

    def test_mylibgo_return_string(self):
        """Test call go func"""

        self.assertEqual(mygolib.Test3("hello", "world", 42), "hello world 42")

    def test_mylibgo_pass_int_array(self):
        """Test call go func"""
        mygolib.Test4("hello", "world", [4, 2])

    def test_mylibgo_return_int_array(self):
        """Test call go func"""
        arg = 4
        expect = [i*2 for i in range(arg)]
        self.assertEqual(mygolib.Test5("hello", "world", arg), expect)

    def test_mylibgo_return_byte_array(self):
        """Test call go func"""
        abyte = b'a'[0]
        self.assertEqual(mygolib.Test6(4, abyte), [abyte, abyte, abyte, abyte ])

    # def test_mylibgo_pass_return_bools(self):
    #     """Test call go func"""
    #     bools = [True, False, False, True]
    #     self.assertEqual(mygolib.Test7(bools, bools), bools+bools)

    def test_mylibgo_pass_return_bools(self):
        """Test call go func"""
        bools = [True, True, True, False, False, True]
        self.assertEqual(mygolib.Test7(bools, bools), bools+bools)


if __name__ == '__main__':
    unittest.main()
