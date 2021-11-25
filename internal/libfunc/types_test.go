package libfunc

import (
	"fmt"
	"testing"
)

func TestMain_Type_T(t *testing.T) {

	tests := []struct {
		Type Type
		T    Type
	}{
		{
			"string",
			"string",
		},
		{
			"*string",
			"string",
		},
		{
			"[]*string",
			"string",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if test.Type.T() != test.T {
				t.Fatalf("match should be %v, was %v", test.T, test.Type.T())
			}
		})
	}
}
