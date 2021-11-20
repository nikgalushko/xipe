package xipe

import "go/ast"

type Func struct {
	n *ast.FuncDecl
}

func (f Func) LastReturnStmt() *ast.ReturnStmt {
	list := f.n.Body.List
	for i := len(list) - 1; i >= 0; i-- {
		if stmt, ok := list[i].(*ast.ReturnStmt); ok {
			return stmt
		}
	}

	return nil
}

func (f Func) AppendBeforeReturn(stmts []ast.Stmt) {
	list := f.n.Body.List
	returnStmt := list[len(list)-1]
	list = append(list[:len(list)-1], stmts...)

	f.n.Body.List = append(list, returnStmt)
}
