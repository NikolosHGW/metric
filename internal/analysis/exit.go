package analysis

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	funcName    = "main"
	packageName = "main"
)

var Analyzer = &analysis.Analyzer{
	Name: "exit",
	Doc:  "Анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != packageName {
			continue
		}
		filename := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(filename, ".go") {
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			fun, ok := n.(*ast.FuncDecl)
			if ok {
				if fun.Name.Name == funcName {
					ast.Inspect(fun.Body, func(n ast.Node) bool {
						if call, ok := n.(*ast.CallExpr); ok {
							if isOsExitCall(call) {
								pass.Reportf(call.Pos(), "вызов Exit функции пакета os не рекомендуется")
							}
						}

						return true
					})
				}
			}

			return true
		})
	}

	return nil, nil
}

func isOsExitCall(call *ast.CallExpr) bool {
	if selectorIdent, ok := call.Fun.(*ast.SelectorExpr); ok {
		if parentIdent, ok := selectorIdent.X.(*ast.Ident); ok {
			if parentIdent.Name == "os" && selectorIdent.Sel.Name == "Exit" {
				return true
			}
		}
	}

	return false
}
