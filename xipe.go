package xipe

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"

	"golang.org/x/tools/go/ast/astutil"
)

type Xipe struct {
	fset *token.FileSet
	file *ast.File
}

func New(path string) (Xipe, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return Xipe{}, err
	}

	return Xipe{fset: fset, file: file}, nil
}

func (x Xipe) FindStructType(name string) (StructType, bool) {
	ret := StructType{}
	found := false

	astutil.Apply(x.file, nil, func(c *astutil.Cursor) bool {
		if node, ok := c.Node().(*ast.TypeSpec); ok && node.Name.Name == name {
			ret.n, found = node.Type.(*ast.StructType)
			return false
		}

		return true
	})

	return ret, found
}

func (x Xipe) FindFunc(name string) (Func, bool) {
	ret := Func{}
	found := false

	astutil.Apply(x.file, nil, func(c *astutil.Cursor) bool {
		if node, ok := c.Node().(*ast.FuncDecl); ok && node.Name.Name == name {
			ret.n = node
			found = true

			return false
		}

		return true
	})

	return ret, found
}

func (x Xipe) GetAllTypeSpecs() map[string]*ast.TypeSpec {
	ret := make(map[string]*ast.TypeSpec)

	astutil.Apply(x.file, nil, func(c *astutil.Cursor) bool {
		if node, ok := c.Node().(*ast.TypeSpec); ok {
			ret[node.Name.Name] = node
		}

		return true
	})

	return ret
}

func (x Xipe) AddImports(imports ...string) {
	for _, s := range imports {
		astutil.AddImport(x.fset, x.file, s)
	}
}

func (x Xipe) AddDecls(dscls ...ast.Decl) {
	x.file.Decls = append(x.file.Decls, dscls...)
}

func (x Xipe) Write(w io.Writer) error {
	return printer.Fprint(w, x.fset, x.file)
}
