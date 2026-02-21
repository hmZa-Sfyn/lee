// interp.go
package main

import "fmt"

func Interpret(prog *Program, env *Environment) (Value, *EsoError) {
	for name, fn := range prog.Functions {
		env.DefineFn(name, fn)
	}

	var last Value
	for _, stmt := range prog.TopStmts {
		val, err := evalStmt(stmt, env)
		if err != nil {
			return nil, err
		}
		if val != nil {
			last = val
		}
	}

	if mainFn, ok := env.GetFn("main"); ok {
		val, err := callFunction(mainFn, nil, env)
		if err != nil {
			return nil, err
		}
		return val, nil
	}

	return last, nil
}

func evalStmt(node Stmt, env *Environment) (Value, *EsoError) {
	switch s := node.(type) {
	case *LetStmt:
		val, err := evalExpr(s.Value, env)
		if err != nil {
			return nil, err
		}
		return nil, env.Define(s.Name, val, s.Mutable)
	case *AssignStmt:
		val, err := evalExpr(s.Value, env)
		if err != nil {
			return nil, err
		}
		return nil, env.Assign(s.Name, val)
	case *ExprStmt:
		return evalExpr(s.Expr, env)
	default:
		return nil, NewEsoErrorf(0, 0, "unsupported stmt %T", node)
	}
}

func evalExpr(node Expr, env *Environment) (Value, *EsoError) {
	switch e := node.(type) {
	case *IdentExpr:
		val, ok := env.Get(e.Value)
		if !ok {
			return nil, NewEsoErrorf(0, 0, "undefined %s", e.Value)
		}
		return val, nil
	case *IntExpr:
		return IntVal{Val: e.Value}, nil
	case *BoolExpr:
		return BoolVal{Val: e.Value}, nil
	case *StringExpr:
		return StringVal{Val: e.Value}, nil
	case *BinOpExpr:
		left, err := evalExpr(e.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := evalExpr(e.Right, env)
		if err != nil {
			return nil, err
		}
		return evalBinOp(left, e.Op, right)
	case *CallExpr:
		ident, ok := e.Callee.(*IdentExpr)
		if !ok {
			return nil, NewEsoErrorf(0, 0, "callee must be ident")
		}
		fn, ok := env.GetFn(ident.Value)
		if ok {
			return callFunction(fn, e.Args, env)
		}
		return callBuiltin(ident.Value, e.Args, env)
	case *IfExpr:
		cond, err := evalExpr(e.Cond, env)
		if err != nil {
			return nil, err
		}
		b, ok := cond.(BoolVal)
		if !ok {
			return nil, NewEsoErrorf(0, 0, "if cond not bool")
		}
		if b.Val {
			return evalExpr(e.Then, env)
		} else if e.Else != nil {
			return evalExpr(e.Else, env)
		}
		return VoidVal{}, nil
	case *WhileExpr:
		for {
			cond, err := evalExpr(e.Cond, env)
			if err != nil {
				return nil, err
			}
			b, ok := cond.(BoolVal)
			if !ok {
				return nil, NewEsoErrorf(0, 0, "while cond not bool")
			}
			if !b.Val {
				break
			}
			_, err = evalExpr(e.Body, env)
			if err != nil {
				return nil, err
			}
		}
		return VoidVal{}, nil
	default:
		return nil, NewEsoErrorf(0, 0, "unsupported expr %T", node)
	}
}

func evalBinOp(left Value, op TokenType, right Value) (Value, *EsoError) {
	l, lok := left.(IntVal)
	r, rok := right.(IntVal)
	if !lok || !rok {
		return nil, NewEsoErrorf(0, 0, "binop on non-int")
	}
	switch op {
	case Plus:
		return IntVal{Val: l.Val + r.Val}, nil
	case Minus:
		return IntVal{Val: l.Val - r.Val}, nil
	case Star:
		return IntVal{Val: l.Val * r.Val}, nil
	case Slash:
		if r.Val == 0 {
			return nil, NewEsoError("division by zero", 0, 0)
		}
		return IntVal{Val: l.Val / r.Val}, nil
	case EqEq:
		return BoolVal{Val: l.Val == r.Val}, nil
	case Less:
		return BoolVal{Val: l.Val < r.Val}, nil
	// add more
	default:
		return nil, NewEsoErrorf(0, 0, "unsupported op %v", op)
	}
}

func callFunction(fn *FunctionDecl, argExprs []Expr, env *Environment) (Value, *EsoError) {
	if len(argExprs) != len(fn.Params) {
		return nil, NewEsoErrorf(0, 0, "arity mismatch: %d vs %d", len(fn.Params), len(argExprs))
	}
	local := NewEnvironment()
	local.outer = env
	for i, param := range fn.Params {
		val, err := evalExpr(argExprs[i], env)
		if err != nil {
			return nil, err
		}
		local.Define(param.Name, val, false) // params immutable
	}
	return evalExpr(fn.Body, local)
}

func callBuiltin(name string, argExprs []Expr, env *Environment) (Value, *EsoError) {
	if name == "print" {
		for _, arg := range argExprs {
			val, err := evalExpr(arg, env)
			if err != nil {
				return nil, err
			}
			fmt.Println(val.String())
		}
		return VoidVal{}, nil
	}
	// add ESO built-ins
	return nil, NewEsoErrorf(0, 0, "unknown builtin %s", name)
}