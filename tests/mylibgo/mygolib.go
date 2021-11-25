// package name is different from dir path on purpose
package mygolib

//go:generate pygo

import (
	"fmt"
)

type MyStruct struct {
	AString string
}
type MyComplexStruct struct {
	AString string
	AStruct *MyStruct
}

/* this func is exported
 * @pygo.export
 */
func Test0() {
	return
}

//@pygo.export
func Test1(arg string) {
	return
}

//@pygo.export
func Test2(arg string) int {
	return 42
}

/* this func is exported
 * @pygo.export
 */
func Test3(arg1, arg2 string, arg3 int) string {
	return fmt.Sprintf("%s %s %d", arg1, arg2, arg3)
}

/* this func is exported
 * @pygo.export
 */
func Test4(arg1, arg2 string, arg3 []int) error {
	return nil
}

/* this func is exported
 * @pygo.export
 */
func Test5(arg1, arg2 string, arg3 int) []int {
	return []int{arg3}
}

/* this func is exported
 * @ pygo.export
 */
func Test6(arg1, arg2 string) []string {
	return []string{arg1, arg2}
}

/* this func is exported
 * @pygo.export
 */
func Test7(arg MyStruct) MyStruct {
	return MyStruct{}
}

/* this func is exported
 * @pygo.export
 */
func Test8(arg MyStruct) *MyComplexStruct {
	return nil
}
