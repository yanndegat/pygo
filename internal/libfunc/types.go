package libfunc

import (
	"regexp"
)

type Type string

func (t Type) ToCType() Type {
	return GoTypeToCTypes[t]
}

var (
	TypeBool   Type = "bool"
	TypeCChar  Type = "*C.char"
	TypeError  Type = "error"
	TypeInt    Type = "int"
	TypeInt64  Type = "int64"
	TypeString Type = "string"
	TypeVoid   Type = "void"

	ValidTypes = []Type{
		TypeBool,
		TypeBool,
		TypeError,
		TypeInt,
		TypeInt64,
		TypeString,
		TypeVoid,
	}

	SupportedTypes = []Type{
		TypeBool,
		TypeError,
		TypeInt,
		TypeInt64,
		TypeString,
		TypeVoid,
	}

	GoTypeToCTypes = map[Type]Type{
		TypeBool:   TypeBool,
		TypeError:  TypeCChar,
		TypeInt:    TypeInt,
		TypeInt64:  TypeInt64,
		TypeString: TypeCChar,
	}
)

var trimTypeRe = regexp.MustCompile(`(\[\]|\*| )`)

func validType(t Type) bool {
	t = Type(trimTypeRe.ReplaceAllString(string(t), ""))

	for _, T := range ValidTypes {
		if t == T {
			return true
		}
	}
	return false
}

func supportedType(t Type) bool {
	for _, T := range SupportedTypes {
		if t == T {
			return true
		}
	}
	return false
}
