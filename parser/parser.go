package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/brandonksides/crazy-types/models"
	"github.com/brandonksides/crazy-types/tokens"
)

// generation rules:
//
// ROOT: EXPRESSION
//
// EXPRESSION => LET DECLARATION EQUAL EXPRESSION IN EXPRESSION
// EXPRESSION => IF EXPRESSION THEN EXPRESSION ELSE EXPRESSION
// EXPRESSION => FOR_EXPRESSION
//
// EXPRESSION => FOR_EXPRESSION FOR DECLARATION IN EXPRESSION
// EXPRESSION => FOR_EXPRESSION
//
// FOR_EXPRESSION => OR_EXPRESSION
// FOR_EXPRESSION => OR_EXPRESSION OR FOR_EXPRESSION
//
// OR_EXPRESSION => AND_EXPRESSION
// OR_EXPRESSION => AND_EXPRESSION AND OR_EXPRESSION
//
// CMP_EXPRESSION => ADD_EXPRESSION
// CMP_EXPRESSION => ADD_EXPRESSION CMP_OPERATOR ADD_EXPRESSION
//
// ADD_EXPRESSION => MUL_EXPRESSION
// ADD_EXPRESSION => MUL_EXPRESSION ADD_OPERATOR ADD_EXPRESSION
//
// MUL_EXPRESSION => PREFIXED_EXPRESSION
// MUL_EXPRESSION => PREFIXED_EXPRESSION MUL_OPERATOR MUL_EXPRESSION
//
// PREFIXED_EXPRESSION -> FN_CALL
// PREFIXED_EXPRESSION => PREFIX_OPERATOR FN_CALL
//
// DOT_EXPRESSION => ATOMIC_EXPRESSION
// DOT_EXPRESSION => DOT_EXPRESSION LEFT_PAREN ARGS RIGHT_PAREN
// DOT_EXPRESSION => DOT_EXPRESSION DOT IDENTIFIER
//
// ATOMIC_EXPRESSION => LEFT_PAREN EXPRESSION RIGHT_PAREN
// ATOMIC_EXPRESSION => IDENTIFIER
// ATOMIC_EXPRESSION => LITERAL
// ATOMIC_EXPRESSION => LEFT_SQUARE_BRACKET EXPRESSION FOR DECLARATION IN EXPRESSION RIGHT_SQUARE_BRACKET
// ATOMIC_EXPRESSION => LEFT_SQUARE_BRACKET ARGS RIGHT_SQUARE_BRACKET
//
// PREFIX_OPERATOR => MINUS
// PREFIX_OPERATOR => TILDE
//
// AND => AMPERSAND AMPERSAND
// OR => VERTICAL_BAR VERTICAL_BAR
//
// MUL_OPERATOR => STAR
// MUL_OPERATOR => SLASH
// MUL_OPERATOR => AMPERSAND
//
// ADD_OPERATOR => PLUS
// ADD_OPERATOR => MINUS
// ADD_OPERATOR => VERTICAL_BAR
//
// CMP_OPERATOR => EQUAL
// CMP_OPERATOR => LEFT_ANGLE_BRACKET
// CMP_OPERATOR => RIGHT_ANGLE_BRACKET
// CMP_OPERATOR => LEFT_ANGLE_BRACKET EQUAL
// CMP_OPERATOR => RIGHT_ANGLE_BRACKET EQUAL
//
// ARGS =>
// ARGS => NONEMPTY_ARGS
//
// NONEMPTY_ARGS => EXPRESSION
// NONEMPTY_ARGS => EXPRESSION COMMA NONEMPTY_ARGS
//
// ARGS_DECLARATION =>
// ARGS_DECLARATION => NONEMPTY_ARGS_DECLARATION
// NONEMPTY_ARGS_DECLARATION => DECLARATION
// NONEMPTY_ARGS_DECLARATION => DECLARATION COMMA NONEMPTY_ARGS_DECLARATION
//
// DECLARATION => IDENTIFIER
// DECLARATION => IDENTIFIER LEFT_SQUARE_BRACKET EXPRESSION RIGHT_SQUARE_BRACKET

// SIMPLIFIED:
// EXPRESSION => ADD_EXPRESSION
//
// ADD_EXPRESSION => MUL_EXPRESSION
// ADD_EXPRESSION => MUL_EXPRESSION ADD_OP ADD_EXPRESSION
//
// ADD_OP => +
// ADD_OP => -
//
// MUL_EXPRESSION => ATOMIC_EXPRESSION
// MUL_EXPRESSION => ATOMIC_EXPRESSION MUL_OP MUL_EXPRESSION
//
// MUL_OP => *
// MUL_OP => /
//
// ATOMIC_EXPRESSION => LEFT_PAREN EXPRESSION RIGHT_PAREN
// ATOMIC_EXPRESSION => LITERAL
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
		identifier := rest[0].Value

		if rest[1].Type != tokens.EQUAL {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[1].SourceLocation,
			}
		}

		var exp1 models.Expression
		exp1, rest, err = ParseExpression(rest[2:])
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
	exps := make([]models.Expression, 0)
	for exp, rest, err = ParseExpression(rest); err == nil; exp, rest, err = ParseExpression(rest) {
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

func parseAtomic(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("expected token"),
		}
	}

	switch toks[0].Type {
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
	case tokens.NUMBER:
		ret, innerErr := strconv.Atoi(toks[0].Value)
		if innerErr != nil {
			return nil, toks[1:], &models.InterpreterError{
				Err:            fmt.Errorf("failed to parse number literal: %w", innerErr),
				SourceLocation: toks[0].SourceLocation,
			}
		}

		exp, rest, err = &LiteralExpression{
			val: ret,
			loc: toks[0].SourceLocation,
		}, toks[1:], nil
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
		exp, rest, err = nil, toks, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}
	if err != nil {
		return nil, rest, err
	}

	if len(rest) == 0 {
		return exp, rest, nil
	}

	if rest[0].Type == tokens.FOR {
		return parseForExpression(exp, rest)
	}

	return exp, rest, nil
}
