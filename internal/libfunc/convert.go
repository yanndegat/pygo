package libfunc

import (
	"fmt"
	"go/ast"
	"reflect"

	iast "github.com/yanndegat/pygo/internal/ast"
)

func ConvertFromAstF(lib string, astF *iast.AstFunc) (*Func, error) {
	if astF == nil {
		return nil, nil
	}

	f := &Func{
		Lib:    lib,
		Name:   astF.Name,
		Args:   []Arg{},
		Result: TypeVoid,
	}

	if astF.Params != nil {
		for _, param := range astF.Params {
			t, err := astTypeToType(param.Type)
			if err != nil {
				return nil, err
			}

			for i := 0; i < len(param.Names); i++ {
				f.Args = append(f.Args, Arg{Name: param.Names[i].Name, Type: *t})
			}
		}
	}

	if astF.Results != nil {
		for _, result := range astF.Results {
			t, err := astTypeToType(result.Type)
			if err != nil {
				return nil, err
			}

			nbRets := len(result.Names)
			// if no name, then there's one unnamed return value
			if nbRets == 0 {
				nbRets = 1
			}

			if nbRets > 1 {
				return nil, fmt.Errorf("exported func can have 0 or 1 return value.")
			}

			// if f.Result is set to something else than TypeVoid, it's because we have more
			// than 1 elt in astF.Results.
			if f.Result != TypeVoid {
				return nil, fmt.Errorf("exported func can have 0 or 1 value returned.")
			}
			f.Result = *t
		}
	}

	return f, nil
}

func checkType(t interface{}) (bool, error) {
	if expr, ok := t.(*ast.Ident); ok {
		return validType(Type(expr.Name)), nil
	} else if expr, ok := t.(*ast.SelectorExpr); ok {
		xT, err := astTypeToType(expr.X)
		if err != nil {
			return false, err
		}
		selT, err := astTypeToType(expr.Sel)
		if err != nil {
			return false, err
		}
		return validType(Type(fmt.Sprintf("%v.%v", xT, selT))), nil
	} else if expr, ok := t.(*ast.StarExpr); ok {
		return checkType(expr.X)
	} else if expr, ok := t.(*ast.ArrayType); ok {
		return checkType(expr.Elt)
	} else if expr, ok := t.(*ast.MapType); ok {
		kT, err := checkType(expr.Key)
		if err != nil {
			return false, err
		}
		vT, err := checkType(expr.Value)
		if err != nil {
			return false, err
		}
		return kT && vT, nil
	} else {
		return false, fmt.Errorf("unsupported ast type: %v", reflect.TypeOf(t))
	}

	return false, nil
}

func astTypeToType(t interface{}) (*Type, error) {
	var res Type
	if expr, ok := t.(*ast.Ident); ok {
		res = Type(expr.Name)
	} else if expr, ok := t.(*ast.SelectorExpr); ok {
		tStr, err := astTypeToType(expr.X)
		if err != nil {
			return nil, err
		}
		tStr2, err := astTypeToType(expr.Sel)
		if err != nil {
			return nil, err
		}
		res = Type(fmt.Sprintf("%s.%s", *tStr, *tStr2))
	} else if expr, ok := t.(*ast.StarExpr); ok {
		tStr, err := astTypeToType(expr.X)
		if err != nil {
			return nil, err
		}
		res = Type(fmt.Sprintf("*%s", *tStr))
	} else if expr, ok := t.(*ast.ArrayType); ok {
		tStr, err := astTypeToType(expr.Elt)
		if err != nil {
			return nil, err
		}
		res = Type(fmt.Sprintf("[]%s", *tStr))
	} else if expr, ok := t.(*ast.MapType); ok {
		keyTStr, err := astTypeToType(expr.Key)
		if err != nil {
			return nil, err
		}
		valueTStr, err := astTypeToType(expr.Value)
		if err != nil {
			return nil, err
		}
		res = Type(fmt.Sprintf("map[%s]%s", *keyTStr, *valueTStr))
	} else {
		return nil, fmt.Errorf("unsupported ast type: %v", reflect.TypeOf(t))
	}

	return &res, nil
}

// func (f *Func) PySig() string {
// 	sigArgs := []string{}

// 	if f.Args != nil {
// 		for _, arg := range f.Args {
// 			typeStr, _ := astTypeToType(arg.Type)
// 			for i := 0; i < len(arg.Names); i++ {
// 				sigArgs = append(sigArgs, typeStr)
// 			}
// 		}
// 	}

// 	sigResults := []string{}
// 	if f.Results != nil {
// 		for _, result := range f.Results {
// 			typeStr, _ := astTypeToType(result.Type)
// 			for i := 0; i < len(result.Names); i++ {
// 				if typeStr != "error" {
// 					sigResults = append(sigResults, typeStr)
// 				}
// 			}
// 		}
// 	}

// 	return strings.Join(args, ",")

// }
