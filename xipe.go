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
	FileSet *token.FileSet
	file    *ast.File
}

func New(path string) (Xipe, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return Xipe{}, err
	}

	return Xipe{FileSet: fset, file: file}, nil
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

func (x Xipe) FindInterface(name string) (Interface, bool) {
	ret := Interface{}
	found := false

	astutil.Apply(x.file, nil, func(c *astutil.Cursor) bool {
		if node, ok := c.Node().(*ast.TypeSpec); ok && node.Name.Name == name {
			ret.n, found = node.Type.(*ast.InterfaceType)
			return false
		}

		return true
	})

	return ret, found
}

func (x Xipe) FindInterfaces() []Interface {
	var ret []Interface

	astutil.Apply(x.file, nil, func(c *astutil.Cursor) bool {
		if node, ok := c.Node().(*ast.TypeSpec); ok {
			n, ok := node.Type.(*ast.InterfaceType)
			if ok {
				ret = append(ret, Interface{n: n})
			}
		}

		return true
	})

	return ret
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
		astutil.AddImport(x.FileSet, x.file, s)
	}
}

type NamedImport struct {
	Name    string
	Package string
}

func (x Xipe) AddNamedImports(namedImports ...NamedImport) {
	for _, s := range namedImports {
		astutil.AddNamedImport(x.FileSet, x.file, s.Name, s.Package)
	}
}

func (x Xipe) AddDecls(dscls ...ast.Decl) {
	x.file.Decls = append(x.file.Decls, dscls...)
}

func (x Xipe) Write(w io.Writer) error {
	return printer.Fprint(w, x.FileSet, x.file)
}
