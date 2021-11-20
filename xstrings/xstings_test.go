package xstrings_test

import (
	"go/ast"
	"testing"

	"github.com/nikgalushko/xipe/xstrings"

	"github.com/stretchr/testify/require"
)

func TestToStmts(t *testing.T) {
	const s = `
	db, err := postgres.NewDB(cfg.Postgres, srv.Prometheus)
	if err != nil {
		panic(err)
	}
	`

	actual, err := xstrings.ToStmts(s)
	require.NoError(t, err)

	require.NotNil(t, actual)
	require.Len(t, actual, 2)
	require.IsType(t, &ast.AssignStmt{}, actual[0])
	require.IsType(t, &ast.IfStmt{}, actual[1])
}

func TestToStmts_InvalidInput(t *testing.T) {
	actual, err := xstrings.ToStmts("blah+{kek}")
	require.ErrorIs(t, err, xstrings.ErrIncorrectStmts)
	require.Nil(t, actual)
}

func TestToDecl(t *testing.T) {
	const s = `
	type MyType struct {
		Hosts []string
		Timeout time.Duration
	}

	var i = 0
	`
	actual, err := xstrings.ToDecl(s)
	require.NoError(t, err)

	require.NotNil(t, actual)
	require.Len(t, actual, 2)

	for i := 0; i < len(actual); i++ {
		require.IsType(t, &ast.GenDecl{}, actual[i])
	}

	require.Len(t, actual[0].(*ast.GenDecl).Specs, 1)
	require.IsType(t, &ast.TypeSpec{}, actual[0].(*ast.GenDecl).Specs[0])
	require.IsType(t, &ast.ValueSpec{}, actual[1].(*ast.GenDecl).Specs[0])
}

func TestToDecl_InvalidInput(t *testing.T) {
	actual, err := xstrings.ToDecl("var d = ∂")
	require.ErrorIs(t, err, xstrings.ErrIncorrectDecl)
	require.Nil(t, actual)
}

func TestToExpr(t *testing.T) {
	const s = `strings.Split(str, " ")[0] + "blah"`
	actual, err := xstrings.ToExpr(s)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.IsType(t, &ast.BinaryExpr{}, actual)
	require.Equal(t, s, xstrings.FromNode(actual))
}

func TestToExp_InvalidInput(t *testing.T) {
	actual, err := xstrings.ToExpr("~(d >> i)")
	require.ErrorIs(t, err, xstrings.ErrIncorrectExpr)
	require.Nil(t, actual)
}
