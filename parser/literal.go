package parser

import "github.com/brandonksides/phonk/models"

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
