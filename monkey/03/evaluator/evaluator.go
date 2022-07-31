package evaluator

import (
	"github.com/kuromitsu0104/go-interpreter/monkey/03/ast"
	"github.com/kuromitsu0104/go-interpreter/monkey/03/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}
