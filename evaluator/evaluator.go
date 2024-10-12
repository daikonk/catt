package evaluator

import (
	"fmt"
	"go_interpreter/ast"
	"go_interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

var builtins = map[string]*object.BuiltIn{
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return nil
			}

			switch arg := args[0].(type) {
			case *object.String:
				fmt.Printf(arg.Value)
				return &object.String{Value: arg.Value}

			case *object.Integer:
				fmt.Printf(arg.Inspect())
				return &object.String{Value: arg.Inspect()}

			default:
				return nil
			}
		},
	},
	"putsln": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return nil
			}

			switch arg := args[0].(type) {
			case *object.String:
				fmt.Println(arg.Value)
				return &object.String{Value: arg.Value}

			case *object.Integer:
				fmt.Println(arg.Inspect())
				return &object.String{Value: arg.Inspect()}

			default:
				return nil
			}
		},
	},
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.String:
		return &object.String{node.Value}
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
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		left := Eval(node.Left, env)
		return evalInfixExpression(node.Operator, right, left)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if val == NULL {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		return &object.ReturnValue{Value: val}

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if function == nil {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && args[0] == nil {
			return args[0]
		}
		return applyFunction(function, args)
	}

	return nil
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
		return nil
	}
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
		if evaluated != nil {
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

		if result != nil && result.Type() == object.RETURN_VAL_OBJ {
			return result
		}
	}

	return result
}

func evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)

		if value, ok := result.(*object.ReturnValue); ok {
			return value.Value
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

	return nil
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	evaluation := Eval(ie.Condition, env)

	if isTruthy(evaluation) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	for eval := Eval(we.Condition, env); isTruthy(eval); eval = Eval(we.Condition, env) {
		Eval(we.Consequence, env)
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
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(op, right, left)
	case op == "==":
		return nativeToObjectBool(right == left)
	case op == "!=":
		return nativeToObjectBool(right != left)
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
		return NULL
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
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
