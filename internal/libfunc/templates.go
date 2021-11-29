package libfunc

import (
	"bytes"
	"text/template"
)

var resToSliceTpl = template.Must(template.New("").Parse(`
    res := {{.FunCall}}
    s := len(res)
    p := C.malloc(C.size_t(s) * C.size_t(unsafe.Sizeof(uintptr(0))))
    cslice := (C.CSliceP)(C.malloc(C.sizeof_CSlice))
    cslice.len = C.CInt64(s)
    cslice.cap = C.CInt64(s)
    cslice.data = p
    pp := (*[1<<30 - 1]{{.SliceCType}})(p)
    copy(pp[:], res)
    ptr := (C.CSliceP)(unsafe.Pointer(cslice))
    return ptr
`))

func (f *Func) convertResToSlice() (string, error) {
	data := struct {
		SliceCType Type
		FunCall    string
	}{
		SliceCType: f.Result.T().ToCType(),
		FunCall:    f.GoFuncCall(),
	}

	var tpl bytes.Buffer
	if err := resToSliceTpl.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
