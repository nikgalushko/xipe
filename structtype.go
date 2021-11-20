package xipe

import "go/ast"

type StructType struct {
	n *ast.StructType
}

type Field struct {
	Name, Type, Tag string
}

func (s StructType) AppendField(f Field) {
	s.n.Fields.List = append(s.n.Fields.List, &ast.Field{
		Names: []*ast.Ident{{Name: f.Name}},
		Type:  &ast.Ident{Name: f.Name},
		Tag:   &ast.BasicLit{Value: f.Tag},
	})
}

func (s StructType) FieldExists(name string) bool {
	for _, f := range s.n.Fields.List {
		if f.Names[0].Name == name {
			return true
		}
	}

	return false
}
