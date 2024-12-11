package evaluator

import (
	"fmt"
	"go_interpreter/ast"
	"go_interpreter/object"
	"go_interpreter/utils"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

var builtins = map[string]*object.BuiltIn{
	"meow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("supports 1 argument, got: %d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				fmt.Print(arg.Value)
				return &object.String{Value: arg.Value}

			case *object.Integer:
				fmt.Print(arg.Inspect())
				return &object.String{Value: arg.Inspect()}

			case *object.Boolean:
				fmt.Print(arg.Inspect())
				return &object.String{Value: arg.Inspect()}

			case *object.Array:
				fmt.Print(arg.Inspect())
				return &object.String{Value: arg.Inspect()}

			default:
				return newError("argument type is not supported: %s", arg.Type())

			}
		},
	},
	"meowln": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("supports 1 argument, got: %d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				fmt.Println(arg.Value)
				return &object.String{Value: arg.Value}

			case *object.Integer:
				fmt.Println(arg.Inspect())
				return &object.String{Value: arg.Inspect()}

			case *object.Boolean:
				fmt.Println(arg.Inspect())
				return &object.String{Value: arg.Inspect()}

			case *object.Array:
				fmt.Print(arg.Inspect())
				return &object.String{Value: arg.Inspect()}

			default:
				return newError("argument type is not supported: %s", arg.Type())
			}
		},
	},
	"cattfusion": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("supports 1 argument, got: %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				res := utils.Cattfusion(arg.Value)
				return &object.String{Value: res}
			default:
				return newError("argument type is not supported: %s", arg.Type())
			}
		},
	},
	"cattify": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("supports 1 argument, got: %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				res, err := utils.Cattify(arg.Value)
				if res == "" {
					return newError(err)
				}
				return &object.String{Value: res}
			default:
				return newError("argument type is not supported: %s", arg.Type())

			}
		},
	},
	"cattsort": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("supports 1 argument, got: %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				res, err := utils.Cattsort(arg)
				if res == nil {
					if err == "" {
						return newError(".., too lazy to sort imma nap instead")
					} else {
						return newError(err)
					}
				}
				return res

			case *object.Integer:
				return &object.String{Value: arg.Inspect()}

			case *object.Boolean:
				return &object.String{Value: arg.Inspect()}

			default:
				return newError("argument type is not supported: %s", arg.Type())

			}
		},
	},
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeToObjectBool(node.Value)

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
		return evalInfixExpression(node.Operator, right, left)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

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
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

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
	}

	return nil
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		newEnvironment := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, newEnvironment)
		return unwrapReturnValue(evaluated)

	case *object.BuiltIn:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
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
		return newError("Index %d is null", idx)
	}

	return arrayObject.Elements[idx]
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for id, param := range fn.Parameters {
		env.Set(param.Value, args[id])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
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

func evalBlockStatement(node *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)

		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VAL_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
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

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	evaluation := Eval(ie.Condition, env)
	if isError(evaluation) {
		return evaluation
	}

	if isTruthy(evaluation) {
		return Eval(ie.Consequence, env)
	} else if ie.CondAlternative != nil {
		return Eval(ie.CondAlternative, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	condition := Eval(we.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition.Type() != object.BOOL_OBJ {
		return newError("expected type BOOL got: %s", condition.Type())
	}

	for isTruthy(condition) {
		blockstmt := Eval(we.Consequence, env)
		if isError(blockstmt) {
			return blockstmt
		}

		condition = Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if condition.Type() != object.BOOL_OBJ {
			return newError("expected type BOOL got: %s", condition.Type())
		}

	}

	return nil
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	decl := Eval(fe.Declaration, env)
	if isError(decl) {
		return decl
	}

	condition := Eval(fe.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition.Type() != object.BOOL_OBJ {
		return newError("expected type BOOL got: %s", condition.Type())
	}

	for isTruthy(condition) {

		blockstmt := Eval(fe.Consequence, env)
		if isError(blockstmt) {
			return blockstmt
		}

		incr := Eval(fe.Increment, env)
		if isError(incr) {
			return incr
		}

		condition = Eval(fe.Condition, env)
		if isError(condition) {
			return condition
		}

		if condition.Type() != object.BOOL_OBJ {
			return newError("expected type BOOL got: %s", condition.Type())
		}
	}

	return nil
}

func isTruthy(eval object.Object) bool {
	switch eval {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalInfixExpression(op string, right object.Object, left object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalInfixStringExpression(op, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(op, right, left)
	case op == "==":
		return nativeToObjectBool(right == left)
	case op == "!=":
		return nativeToObjectBool(right != left)
	case op == "&&":
		if left == TRUE && right == TRUE {
			return TRUE
		} else {
			return FALSE
		}
	case op == "||":
		if left == TRUE || right == TRUE {
			return TRUE
		} else {
			return FALSE
		}
	default:
		return newError("unknown infix operator: %s%s%s", left.Type(), op, right.Type())
	}
}

func evalInfixStringExpression(op string, right object.Object, left object.Object) object.Object {
	left_val := left.(*object.String).Value
	right_val := right.(*object.String).Value

	switch op {
	case "+":
		return &object.String{Value: left_val + right_val}
	default:
		return NULL
	}
}

func evalInfixIntegerExpression(op string, right object.Object, left object.Object) object.Object {
	left_val := left.(*object.Integer).Value
	right_val := right.(*object.Integer).Value

	switch op {
	case "%":
		return &object.Integer{Value: left_val % right_val}
	case "+":
		return &object.Integer{Value: left_val + right_val}
	case "-":
		return &object.Integer{Value: left_val - right_val}
	case "*":
		return &object.Integer{Value: left_val * right_val}
	case "/":
		return &object.Integer{Value: left_val / right_val}
	case ">":
		return nativeToObjectBool(left_val > right_val)
	case "<":
		return nativeToObjectBool(left_val < right_val)
	case "==":
		return nativeToObjectBool(left_val == right_val)
	case "!=":
		return nativeToObjectBool(left_val != right_val)
	default:
		return NULL
	}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown prefix operator: %s%s", op, right.Type())
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("not of type INT: %s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
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

func nativeToObjectBool(input bool) object.Object {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)
	}

	return result
}
