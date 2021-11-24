package libfunc

import (
	"fmt"
	"strings"
)

type Func struct {
	Lib    string
	Name   string
	Args   []Arg
	Result Type
}

func (f Func) IsSupported() bool {
	for _, a := range f.Args {
		if !supportedType(a.Type) {
			return false
		}
	}

	return supportedType(f.Result)
}

func (f *Func) String() string {
	return fmt.Sprintf("%s.%s: %v -> %v", f.Lib, f.Name, f.Args, f.Result)
}

func (f *Func) PySig() string {
	sig := []string{}
	for i, arg := range f.Args {
		sig = append(sig, fmt.Sprintf("%s_%d", arg.Type, i))
	}
	if f.Result != TypeVoid {
		// star means kwargs, which is a special case
		// interpreted by pygo to infer return type
		sig = append(sig, fmt.Sprintf("*%s", f.Result))
	}

	return strings.Join(sig, ", ")
}

func (f *Func) GoSigArgs() string {
	sig := make([]string, len(f.Args))
	for i, arg := range f.Args {
		sig[i] = fmt.Sprintf("%s %s", arg.Name, arg.Type.ToCType())
	}
	return strings.Join(sig, ", ")
}

func (f *Func) GoFuncCall() string {
	args := make([]string, len(f.Args))
	for i, arg := range f.Args {
		args[i] = string(arg.ToGoType())
	}

	call := fmt.Sprintf("%s.%s(%s)", f.Lib, f.Name, strings.Join(args, ", "))
	if f.Result == TypeString {
		call = fmt.Sprintf("C.CString(%s)", call)
	}
	return call

}

func (f *Func) IsVoid() bool {
	return f.Result == TypeVoid
}

func (f *Func) GoSigRet() string {
	if f.Result == TypeVoid {
		return ""
	}

	return fmt.Sprintf("%s", f.Result.ToCType())
}

type Arg struct {
	Name string
	Type Type
}

func (a Arg) String() string {
	return fmt.Sprintf("%s:%s", a.Name, a.Type)
}

func (a Arg) ToGoType() string {
	if a.Type == TypeString {
		return fmt.Sprintf("C.GoString(%s)", a.Name)
	}

	return a.Name
}
