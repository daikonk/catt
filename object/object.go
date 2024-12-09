package object

import (
	"bytes"
	"fmt"
	"go_interpreter/ast"
	"strings"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ    = "INTEGER"
	BOOL_OBJ       = "BOOLEAN"
	NULL_OBJ       = "NULL"
	ERROR_OBJ      = "ERROR"
	STRING_OBJ     = "STRING_OBJ"
	RETURN_VAL_OBJ = "RETURN_VAL_OBJ"
	FUNCTION_OBJ   = "FUNCTION_OBJ"
	BUILTIN_OBJ    = "BUILTIN"
	CHANNEL_OBJ    = "CHANNEL"
	ROUTINE_OBJ    = "ROUTINE"
)

type Channel struct {
	Value chan Object
}

func (c *Channel) Type() ObjectType {
	return CHANNEL_OBJ
}

func (c *Channel) Inspect() string {
	return "channel"
}

type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInFunction
}

func (bi *BuiltIn) Type() ObjectType {
	return BUILTIN_OBJ
}

func (bi *BuiltIn) Inspect() string {
	return "builtin"
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VAL_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOL_OBJ
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
