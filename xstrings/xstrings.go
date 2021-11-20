package xstrings

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

var (
	ErrIncorrectStmts = errors.New("incorrect statements")
	ErrIncorrectDecl  = errors.New("incorrect declaration")
	ErrIncorrectExpr  = errors.New("incorrect expration")
)

func ToStmts(s string) ([]ast.Stmt, error) {
	s = "package p\nfunc f(){\n" + s + "\n}"

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", s, 0)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrIncorrectStmts, err.Error())
	}

	_func, ok := f.Decls[0].(*ast.FuncDecl)
	if !ok {
		return nil, ErrIncorrectStmts
	}

	return _func.Body.List, nil
}

func ToDecl(s string) ([]ast.Decl, error) {
	fset := token.NewFileSet()
	s = "package p\n" + s

	expr, err := parser.ParseFile(fset, "", s, 0)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrIncorrectDecl, err.Error())
	}

	return expr.Decls, nil
}

func ToExpr(s string) (ast.Expr, error) {
	ret, err := parser.ParseExpr(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrIncorrectExpr, err.Error())
	}

	return ret, nil
}

func FromNode(n ast.Node) string {
	var buf bytes.Buffer

	if err := format.Node(&buf, token.NewFileSet(), n); err != nil {
		panic(fmt.Sprintf("format error: %s", err.Error()))
	}

	return buf.String()
}
