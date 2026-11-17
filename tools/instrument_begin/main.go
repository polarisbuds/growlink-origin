//go:build ignore

// One-off: go run tools/instrument_begin/main.go <dir> [usecase|repository]
package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	dir := os.Args[1]
	layer := "usecase"
	if len(os.Args) > 2 {
		layer = os.Args[2]
	}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		return instrument(path, layer)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func instrument(path, layer string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	opPrefix := packageOpPrefix(f.Name.Name)
	changed := false
	needCtxImport := false

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv == nil || fd.Body == nil {
			continue
		}
		if !fd.Name.IsExported() || strings.HasPrefix(fd.Name.Name, "New") {
			continue
		}
		if bodyHasBegin(fd.Body) {
			continue
		}
		ctxExpr := ctxExpression(fd, layer, &needCtxImport)
		op := opForMethod(f, fd, opPrefix, layer)
		returnsErr := lastResultIsError(fd.Type)
		if returnsErr {
			ensureNamedErr(fd.Type)
		}
		prefix := beginStmts(ctxExpr, layer, op, returnsErr)
		fd.Body.List = append(prefix, fd.Body.List...)
		if returnsErr {
			fixErrShortDeclInFunc(fd)
		}
		changed = true
	}

	if !changed {
		return nil
	}
	if needCtxImport {
		astutil.AddImport(fset, f, "context")
	}
	astutil.AddNamedImport(fset, f, "logutils", "growlink/internal/shared/logUtils")

	var buf strings.Builder
	if err := format.Node(&buf, fset, f); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return os.WriteFile(path, []byte(buf.String()), 0o644)
}

func packageOpPrefix(pkg string) string {
	if pkg == "authenticate_usecase" {
		return "auth_usecase"
	}
	return pkg
}

func opForMethod(f *ast.File, fd *ast.FuncDecl, prefix, layer string) string {
	if layer == "repository" {
		recvType := receiverTypeName(fd.Recv)
		if recvType == "Commenter" {
			return "commenter_repository." + fd.Name.Name
		}
		snake := typeToSnakeRepo(recvType)
		return snake + "." + fd.Name.Name
	}
	return prefix + "." + fd.Name.Name
}

func receiverTypeName(fl *ast.FieldList) string {
	if fl == nil || len(fl.List) == 0 {
		return ""
	}
	t := fl.List[0].Type
	switch x := t.(type) {
	case *ast.StarExpr:
		return typeIdent(x.X)
	case *ast.IndexExpr:
		return typeIdent(x.X)
	default:
		return typeIdent(t)
	}
}

func typeIdent(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	default:
		return ""
	}
}

func typeToSnakeRepo(typeName string) string {
	typeName = strings.TrimSuffix(typeName, "Repository")
	var parts []string
	runes := []rune(typeName)
	start := 0
	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) {
			parts = append(parts, strings.ToLower(string(runes[start:i])))
			start = i
		}
	}
	parts = append(parts, strings.ToLower(string(runes[start:])))
	return strings.Join(parts, "_") + "_repository"
}

func ctxExpression(fd *ast.FuncDecl, layer string, needCtx *bool) string {
	if layer == "repository" {
		rv := "repo"
		if len(fd.Recv.List[0].Names) > 0 && fd.Recv.List[0].Names[0].Name != "" {
			rv = fd.Recv.List[0].Names[0].Name
		}
		return rv + ".ctx"
	}
	if fd.Type.Params != nil {
		for _, field := range fd.Type.Params.List {
			if sel, ok := field.Type.(*ast.SelectorExpr); ok {
				if pkg, _ := sel.X.(*ast.Ident); pkg != nil && pkg.Name == "context" && sel.Sel.Name == "Context" {
					if len(field.Names) > 0 {
						return field.Names[0].Name
					}
					return "ctx"
				}
			}
		}
	}
	*needCtx = true
	return "context.Background()"
}

func lastResultIsError(ft *ast.FuncType) bool {
	if ft.Results == nil || len(ft.Results.List) == 0 {
		return false
	}
	last := ft.Results.List[len(ft.Results.List)-1].Type
	id, ok := last.(*ast.Ident)
	return ok && id.Name == "error"
}

func ensureNamedErr(ft *ast.FuncType) {
	if ft.Results == nil {
		return
	}
	// Name every result so go/printer accepts mixed qualified types with a named err.
	defaultNames := []string{"ret0", "ret1", "ret2", "ret3", "ret4"}
	ri := 0
	for _, field := range ft.Results.List {
		if len(field.Names) == 0 {
			name := "ret"
			if ri < len(defaultNames) {
				name = defaultNames[ri]
			} else {
				name = fmt.Sprintf("ret%d", ri)
			}
			// single anonymous type in field
			if len(field.Names) == 0 {
				field.Names = []*ast.Ident{ast.NewIdent(name)}
			}
			ri++
		}
	}
	last := ft.Results.List[len(ft.Results.List)-1]
	if id, ok := last.Type.(*ast.Ident); ok && id.Name == "error" {
		last.Names = []*ast.Ident{ast.NewIdent("err")}
	}
}

func beginStmts(ctxExpr, layer, op string, returnsErr bool) []ast.Stmt {
	layerIdent := "LayerUsecase"
	if layer == "repository" {
		layerIdent = "LayerRepository"
	}
	constOp := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.CONST,
			Specs: []ast.Spec{&ast.ValueSpec{
				Names:  []*ast.Ident{ast.NewIdent("op")},
				Values: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", op)}},
			}},
		},
	}
	done := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{ast.NewIdent("done")},
		Rhs: []ast.Expr{&ast.CallExpr{
			Fun: ast.NewIdent("logutils.Begin"),
			Args: []ast.Expr{
				ast.NewIdent(ctxExpr),
				&ast.SelectorExpr{X: ast.NewIdent("logutils"), Sel: ast.NewIdent(layerIdent)},
				ast.NewIdent("op"),
			},
		}},
	}
	doneArg := ast.NewIdent("nil")
	if returnsErr {
		doneArg = ast.NewIdent("err")
	}
	deferStmt := &ast.DeferStmt{Call: &ast.CallExpr{
		Fun: &ast.FuncLit{
			Type: &ast.FuncType{},
			Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{
				Fun: ast.NewIdent("done"), Args: []ast.Expr{doneArg},
			}}}},
		},
	}}
	return []ast.Stmt{constOp, done, deferStmt}
}

func bodyHasBegin(body *ast.BlockStmt) bool {
	for _, st := range body.List {
		as, ok := st.(*ast.AssignStmt)
		if !ok || len(as.Rhs) == 0 {
			continue
		}
		call, ok := as.Rhs[0].(*ast.CallExpr)
		if !ok {
			continue
		}
		if fun, ok := call.Fun.(*ast.Ident); ok && fun.Name == "logutils.Begin" {
			return true
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Begin" {
			return true
		}
	}
	return false
}

func fixErrShortDeclInFunc(fd *ast.FuncDecl) {
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.AssignStmt:
			if s.Tok != token.DEFINE {
				return true
			}
			if assignShouldUseEquals(s.Lhs) {
				s.Tok = token.ASSIGN
			}
		case *ast.IfStmt:
			if s.Init != nil {
				if as, ok := s.Init.(*ast.AssignStmt); ok && as.Tok == token.DEFINE && assignShouldUseEquals(as.Lhs) {
					as.Tok = token.ASSIGN
				}
			}
		}
		return true
	})
}

func assignShouldUseEquals(lhs []ast.Expr) bool {
	hasErr := false
	hasNewIdent := false
	for _, e := range lhs {
		id, ok := e.(*ast.Ident)
		if !ok {
			continue
		}
		if id.Name == "err" {
			hasErr = true
			continue
		}
		if id.Name != "_" {
			hasNewIdent = true
		}
	}
	if !hasErr {
		return false
	}
	if hasNewIdent {
		return false
	}
	return true
}