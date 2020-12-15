import unittest
import ctypes
from pygo import gofunc
from pygo import _map_ctype


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

    def test_func(self):
        """Test call go func"""

        @gofunc(lib="tests/mygolib.so", sig="c_char_p,c_char_p")
        def test1(hello):
            return

        result = test1("world".encode('utf-8'))
        self.assertEqual(result.decode('utf-8'), "hello world")


if __name__ == '__main__':
    unittest.main()