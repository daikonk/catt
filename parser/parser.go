package parser

import (
	"fmt"
	"go_interpreter/ast"
	"go_interpreter/lexer"
	"go_interpreter/token"
	"strconv"
)

type (
	prefixParserFn func() ast.Expression
	infixParserFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.MINUS:  SUM,
	token.PLUS:   SUM,
	token.ASTR:   PRODUCT,
	token.SLASH:  PRODUCT,
	token.MODULO: PRODUCT,
	token.LPAR:   CALL,
}

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParserFn) {
	p.prefixParserFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParserFn) {
	p.infixParserFns[tokenType] = fn
}

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	aftToken  token.Token
	errors    []string

	prefixParserFns map[token.TokenType]prefixParserFn
	infixParserFns  map[token.TokenType]infixParserFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.NextToken()
	p.NextToken()

	p.prefixParserFns = make(map[token.TokenType]prefixParserFn)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.ASSIGN, p.parsePrefixExpression)
	p.registerPrefixFn(token.FALSE, p.parseBoolean)
	p.registerPrefixFn(token.TRUE, p.parseBoolean)
	p.registerPrefixFn(token.LPAR, p.parseParExpression)
	p.registerPrefixFn(token.IF, p.parseIfExpression)
	p.registerPrefixFn(token.STRING, p.ParseString)
	p.registerPrefixFn(token.WHILE, p.ParseWhileExpression)
	p.registerPrefixFn(token.FUNCTION, p.ParseFunctionLiteral)

	p.infixParserFns = make(map[token.TokenType]infixParserFn)
	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTR, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.MODULO, p.parseInfixExpression)
	p.registerInfixFn(token.LPAR, p.ParseCallExpression)

	return p
}

func (p *Parser) ParseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: p.currToken, Function: function}
	expr.Arguments = p.ParseCallArguments()

	return expr
}

func (p *Parser) ParseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.aftToken.Type == token.RPAR {
		p.NextToken()
		return args
	}

	p.NextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.aftToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.PeekAndMove(token.RPAR) {
		return nil
	}

	return args
}

func (p *Parser) ParseFunctionLiteral() ast.Expression {
	fnc := &ast.FunctionLiteral{Token: p.currToken}

	if !p.PeekAndMove(token.LPAR) {
		return nil
	}

	fnc.Parameters = p.ParseFunctionParameters()

	if !p.PeekAndMove(token.LBRAC) {
		return nil
	}

	fnc.Body = p.parseBlockStatement()

	return fnc
}

func (p *Parser) ParseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.aftToken.Type == token.RPAR {
		p.NextToken()
		return identifiers
	}

	p.NextToken()

	identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, identifier)

	for p.aftToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, identifier)

	}

	if !p.PeekAndMove(token.RPAR) {
		return nil
	}

	return identifiers
}

func (p *Parser) ParseString() ast.Expression {
	return &ast.String{p.currToken, p.currToken.Literal}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}

	p.NextToken()

	for (p.currToken.Type != token.RBRAC) && (p.currToken.Type != token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}

	return block
}

func (p *Parser) ParseWhileExpression() ast.Expression {
	expr := &ast.WhileExpression{Token: p.currToken}

	if !p.PeekAndMove(token.LPAR) {
		return nil
	}

	p.NextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.PeekAndMove(token.RPAR) {
		return nil
	}

	if !p.PeekAndMove(token.LBRAC) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currToken}

	if !p.PeekAndMove(token.LPAR) {
		return nil
	}

	p.NextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.PeekAndMove(token.RPAR) {
		return nil
	}

	if !p.PeekAndMove(token.LBRAC) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.aftToken.Type == token.ELSE {

		p.NextToken()

		if !p.PeekAndMove(token.LBRAC) {
			return nil
		}

		expr.Alternative = p.parseBlockStatement()

	}

	return expr
}

func (p *Parser) parseParExpression() ast.Expression {
	p.NextToken()
	expr := p.parseExpression(LOWEST)

	if !p.PeekAndMove(token.RPAR) {
		return nil
	}

	return expr
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: (p.currToken.Type == token.TRUE)}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	precedence := p.CurPrecedence()
	p.NextToken()

	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.NextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected %s as aftToken, got %s instead", t, p.aftToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function found for %s", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) NextToken() {
	p.currToken = p.aftToken
	p.aftToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.NextToken()
	}

	return program
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{p.currToken, p.currToken.Literal}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.VAR:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.NextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.aftToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.aftToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParserFns[p.currToken.Type]

	if prefix == nil {

		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefix()

	for (p.aftToken.Type != token.SEMICOLON) && (precedence < p.PeekPrecedence()) {
		infix := p.infixParserFns[p.aftToken.Type]
		if infix == nil {
			return leftExp
		}

		p.NextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}
	if !p.PeekAndMove(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.PeekAndMove(token.ASSIGN) {
		return nil
	}

	p.NextToken()

	stmt.Value = p.parseExpression(LOWEST)

	for p.currToken.Type != token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as an int64", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) PeekAndMove(tkn token.TokenType) bool {
	if p.aftToken.Type == tkn {
		p.NextToken()
		return true
	} else {
		p.peekError(tkn)
		return false
	}
}

func (p *Parser) PeekPrecedence() int {
	if p, ok := precedences[p.aftToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) CurPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}
