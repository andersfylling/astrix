package astrix

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

type primitiveType int

const (
	TypeUnknown primitiveType = iota
	TypeInt
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUint
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
	TypeFloat32
	TypeFloat64
	TypeStruct
	TypeMap
	TypeString
	TypeSlice
	TypeInterface
	TypeBool
	TypePointer
)

// GenerateTypes generates a list of types that can be reflected upon.
func GetTypes(path string, excludeFilesRegexp *regexp.Regexp, excludeTypesRegexp *regexp.Regexp) (types []*TypeInfo, err error) {
	files, err := FindFiles(path, excludeFilesRegexp)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileData, err := parser.ParseFile(token.NewFileSet(), file.Name(), nil, 0)
		if err != nil {
			return nil, err
		}

		fileTypes, err := GetTypesFromFile(fileData, excludeTypesRegexp)
		if err != nil {
			return nil, err
		}

		types = append(types, fileTypes...)
	}

	return types, nil
}

func GetTypesFromFile(file *ast.File, exclude *regexp.Regexp) (types []*TypeInfo, err error) {
	for _, item := range file.Decls {
		var gdecl *ast.GenDecl
		var ok bool
		if gdecl, ok = item.(*ast.GenDecl); !ok || gdecl.Tok != token.TYPE {
			continue
		}

		specs := item.(*ast.GenDecl).Specs
		for i := range specs {
			ts := specs[i].(*ast.TypeSpec)
			if ts.Name.Name == "" {
				continue
			}

			t := &TypeInfo{
				T:                  nil,
				PrimitiveType:      TypeUnknown,
				DefaultValueString: "",
				Name:               ts.Name.Name,
				Fields:             nil,
				Exported:           false,
			}
			if strings.ToUpper(t.Name[:1]) == t.Name[:1] {
				t.Exported = true
			}

			types = append(types, t)
		}
	}

	return types, nil
}

type TypeInfo struct {
	T                  *TypeInfo
	PrimitiveType      primitiveType
	DefaultValueString string
	Name               string
	Fields             map[string]*FieldInfo
	Exported           bool
	//Interfaces []*TypeInfo
}

func (t *TypeInfo) Is(pt primitiveType) bool {
	if t.T == nil {
		return t.PrimitiveType == pt
	} else {
		return t.T.Is(pt)
	}
}

type FieldInfo struct {
	Pointer bool
	*TypeInfo
}
