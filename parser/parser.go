package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/brandonksides/phonk/models"
	"github.com/brandonksides/phonk/tokens"
)

func ParseExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("expected token"),
		}
	}

	return parseAddExpression(toks)
}

type BindingExpression struct {
	Identifier string
	Expression models.Expression
}

type LetExpression struct {
	loc                models.SourceLocation
	BindingExpressions []BindingExpression
	Expression2        models.Expression
}

func (le *LetExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	newBindings := make(models.Bindings)
	for k, v := range bindings {
		newBindings[k] = v
	}

	for _, bindingExp := range le.BindingExpressions {
		k, v := bindingExp.Identifier, bindingExp.Expression
		val, err := v.Evaluate(newBindings)
		if err != nil {
			return nil, err
		}

		newBindings[k] = val

		if funcVal, ok := val.(*models.ExpFunction); ok {
			funcVal.Bindings[k] = val
		}
	}

	return le.Expression2.Evaluate(newBindings)
}

func (le *LetExpression) SourceLocation() models.SourceLocation {
	return le.loc
}

func parseLetExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) < 3 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.LET {
		return nil, toks, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]

	bindingExpressions := make([]BindingExpression, 0)
	for len(rest) > 0 {
		if rest[0].Type != tokens.IDENTIFIER {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}

		identifier := rest[0].Value
		rest = rest[1:]

		if rest[0].Type != tokens.EQUAL {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}

		rest = rest[1:]

		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("unexpected end of input"),
			}
		}

		var exp1 models.Expression
		exp1, rest, err = ParseExpression(rest[0:])
		if err != nil {
			return nil, rest, err
		}

		bindingExpressions = append(bindingExpressions, BindingExpression{
			Identifier: identifier,
			Expression: exp1,
		})

		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("unexpected end of input"),
			}
		}

		if rest[0].Type != tokens.COMMA {
			break
		}

		rest = rest[1:]
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type != tokens.IN {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	exp2, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &LetExpression{
		BindingExpressions: bindingExpressions,
		loc:                toks[0].SourceLocation,
		Expression2:        exp2,
	}, rest, nil
}

type IfExpression struct {
	Condition models.Expression
	Then      models.Expression
	Else      models.Expression
	loc       models.SourceLocation
}

func (ie *IfExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	cond, err := ie.Condition.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	condBool, ok := cond.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("if condition must evaluate to a boolean"),
			SourceLocation: ie.Condition.SourceLocation(),
		}
	}

	if condBool {
		return ie.Then.Evaluate(bindings)
	}

	return ie.Else.Evaluate(bindings)
}

func (ie *IfExpression) SourceLocation() models.SourceLocation {
	return ie.loc
}

func parseIfExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.IF {
		return nil, toks, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]
	exp1, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type != tokens.THEN {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	exp2, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type != tokens.ELSE {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	exp3, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &IfExpression{
		Condition: exp1,
		Then:      exp2,
		Else:      exp3,
		loc:       toks[0].SourceLocation,
	}, rest, nil
}

type ForExpression struct {
	Expression1 models.Expression
	Identifier  string
	Expression2 models.Expression
	loc         models.SourceLocation
}

func (fe *ForExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	ret := make([]any, 0)

	innerBindings := make(models.Bindings)
	for k, v := range bindings {
		innerBindings[k] = v
	}

	iterableExp, err := fe.Expression2.Evaluate(innerBindings)
	if err != nil {
		return nil, err
	}

	iterableExpArr, ok := iterableExp.([]any)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("for expression must evaluate to an array"),
			SourceLocation: fe.Expression2.SourceLocation(),
		}
	}

	for _, v := range iterableExpArr {
		innerBindings[fe.Identifier] = v
		retVal, err := fe.Expression1.Evaluate(innerBindings)
		if err != nil {
			return nil, err
		}

		ret = append(ret, retVal)
	}

	return ret, nil
}

func (fe *ForExpression) SourceLocation() models.SourceLocation {
	return fe.loc
}

func parseForExpression(exp1 models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return exp1, toks, nil
	}

	if toks[0].Type != tokens.FOR {
		return exp1, toks, nil
	}

	rest = toks[1:]
	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	identifier := rest[0].Value
	rest = rest[1:]
	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type != tokens.IN {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	exp2, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &ForExpression{
		Expression1: exp1,
		Identifier:  identifier,
		Expression2: exp2,
		loc:         toks[0].SourceLocation,
	}, rest, nil
}

type FunctionExpression struct {
	args []string
	exp  models.Expression
}

func (fe *FunctionExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	retBindings := make(models.Bindings)
	for k, v := range bindings {
		retBindings[k] = v
	}

	return &models.ExpFunction{
		Args:     fe.args,
		Bindings: retBindings,
		Exp:      fe.exp,
	}, nil
}

func (fe *FunctionExpression) SourceLocation() models.SourceLocation {
	return fe.exp.SourceLocation()
}

func parseFunction(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.FUNC {
		return nil, toks, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]
	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type != tokens.LEFT_PAREN {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	args := make([]string, 0)
	for len(rest) > 0 {
		if rest[0].Type != tokens.IDENTIFIER {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}

		args = append(args, rest[0].Value)
		rest = rest[1:]
		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("unexpected end of input"),
			}
		}

		if rest[0].Type == tokens.RIGHT_PAREN {
			break
		}

		if rest[0].Type != tokens.COMMA {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}

		rest = rest[1:]
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	rest = rest[1:]

	exp, rest, err = ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &FunctionExpression{
		args: args,
		exp:  exp,
	}, rest, nil
}

type AddExpression struct {
	op     tokens.Token
	first  models.Expression
	second models.Expression
}

func (ae *AddExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := ae.first.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Adder, ok := v1.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            fmt.Errorf("operator '%s' cannot be applied to first operand", ae.op.Value),
			SourceLocation: ae.first.SourceLocation(),
		}
	}

	v2, err := ae.second.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v2Adder, ok := v2.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            fmt.Errorf("operator '%s' cannot be applied to second operand", ae.op.Value),
			SourceLocation: ae.second.SourceLocation(),
		}
	}

	switch ae.op.Type {
	case tokens.PLUS:
		return v1Adder + v2Adder, nil
	case tokens.MINUS:
		return v1Adder - v2Adder, nil
	default:
		return nil, &models.InterpreterError{
			Err:            errors.New("invalid operator"),
			SourceLocation: ae.op.SourceLocation,
		}
	}
}

func (ae *AddExpression) SourceLocation() models.SourceLocation {
	return ae.first.SourceLocation()
}

func parseAddExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	exp, rest, err = parseMulExpression(toks)
	if err != nil {
		return nil, rest, err
	}

	return foldAdd(exp, rest)
}

func foldAdd(first models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return first, toks, nil
	}

	if toks[0].Type != tokens.PLUS && toks[0].Type != tokens.MINUS {
		return first, toks, nil
	}
	op := toks[0]

	rest = toks[1:]

	var withNext models.Expression
	var next models.Expression
	next, rest, err = parseMulExpression(rest)
	if err != nil {
		return first, rest, err
	}

	withNext = &AddExpression{
		op:     op,
		first:  first,
		second: next,
	}

	return foldAdd(withNext, rest)
}

type MulExpression struct {
	op     tokens.Token
	first  models.Expression
	second models.Expression
}

func (me *MulExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := me.first.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Muller, ok := v1.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("operator '*' cannot be applied to first operand"),
			SourceLocation: me.first.SourceLocation(),
		}
	}

	v2, err := me.second.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v2Muller, ok := v2.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("operator '*' cannot be applied to second operand"),
			SourceLocation: me.second.SourceLocation(),
		}
	}

	switch me.op.Type {
	case tokens.STAR:
		return v1Muller * v2Muller, nil
	case tokens.SLASH:
		return v1Muller / v2Muller, nil
	default:
		return nil, &models.InterpreterError{
			Err:            errors.New("invalid operator"),
			SourceLocation: me.op.SourceLocation,
		}
	}
}

func (me *MulExpression) SourceLocation() models.SourceLocation {
	return me.first.SourceLocation()
}

func parseMulExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	exp, rest, err = parseAtomic(toks)
	if err != nil {
		return nil, rest, err
	}

	return foldMul(exp, rest)
}

func foldMul(first models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return first, toks, nil
	}

	if toks[0].Type != tokens.STAR && toks[0].Type != tokens.SLASH {
		return first, toks, nil
	}

	rest = toks[1:]

	var withNext models.Expression
	var next models.Expression
	next, rest, err = parseAtomic(rest)
	if err != nil {
		return first, rest, err
	}

	withNext = &MulExpression{
		first:  first,
		second: next,
		op:     toks[0],
	}

	return foldMul(withNext, rest)
}

type LiteralExpression struct {
	val any
	loc models.SourceLocation
}

func (le *LiteralExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	if le == nil {
		return nil, nil
	}
	return le.val, nil
}

func (le *LiteralExpression) SourceLocation() models.SourceLocation {
	return le.loc
}

type IdentifierExpression struct {
	name string
	loc  models.SourceLocation
}

func (ie *IdentifierExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	ret, ok := map[string]any(bindings)[ie.name]
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("cannot evaluate unbound identifier"),
			SourceLocation: ie.loc,
		}
	}

	return ret, nil
}

func (ie *IdentifierExpression) SourceLocation() models.SourceLocation {
	return ie.loc
}

/*
func parseIdentifierExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.IDENTIFIER {
		return nil, toks[1:], &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}
	id := toks[0].Value
	rest = toks[1:]

	return &IdentifierExpression{
		name: id,
		loc:  toks[0].SourceLocation,
	}, rest, nil
}
*/

type ArrayLiteralExpression struct {
	val []models.Expression
	loc models.SourceLocation
}

func (ale *ArrayLiteralExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	ret := make([]any, 0)
	for _, v := range ale.val {
		retVal, err := v.Evaluate(bindings)
		if err != nil {
			return nil, err
		}

		ret = append(ret, retVal)
	}

	return ret, nil
}

func (ale *ArrayLiteralExpression) SourceLocation() models.SourceLocation {
	return ale.loc
}

func parseArrayLiteral(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.LEFT_SQUARE_BRACKET {
		return nil, toks, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]
	exps, rest, err := parseExpressions(rest)
	if err != nil {
		return nil, rest, err
	}
	if rest[0].Type != tokens.RIGHT_SQUARE_BRACKET {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}
	rest = rest[1:]
	exp = &ArrayLiteralExpression{
		val: exps,
		loc: toks[0].SourceLocation,
	}

	return exp, rest, nil
}

type FunctionCallExpression struct {
	Function models.Expression
	Args     []models.Expression
	loc      models.SourceLocation
}

func (fce *FunctionCallExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	f, err := fce.Function.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	fun, ok := f.(models.Function)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("cannot call non-function"),
			SourceLocation: fce.Function.SourceLocation(),
		}
	}

	argArray := make([]any, len(fce.Args))
	for i, arg := range fce.Args {
		val, err := arg.Evaluate(bindings)
		if err != nil {
			return nil, err
		}

		argArray[i] = val
	}

	return fun.Call(argArray)
}

func (fce *FunctionCallExpression) SourceLocation() models.SourceLocation {
	return fce.loc
}

func parseAtomic(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("expected token"),
		}
	}

	switch toks[0].Type {
	case tokens.FUNC:
		exp, rest, err = parseFunction(toks)
	case tokens.LEFT_PAREN:
		rest = toks[1:]
		exp, rest, err = ParseExpression(rest)
		if err != nil {
			return nil, rest, err
		}
		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("expected closing parenthesis"),
			}
		}
		if rest[0].Type != tokens.RIGHT_PAREN {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("expected closing parenthesis"),
				SourceLocation: rest[0].SourceLocation,
			}
		}
		rest = rest[1:]
	case tokens.NUMBER, tokens.MINUS:
		var numStr string
		rest = toks
		if toks[0].Type == tokens.MINUS {
			rest = toks[1:]
			if len(rest) == 0 {
				return nil, rest, &models.InterpreterError{
					Err: errors.New("unexpected end of input"),
				}
			}

			numStr = "-"
		}

		if rest[0].Type != tokens.NUMBER {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}

		numStr += rest[0].Value
		rest = rest[1:]

		ret, innerErr := strconv.Atoi(numStr)
		if innerErr != nil {
			return nil, toks[1:], &models.InterpreterError{
				Err:            fmt.Errorf("failed to parse number literal: %w", innerErr),
				SourceLocation: toks[0].SourceLocation,
			}
		}

		exp = &LiteralExpression{
			val: ret,
			loc: toks[0].SourceLocation,
		}
	case tokens.IDENTIFIER:
		if toks[0].Value == "true" {
			exp, rest, err = &LiteralExpression{
				val: true,
				loc: toks[0].SourceLocation,
			}, toks[1:], nil
		} else if toks[0].Value == "false" {
			exp, rest, err = &LiteralExpression{
				val: false,
				loc: toks[0].SourceLocation,
			}, toks[1:], nil
		} else {
			exp, rest, err = &IdentifierExpression{
				name: toks[0].Value,
				loc:  toks[0].SourceLocation,
			}, toks[1:], nil
		}
	case tokens.LEFT_SQUARE_BRACKET:
		exp, rest, err = parseArrayLiteral(toks)
	case tokens.STRING:
		exp, rest, err = &LiteralExpression{
			val: toks[0].Value,
			loc: toks[0].SourceLocation,
		}, toks[1:], nil
	case tokens.LET:
		exp, rest, err = parseLetExpression(toks)
	case tokens.IF:
		exp, rest, err = parseIfExpression(toks)
	default:
		return nil, toks, nil
	}
	if err != nil {
		return nil, rest, err
	}

	if len(rest) == 0 {
		return exp, rest, nil
	}

	if rest[0].Type == tokens.LEFT_PAREN {
		rest = rest[1:]

		var exps []models.Expression
		exps, rest, err = parseExpressions(rest)
		if err != nil {
			return nil, rest, err
		}
		if rest[0].Type != tokens.RIGHT_PAREN {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}
		rest = rest[1:]
		exp = &FunctionCallExpression{
			Function: exp,
			Args:     exps,
			loc: models.SourceLocation{
				LineNumber:   exp.SourceLocation().LineNumber,
				ColumnNumber: exp.SourceLocation().ColumnNumber,
			},
		}
	}

	if len(rest) == 0 {
		return exp, rest, nil
	}

	if rest[0].Type == tokens.FOR {
		exp, rest, err = parseForExpression(exp, rest)
	}

	return exp, rest, err
}

func parseExpressions(toks []tokens.Token) (exps []models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	exps = make([]models.Expression, 0)
	var exp models.Expression
	for exp, rest, err = ParseExpression(toks); err == nil; exp, rest, err = ParseExpression(rest) {
		if exp == nil {
			return exps, rest, nil
		}

		exps = append(exps, exp)
		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("unexpected end of input"),
			}
		}

		if rest[0].Type != tokens.COMMA {
			break
		}
		rest = rest[1:]
	}
	if err != nil {
		return nil, rest, err
	}
	return exps, rest, nil
}
