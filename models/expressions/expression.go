package expressions

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/types"
)

type Bindings map[string]any

type Expression interface {
	Evaluate(Bindings) (any, *models.InterpreterError)
	Type(types.TypeBindings) (types.Type, *models.InterpreterError)
	SourceLocation() *models.SourceLocation
}
