package evaluator

import (
	"fmt"

	"github.com/kuroweb/go-interpreter/practice/monkey/03/ast"
	"github.com/kuroweb/go-interpreter/practice/monkey/03/object"
)

func Eval(node ast.Node) object.Object {
	// データ型に合わせて処理を分岐
	switch node := node.(type) {
	case *ast.Program:
		// Program型のときに実行. Programが持つStatements単位での評価を evalStatements関数 で実行
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		// ReturnStatement型のときは、object.ReturnValueでラップした値を返却
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
	}

	return nil
}

var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpresison(right)
	default:
		// エラーオブジェクトを返却
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpresison(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		// エラーオブジェクトを返却
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		// エラーオブジェクトを返却
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		// エラーオブジェクトを返却
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		// エラーオブジェクトを返却
		return newError("unknown oeprator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// Program型を評価するための関数.
func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	// Programが持つStatements単位でStatementの評価を行う.
	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			// object.ReturnValueが返却されたとき、後続のStatementsの評価をスキップする.
			return result.Value
		case *object.Error:
			// object.Errorが返却されたとき、後続のStatementsの評価をスキップする.
			return result
		}
	}

	return result
}

// BlockStatement型を評価するための関数.
// NOTE: Program型と分離して、Blockの評価に特化させたもの.
func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	// BlockStatementが持つStatements単位でStatementの評価を行う.
	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			// ブロック内に return が存在していても後続の処理を継続する.
			// (evalProgramでは、トップレベルでobject.ReturnValueが返却されると後続のStatementsの評価をスキップする)
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// エラーオブジェクトを返却する関数
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// 評価結果がエラーであるのかチェックする関数
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
