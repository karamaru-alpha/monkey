package evaluator

import (
	"github.com/karamaru-alpha/monkey/ast"
	"github.com/karamaru-alpha/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	if operator == "!" {
		return evalBangOperatorExpression(right)
	}
	if operator == "-" {
		return evalMinusOperatorExpression(right)
	}
	return NULL
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
		return NULL
	}
	return &object.Integer{Value: -rightObj.Value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
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
	}
	return NULL
}
