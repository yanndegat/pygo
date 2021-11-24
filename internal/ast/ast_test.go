package ast

import (
	"fmt"
	"testing"
)

func TestMain_commentFuncExport(t *testing.T) {

	tests := []struct {
		Text  string
		Match bool
	}{
		{
			`//@pygo.export`,
			true,
		},
		{

			`// @pygo.export`,
			true,
		},
		{
			`/* this func is exported
 * @pygo.export
 */`,
			true,
		},
		{
			`/* this func is not exported
 * @pygo.exporti
 */`,
			false,
		},
		{
			`/* this func is exported
 * hi, @pygo.export, ok
 */`,
			true,
		},
		{
			`/* this func is not exported
 * hi, -@pygo.export, ok
 */`,
			false,
		},
		{
			`/* this func is not exported
 * hi, //@pygo.export, ok
 */`,
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			match, err := commentFuncExport(test.Text)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if match != test.Match {
				t.Fatalf("match should be %v, was %v", test.Match, match)
			}
		})
	}
}
