package types

import (
	"errors"

	"github.com/brandonksides/crazy-types/models"
	"golang.org/x/exp/constraints"
)

type EmptyExpression[T any] struct{}

func (e *EmptyExpression[T]) Evaluate(models.Bindings) (*T, *models.InterpreterError) {
	return nil, nil
}

type LetExpression[T, U any] struct {
	Declaration Declaration
	Expression1 models.Expression
	Expression2 models.Expression
}

func (e *LetExpression[T, U]) Evaluate(bindings models.Bindings) (*U, *models.InterpreterError) {
	newBindings := make(models.Bindings)
	for k, v := range bindings {
		newBindings[k] = v
	}
	var err *models.InterpreterError
	newBindings[e.Declaration.Identifier], err = e.Expression1.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	return e.Expression2.Evaluate(newBindings)
}

type Declaration struct {
	Identifier string
	Type       models.Expression[bool]
}

type IfExpression[T, U any] struct {
	Condition models.Expression[bool]
	Then      models.Expression[T]
	Else      models.Expression[U]
}

type Either[T, U any] struct {
	val interface{}
}

func (e *Either[T, U]) IsLeft() bool {
	_, ok := e.val.(T)
	return ok
}

func (e *Either[T, U]) IsRight() bool {
	_, ok := e.val.(U)
	return ok
}

func (e *Either[T, U]) Left() (T, bool) {
	t, ok := e.val.(T)
	return t, ok
}

func (e *Either[T, U]) Right() (U, bool) {
	u, ok := e.val.(U)
	return u, ok
}

func (e *IfExpression[T, U]) Evaluate(bindings models.Bindings) (*Either[T, U], *models.InterpreterError) {
	cond, err := e.Condition.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	if cond != nil && *cond {
		ret, err := e.Then.Evaluate(bindings)
		if err != nil {
			return nil, err
		}

		return &Either[T, U]{val: ret}, nil
	}

	ret, err := e.Else.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	return &Either[T, U]{val: ret}, nil
}

type HomogenousIfExpression[T any] struct {
	Underlying models.Expression[Either[T, T]]
}

func (e *HomogenousIfExpression[T]) Evaluate(bindings models.Bindings) (*T, *models.InterpreterError) {
	eith, err := e.Underlying.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	if eith.IsLeft() {
		ret, ok := eith.Left()
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("failed to cast left value"),
				SourceLocation: models.SourceLocation{},
			}
		}

		return &ret, nil
	}

	ret, ok := eith.Right()
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("failed to cast right value"),
			SourceLocation: models.SourceLocation{},
		}
	}
	return &ret, nil
}

type FuncExpression[T any] struct {
	Declaration Declaration
	Args        ArgsDeclaration
	Expression  models.Expression[T]
}

func (e *FuncExpression[T]) Evaluate(bindings models.Bindings) (*func(args []any) (*T, error), *models.InterpreterError) {
	ret := func(args []any) (*T, error) {
		newBindings := make(models.Bindings)
		for k, v := range bindings {
			newBindings[k] = v
		}

		for i, arg := range e.Args.Args {
			newBindings[arg.Identifier] = args[i]
		}

		result, err := e.Expression.Evaluate(newBindings)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return &ret, nil
}

type ForExpression[T, U any] struct {
	Declaration Declaration
	Expression1 models.Expression[U]
	Expression2 models.Expression[[]T]
}

func (e *ForExpression[T, U]) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v2, err := e.Expression2.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	results := make([]any, 0)
	for _, v := range *v2 {
		newBindings := make(models.Bindings)
		for k, v := range bindings {
			newBindings[k] = v
		}
		newBindings[e.Declaration.Identifier] = v
		result, err := e.Expression1.Evaluate(newBindings)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

type ArgsDeclaration struct {
	Args []Declaration
}

type OrExpression struct {
	Expression1 models.Expression[bool]
	Expression2 models.Expression[bool]
}

func (e *OrExpression) Evaluate(bindings models.Bindings) (*bool, *models.InterpreterError) {
	v1, err := e.Expression1.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	if *v1 {
		return v1, nil
	}

	return e.Expression2.Evaluate(bindings)
}

type AndExpression struct {
	Expression1 models.Expression[bool]
	Expression2 models.Expression[bool]
}

func (e *AndExpression) Evaluate(bindings models.Bindings) (*bool, *models.InterpreterError) {
	v1, err := e.Expression1.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	if !*v1 {
		return v1, nil
	}

	return e.Expression2.Evaluate(bindings)
}

type EqualsExpression[C comparable] struct {
	Expression1 models.Expression[C]
	Expression2 models.Expression[C]
	not         bool
}

func (e *EqualsExpression[C]) Evaluate(bindings models.Bindings) (*bool, *models.InterpreterError) {
	v1, err := e.Expression1.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	v2, err := e.Expression2.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	eq := *v1 == *v2
	if e.not {
		eq = !eq
	}

	return &eq, nil
}

type ComparisonExpression[O constraints.Ordered] struct {
	Expression1 models.Expression[O]
	Expression2 models.Expression[O]
	Operator    ComparisonOperator
}

type ComparisonOperator uint8

const (
	OperatorGreaterThan ComparisonOperator = iota
	OperatorGreaterThanEqual
	OperatorLessThan
	OperatorLessThanEqual
)

func (e *ComparisonExpression[C]) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := e.Expression1.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	v2, err := e.Expression2.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	var result bool
	switch e.Operator {
	case OperatorGreaterThan:
		result = *v1 > *v2
	case OperatorGreaterThanEqual:
		result = *v1 >= *v2
	case OperatorLessThan:
		result = *v1 < *v2
	case OperatorLessThanEqual:
		result = *v1 <= *v2
	default:
		return nil, &models.InterpreterError{
			Err:            errors.New("unknown operator"),
			SourceLocation: models.SourceLocation{},
		}
	}

	return &result, nil
}
