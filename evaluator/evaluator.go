package evaluator

import (
	"fmt"
	"github.com/user/courseWork/ast"
	"github.com/user/courseWork/object"
)

var (
	NIL   = &object.NIL{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.DoubleLiteral:
		return &object.Double{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
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
	case *ast.ForLoopExpression:
		return evalForLoopxpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if _, ok := env.Get(node.Name.String()); !ok {
			return newError("identifier not found: %s", node.Name.String())
		}
		env.Set(node.Name.String(), val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.ArrayLiteral:
		el := evalExpressions(node.Elements, env)
		if len(el) == 1 && isError(el[0]) {
			return el[0]
		}
		return &object.Array{Elements: el}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var res object.Object

	for _, statement := range stmts {
		res = Eval(statement, env)

		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}

	}

	return res
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var res object.Object

	for _, statement := range block.Statements {
		res = Eval(statement, env)

		if res != nil {
			if isReturn(res) || isError(res) { //return can propagate further up
				return res
			}
		}
	}

	return res
}

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
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case isNumeric(left) && isNumeric(right):
		return evalNumericInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NIL
	}

	return arrayObject.Elements[idx]
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if !isNumeric(right) {
		return newError("unknown operator: -%s", right.Type())
	}

	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	} else {
		value := right.(*object.Double).Value
		return &object.Double{Value: -value}
	}
}

func evalNumericInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	} else {
		if left.Type() == object.INTEGER_OBJ {
			value := left.(*object.Integer).Value
			left = &object.Double{float64(value)}
		}
		if right.Type() == object.INTEGER_OBJ {
			value := right.(*object.Integer).Value
			right = &object.Double{float64(value)}
		}
		return evalDoubleInfixExpression(operator, left, right)
	}
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
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalDoubleInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.Double).Value

	switch operator {
	case "+":
		return &object.Double{Value: leftVal + rightVal}
	case "-":
		return &object.Double{Value: leftVal - rightVal}
	case "*":
		return &object.Double{Value: leftVal * rightVal}
	case "/":
		return &object.Double{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value
	switch operator {
	case "&&":
		if leftVal && rightVal {
			return TRUE
		} else {
			return FALSE
		}
	case "||":
		if leftVal || rightVal {
			return TRUE
		} else {
			return FALSE
		}
	case "==":
		if leftVal == rightVal {
			return TRUE
		} else {
			return FALSE
		}
	case "!=":
		if leftVal != rightVal {
			return TRUE
		} else {
			return FALSE
		}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTrue(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NIL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func isTrue(obj object.Object) bool {
	switch obj {
	case NIL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isNumeric(obj object.Object) bool {
	return obj.Type() == object.INTEGER_OBJ || obj.Type() == object.DOUBLE_OBJ
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isReturn(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.RETURN_VALUE_OBJ
	}
	return false
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendedFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrawpReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrawpReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalForLoopxpression(fle *ast.ForLoopExpression, env *object.Environment) object.Object {
	if rtn, ok := fle.Init.(*ast.ReturnStatement); ok {
		return newError("for loop initial declaration is invalid: %s", rtn.String())
	}

	var (
		obj         object.Object = NIL
		extendedEnv *object.Environment
	)
	if _, ok := fle.Init.(*ast.VarStatement); ok {
		extendedEnv = object.NewEnclosedEnvironment(env)
	} else {
		extendedEnv = env
	}

	init := Eval(fle.Init, extendedEnv)
	if isError(init) {
		return init
	}
	condition := Eval(fle.Condition, extendedEnv)
	if isError(condition) {
		return condition
	}
	for isTrue(condition) {
		obj = Eval(fle.Body, extendedEnv)
		if isError(obj) || isReturn(obj) {
			return obj
		}
		step := Eval(fle.Step, extendedEnv)
		if isError(step) {
			return step
		}
		condition = Eval(fle.Condition, extendedEnv)
		if isError(condition) {
			return condition
		}
	}
	return obj
}
