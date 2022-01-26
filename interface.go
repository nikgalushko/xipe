package xipe

import (
	"go/ast"

	"github.com/nikgalushko/xipe/xstrings"
)

type Interface struct {
	n *ast.InterfaceType
}

func (i Interface) AppendMethod(name string, params []Field, results []Field) error {
	function := &ast.FuncType{
		Params:  &ast.FieldList{},
		Results: &ast.FieldList{},
	}

	for _, p := range params {
		expr, err := xstrings.ToExpr(p.Type)
		if err != nil {
			return err
		}
		function.Params.List = append(function.Params.List, &ast.Field{
			Names: []*ast.Ident{{Name: p.Name}},
			Type:  expr,
		})
	}

	i.n.Methods.List = append(i.n.Methods.List, &ast.Field{
		Names: []*ast.Ident{{Name: name}},
		Type:  function,
	})

	return nil
}
