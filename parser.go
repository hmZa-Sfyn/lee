// parser.go
package main

import (
	"strconv"
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(input string) *Parser {
	p := &Parser{l: NewLexer(input)}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() (*Program, *EsoError) {
	prog := &Program{Functions: make(map[string]*FunctionDecl), TopStmts: []Stmt{}}
	for p.curToken.Type != EOF {
		if decl, err := p.parseFunctionDecl(); err == nil {
			prog.Functions[decl.Name] = decl
		} else if stmt, err := p.parseStmt(); err == nil {
			prog.TopStmts = append(prog.TopStmts, stmt)
		} else {
			return nil, err
		}
		p.nextToken()
	}
	return prog, nil
}

func (p *Parser) parseFunctionDecl() (*FunctionDecl, *EsoError) {
	if p.curToken.Type != Ident {
		return nil, NewEsoErrorf(p.curToken.Line, p.curToken.Col, "expected type for function return")
	}
	returnType := p.curToken.Value
	if p.peekToken.Type != Colon {
		return nil, NewEsoErrorf(p.peekToken.Line, p.peekToken.Col, "expected : after return type")
	}
	p.nextToken() // :
	p.nextToken()
	if p.curToken.Type != Ident {
		return nil, NewEsoErrorf(p.curToken.Line, p.curToken.Col, "expected function name")
	}
	name := p.curToken.Value
	if p.peekToken.Type != Eq {
		return nil, NewEsoErrorf(p.peekToken.Line, p.peekToken.Col, "expected = after function name")
	}
	p.nextToken() // =

	var params []Param
	p.nextToken()
	for p.curToken.Type != Arrow {
		if p.curToken.Type == EOF {
			return nil, NewEsoErrorf(p.curToken.Line, p.curToken.Col, "unexpected EOF in params")
		}
		paramType := p.curToken.Value
		if p.peekToken.Type != Colon {
			return nil, NewEsoErrorf(p.peekToken.Line, p.peekToken.Col, "expected : after param type")
		}
		p.nextToken() // :
		p.nextToken()
		paramName := p.curToken.Value
		params = append(params, Param{Type: paramType, Name: paramName})
		if p.peekToken.Type == Pipe {
			p.nextToken() // |
			p.nextToken()
		}
	}

	p.nextToken() // ->
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &FunctionDecl{Name: name, Params: params, Return: returnType, Body: body}, nil
}

func (p *Parser) parseStmt() (Stmt, *EsoError) {
	switch p.curToken.Type {
	case Let:
		return p.parseLetStmt()
	default:
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ExprStmt{Expr: expr}, nil
	}
}

func (p *Parser) parseLetStmt() (*LetStmt, *EsoError) {
	mutable := false
	p.nextToken()
	if p.curToken.Type == Mut {
		mutable = true
		p.nextToken()
	}
	if p.curToken.Type != Ident {
		return nil, NewEsoErrorf(p.curToken.Line, p.curToken.Col, "expected type after let/mut")
	}
	typ := p.curToken.Value
	if p.peekToken.Type != Colon {
		return nil, NewEsoErrorf(p.peekToken.Line, p.peekToken.Col, "expected : after type")
	}
	p.nextToken() // :
	p.nextToken()
	if p.curToken.Type != Ident {
		return nil, NewEsoErrorf(p.curToken.Line, p.curToken.Col, "expected name after :")
	}
	name := p.curToken.Value
	if p.peekToken.Type != Eq {
		return nil, NewEsoErrorf(p.peekToken.Line, p.peekToken.Col, "expected = after name")
	}
	p.nextToken() // =
	p.nextToken()
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &LetStmt{Mutable: mutable, Type: typ, Name: name, Value: value}, nil
}

func (p *Parser) parseExpr() (Expr, *EsoError) {
	return p.parseBinaryExpr(0)
}

var precedences = map[TokenType]int{
	EqEq:      1,
	NotEq:     1,
	Less:      2,
	Greater:   2,
	LessEq:    2,
	GreaterEq: 2,
	Plus:      3,
	Minus:     3,
	Star:      4,
	Slash:     4,
	Percent:   4,
}

func (p *Parser) parseBinaryExpr(prec int) (Expr, *EsoError) {
	left, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	for {
		opPrec := precedences[p.peekToken.Type]
		if opPrec <= prec {
			break
		}
		op := p.peekToken.Type
		p.nextToken()
		p.nextToken()
		right, err := p.parseBinaryExpr(opPrec)
		if err != nil {
			return nil, err
		}
		left = &BinOpExpr{Left: left, Op: op, Right: right}
	}
	return left, nil
}

func (p *Parser) parsePrimaryExpr() (Expr, *EsoError) {
	switch p.curToken.Type {
	case Ident:
		if p.peekToken.Type == LParen || p.peekToken.Type == Pipe || p.peekToken.Type == Ident { // rough check for call
			return p.parseCallExpr()
		}
		return &IdentExpr{Value: p.curToken.Value}, nil
	case Int:
		val, err := strconv.ParseInt(p.curToken.Value, 10, 64)
		if err != nil {
			return nil, NewEsoErrorf(p.curToken.Line, p.curToken.Col, "invalid int: %s", p.curToken.Value)
		}
		return &IntExpr{Value: val}, nil
	case True:
		return &BoolExpr{Value: true}, nil
	case False:
		return &BoolExpr{Value: false}, nil
	case String:
		return &StringExpr{Value: p.curToken.Value}, nil
	case If:
		return p.parseIfExpr()
	case While:
		return p.parseWhileExpr()
	case LParen:
		p.nextToken()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.peekToken.Type != RParen {
			return nil, NewEsoErrorf(p.peekToken.Line, p.peekToken.Col, "expected )")
		}
		p.nextToken()
		return expr, nil
	// Add LBracket for vec, etc.
	default:
		return nil, NewEsoErrorf(p.curToken.Line, p.curToken.Col, "unexpected token %v", p.curToken.Type)
	}
}

func (p *Parser) parseCallExpr() (*CallExpr, *EsoError) {
	callee := &IdentExpr{Value: p.curToken.Value}
	p.nextToken()
	var args []Expr
	for p.curToken.Type != EOF && p.curToken.Type != Arrow && p.curToken.Type != Else {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.peekToken.Type == Pipe {
			p.nextToken()
			p.nextToken()
		} else {
			break
		}
	}
	return &CallExpr{Callee: callee, Args: args}, nil
}

func (p *Parser) parseIfExpr() (*IfExpr, *EsoError) {
	p.nextToken()
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type != Arrow {
		return nil, NewEsoErrorf(p.peekToken.Line, p.peekToken.Col, "expected -> after if cond")
	}
	p.nextToken() // ->
	p.nextToken()
	then, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	var els Expr
	if p.peekToken.Type == Else {
		p.nextToken() // else
		p.nextToken()
		els, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	return &IfExpr{Cond: cond, Then: then, Else: els}, nil
}

func (p *Parser) parseWhileExpr() (*WhileExpr, *EsoError) {
	p.nextToken()
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type != Arrow {
		return nil, NewEsoErrorf(p.peekToken.Line, p.peekToken.Col, "expected -> after while cond")
	}
	p.nextToken() // ->
	p.nextToken()
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &WhileExpr{Cond: cond, Body: body}, nil
}
