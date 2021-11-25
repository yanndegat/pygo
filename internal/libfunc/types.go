package libfunc

import (
	"fmt"
	"regexp"
)

var (
	trimTypeRe    = regexp.MustCompile(`(\[\]|\*| )`)
	arrayTypeRe   = regexp.MustCompile(`^(\[\])(.*)`)
	pointerTypeRe = regexp.MustCompile(`^(\*)(.*)`)

	TypeBool   Type = "bool"
	TypeCChar  Type = "*C.char"
	TypeError  Type = "error"
	TypeInt    Type = "int"
	TypeInt32  Type = "int32"
	TypeInt64  Type = "int64"
	TypeString Type = "string"
	TypeVoid   Type = "void"

	ValidTypes = []Type{
		TypeBool,
		TypeBool,
		TypeError,
		TypeInt,
		TypeInt64,
		TypeInt32,
		TypeString,
		TypeVoid,
	}

	SupportedTypes = []Type{
		TypeBool,
		TypeError,
		TypeInt,
		TypeInt32,
		TypeInt64,
		TypeString,
		TypeVoid,
	}

	GoTypeToCTypes = map[Type]Type{
		TypeBool:   TypeBool,
		TypeError:  TypeCChar,
		TypeInt:    TypeInt,
		TypeInt32:  TypeInt32,
		TypeInt64:  TypeInt64,
		TypeString: TypeCChar,
	}
)

type Type string

func (t Type) ToCType() Type {
	if t.IsArray() {
		arrayType := Type(arrayTypeRe.ReplaceAllString(string(t), "$2")).ToCType()
		return Type(fmt.Sprintf("[]%s", arrayType))
	}
	if t.IsPointer() {
		pointerType := Type(pointerTypeRe.ReplaceAllString(string(t), "$2")).ToCType()
		return Type(fmt.Sprintf("*%s", pointerType))
	}

	return GoTypeToCTypes[t.T()]
}

func (t Type) IsArray() bool {
	return arrayTypeRe.MatchString(string(t))
}

func (t Type) IsPointer() bool {
	return pointerTypeRe.MatchString(string(t))
}

func (t Type) T() Type {
	if t.IsArray() {
		return Type(arrayTypeRe.ReplaceAllString(string(t), "$2")).T()
	}
	if t.IsPointer() {
		return Type(pointerTypeRe.ReplaceAllString(string(t), "$2")).T()
	}
	return t
}

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
		if t.T() == T {
			return true
		}
	}
	return false
}
