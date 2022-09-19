package evaluator

import (
	"fmt"

	"github.com/karamaru-alpha/monkey/ast"
	"github.com/karamaru-alpha/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return toBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch res := result.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}
	return result
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil {
			typ := result.Type()
			if typ == object.RETURN_VALUE || typ == object.ERROR {
				return result
			}
		}
	}
	return result
}

func toBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(ident.Value)
	if !ok {
		return newError("identifier not found: %s", ident.Value)
	}
	return val
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	if operator == "!" {
		return evalBangOperatorExpression(right)
	}
	if operator == "-" {
		return evalMinusOperatorExpression(right)
	}
	return newError("unknown operator: %s%s", operator, right.Type())
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE, NULL:
		return TRUE
	}
	return FALSE
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	rightObj, ok := right.(*object.Integer)
	if !ok {
		return newError("unknown operator: -%s", right.Type())
	}
	return &object.Integer{Value: -rightObj.Value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return evalIntegerInfixExpression(operator, left, right)
	}
	if operator == "==" {
		return toBooleanObject(left == right)
	}
	if operator == "!=" {
		return toBooleanObject(left != right)
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
	case ">":
		return toBooleanObject(leftVal > rightVal)
	case "<":
		return toBooleanObject(leftVal < rightVal)
	case "==":
		return toBooleanObject(leftVal == rightVal)
	case "!=":
		return toBooleanObject(leftVal != rightVal)
	}
	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalIfExpression(exp *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(exp.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(exp.Consequence, env)
	}
	if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
	}
	return NULL
}

func isTruthy(obj object.Object) bool {
	if obj == FALSE || obj == NULL {
		return false
	}
	return true

}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}
